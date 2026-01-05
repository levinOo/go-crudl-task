package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/avast/retry-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Конфигурация базы данных
type Config struct {
	URL            string
	PoolMax        int
	ConnectTimeout time.Duration
	RetryAttempts  int
	RetryDelay     time.Duration
}

// Структура базы данных
type Postgres struct {
	Pool *pgxpool.Pool
}

// Подключение к Базе данных
func New(cfg Config, log *slog.Logger) (*Postgres, error) {
	pg := &Postgres{}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.PoolMax)

	err = retry.Do(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
			defer cancel()

			pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
			if err != nil {
				return err
			}

			if err := pool.Ping(ctx); err != nil {
				pool.Close()
				return err
			}

			pg.Pool = pool
			return nil
		},
		retry.Attempts(uint(cfg.RetryAttempts)),
		retry.Delay(cfg.RetryDelay),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Warn("Postgres connection retry",
				slog.Int("attempt", int(n)+1),
				slog.Int("max_attempts", cfg.RetryAttempts),
				slog.String("error", err.Error()),
			)
		}),
	)

	if err != nil {
		return nil, err
	}

	return pg, err
}

// Закрытие подключения к Базе данных
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
