package repository

import (
	"context"

	"github.com/levinOo/go-crudl-task/internal/db"
	"github.com/levinOo/go-crudl-task/internal/domain"

	"time"
)

// Интерфейс репозитория подписок
type SubscriptionRepo interface {
	Create(ctx context.Context, sub domain.Subscription) (string, error)
	Get(ctx context.Context, id string) (domain.Subscription, error)
	Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string) ([]domain.Subscription, error)
	GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error)
}

// Структура слоя репозиториев
type Repositories struct {
	Subscription SubscriptionRepo
}

// Функция конструктор слоя репозиториев
func NewRepositories(pg *db.Postgres) *Repositories {
	return &Repositories{
		Subscription: NewSubscriptionRepository(pg),
	}
}
