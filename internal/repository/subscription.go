package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/levinOo/go-crudl-task/internal/db"
	"github.com/levinOo/go-crudl-task/internal/domain"

	"github.com/jackc/pgx/v5"
)

// Структура репозитория подписок
type SubscriptionRepository struct {
	pg *db.Postgres
}

// Функция конструктор
func NewSubscriptionRepository(pg *db.Postgres) *SubscriptionRepository {
	return &SubscriptionRepository{pg: pg}
}

// Создание подписки
func (r *SubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (string, error) {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id string

	err := r.pg.Pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Ошибка при создании подписки: %w", err)
	}

	return id, nil
}

// Получение подписки
func (r *SubscriptionRepository) Get(ctx context.Context, id string) (domain.Subscription, error) {
	query := `
		SELECT service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	var sub domain.Subscription

	err := r.pg.Pool.QueryRow(ctx, query, id).Scan(
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrSubscriptionNotFound
		}
		return domain.Subscription{}, fmt.Errorf("Ошибка при получении подписки: %w", err)
	}

	return sub, nil
}

// Обновление подписки
func (r *SubscriptionRepository) Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput) error {
	query := "UPDATE subscriptions SET "
	args := []any{}
	argId := 1

	if input.Price != nil {
		query += fmt.Sprintf("price = $%d, ", argId)
		args = append(args, *input.Price)
		argId++
	}
	if input.EndDate != nil {
		query += fmt.Sprintf("end_date = $%d, ", argId)
		args = append(args, *input.EndDate)
		argId++
	}

	if len(args) == 0 {
		return nil
	}

	query = query[:len(query)-2]

	query += fmt.Sprintf(" WHERE id = $%d", argId)
	args = append(args, id)

	result, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Ошибка при обновлении подписки: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

// Удаление подписки
func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	query := `
	DELETE 
	FROM subscriptions 
	WHERE id = $1
	`

	result, err := r.pg.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении подписки: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

// Получение списка подписок
func (r *SubscriptionRepository) List(ctx context.Context, userID string) ([]domain.Subscription, error) {
	query := `
		SELECT service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
	`

	rows, err := r.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении списка подписок: %w", err)
	}
	defer rows.Close()

	subs := make([]domain.Subscription, 0)

	for rows.Next() {
		var sub domain.Subscription

		if err := rows.Scan(
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании списка подписок: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при сканировании списка подписок: %w", err)
	}

	return subs, nil
}

// Получение суммы подписок
func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1
		AND ($2 = '' OR service_name = $2)
		AND start_date <= $4
		AND (end_date IS NULL OR end_date >= $3)
	`

	var total int
	err := r.pg.Pool.QueryRow(ctx, query, userID, serviceName, startDate, endDate).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при подсчете суммы подписок: %w", err)
	}

	return total, nil
}
