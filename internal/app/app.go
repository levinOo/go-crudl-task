package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/levinOo/go-crudl-task/internal/config"
	"github.com/levinOo/go-crudl-task/internal/db"
	"github.com/levinOo/go-crudl-task/internal/handlers"
	"github.com/levinOo/go-crudl-task/internal/repository"
	"github.com/levinOo/go-crudl-task/internal/service"
	"github.com/levinOo/go-crudl-task/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Запускаем приложение
func Run() error {
	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Инициализируем логгер
	log := logger.New(cfg.Env)
	slog.SetDefault(log)

	log.Info("URL Database", slog.String("url", cfg.Postgre.URL))

	// Подключаем базу данных
	pgCfg := db.Config{
		URL:            cfg.Postgre.URL,
		PoolMax:        cfg.Postgre.PoolMax,
		RetryAttempts:  cfg.Postgre.RetryAttempts,
		RetryDelay:     cfg.Postgre.RetryDelay,
		ConnectTimeout: cfg.Postgre.ContextTimeoutValue,
	}

	pg, err := db.New(pgCfg, log)
	if err != nil {
		log.Error("Не удалось подключиться к БД", slog.String("error", err.Error()))
		return err
	}
	defer pg.Close()

	// Выполняем миграции
	if err := db.RunMigrations(pg); err != nil {
		log.Error("Не удалось выполнить миграции", slog.String("error", err.Error()))
		return err
	}

	// Dependency Injection
	repo := repository.NewRepositories(pg)

	deps := service.Deps{
		Repos: *repo,
	}
	services := service.NewServices(deps)
	h := handlers.NewHandler(services.Subscription, log)

	// Устанавливаем режим работы сервера
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
		log.Info("Сервер запущен в production режиме")
	} else {
		gin.SetMode(gin.DebugMode)
		log.Debug("Сервер запущен в debug режиме")
	}

	// Инициализация роутеров
	router := gin.New()
	h.InitRoutes(router)

	// Конфигурация HTTP сервера
	srv := &http.Server{
		Addr:         net.JoinHostPort("", cfg.Server.ServerPort),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Запуске сервера
	go func() {
		log.Info("Запуск HTTP сервера", slog.String("addr", cfg.Server.ServerPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Ошибка запуска сервера", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// контекст для Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Ожидание сигнала остановки
	<-ctx.Done()

	log.Info("Получен сигнал остановки, начинаем graceful shutdown...")

	// Graceful Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownContextValue)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Ошибка при остановке сервера", slog.String("error", err.Error()))
		return err
	}

	log.Info("Сервер успешно остановлен")

	return nil
}
