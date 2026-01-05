package service

import (
	"context"
	"time"

	"github.com/levinOo/go-crudl-task/internal/domain"
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

// Структура сервиса подписок
type SubscriptionServiceImplementation struct {
	repo SubscriptionRepo
}

// Функция конструктор сервиса подписок
func NewSubscriptionService(repo SubscriptionRepository) *SubscriptionServiceImplementation {
	return &SubscriptionServiceImplementation{
		repo: repo,
	}
}

// Функция создания подписки
func (s *SubscriptionServiceImplementation) Create(ctx context.Context, sub domain.Subscription) (string, error) {
	id, err := s.repo.Create(ctx, sub)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Функция получения подписки
func (s *SubscriptionServiceImplementation) Get(ctx context.Context, id string) (domain.Subscription, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Subscription{}, err
	}

	return sub, nil
}

// Функция обновления подписки
func (s *SubscriptionServiceImplementation) Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput) error {
	if input.EndDate != nil {
		currentSub, err := s.repo.Get(ctx, id)
		if err != nil {
			return err
		}

		if input.EndDate.Before(currentSub.StartDate) {
			return domain.ErrInvalidPeriod
		}
	}

	return s.repo.Update(ctx, id, input)
}

// Функция удаления подписки
func (s *SubscriptionServiceImplementation) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Функция получения списка подписок
func (s *SubscriptionServiceImplementation) List(ctx context.Context, userID string) ([]domain.Subscription, error) {
	list, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Функция получения общей стоимости подписок
func (s *SubscriptionServiceImplementation) GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error) {
	total, err := s.repo.GetTotalCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		return 0, err
	}

	return total, nil
}
