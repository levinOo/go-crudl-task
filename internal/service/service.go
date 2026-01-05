package service

import (
	"context"
	"time"

	"github.com/levinOo/go-crudl-task/internal/domain"
	"github.com/levinOo/go-crudl-task/internal/repository"
)

// Интерфейс репозитория подписок
type SubscriptionRepository interface {
	Create(ctx context.Context, sub domain.Subscription) (string, error)
	Get(ctx context.Context, id string) (domain.Subscription, error)
	Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string) ([]domain.Subscription, error)
	GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error)
}

// Интерфейс сервиса подписок
type SubscriptionService interface {
	Create(ctx context.Context, sub domain.Subscription) (string, error)
	Get(ctx context.Context, id string) (domain.Subscription, error)
	Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string) ([]domain.Subscription, error)
	GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error)
}

// Структура сервисов
type Services struct {
	Subscription SubscriptionService
}

// Структура зависимостей
type Deps struct {
	Repos repository.Repositories
}

// Функция конструктор сервисов
func NewServices(deps Deps) *Services {
	return &Services{
		Subscription: NewSubscriptionService(deps.Repos.Subscription),
	}
}
