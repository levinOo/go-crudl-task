package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/levinOo/go-crudl-task/internal/domain"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/levinOo/go-crudl-task/docs"
)

// Интерфейс сервиса подписок
type SubscriptionService interface {
	Create(ctx context.Context, sub domain.Subscription) (string, error)
	Get(ctx context.Context, id string) (domain.Subscription, error)
	Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string) ([]domain.Subscription, error)
	GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error)
}

// Структура хендлера
type Handler struct {
	services SubscriptionService
	log      *slog.Logger
}

// Создание нового хендлера
func NewHandler(services SubscriptionService, log *slog.Logger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

// Инициализация маршрутов
func (h *Handler) InitRoutes(router *gin.Engine) {
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			subs := v1.Group("/subscriptions")
			{
				subs.POST("", h.createSubscription)
				subs.GET("", h.getList)

				subs.GET("/:id", h.getSubscription)
				subs.PATCH("/:id", h.updateSubscription)
				subs.DELETE("/:id", h.deleteSubscription)
				subs.GET("/total-cost", h.getTotalCost)
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
