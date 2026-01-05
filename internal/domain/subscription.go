package domain

import (
	"errors"
	"time"
)

// Ошибки
var (
	ErrSubscriptionNotFound = errors.New("подписка не найдена")
	ErrInvalidPeriod        = errors.New("дана начала подписки должен быть раньше конца")
	ErrInternal             = errors.New("внутренняя ошибка сервера")
)

// Структура для создания подписки
type Subscription struct {
	ServiceName string     `json:"service_name"`       // Название сервиса
	Price       int        `json:"price"`              // Цена в рублях
	UserID      string     `json:"user_id"`            // UUID пользователя
	StartDate   time.Time  `json:"start_date"`         // Дата начала
	EndDate     *time.Time `json:"end_date,omitempty"` // Дата окончания
}

// Структура для обновления подписки
type UpdateSubscriptionInput struct {
	Price   *int64     // Цена в рублях
	EndDate *time.Time // Дата окончания
}

// Структура ответа при ошибке
type ErrorResponse struct {
	Error   string `json:"error" example:"invalid input"`                 // Краткий код или сообщение
	Details string `json:"details,omitempty" example:"email is required"` // (Опционально) Детали
}
