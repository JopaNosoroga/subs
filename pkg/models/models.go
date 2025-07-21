package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscriptions представляет модель подписки пользователя
// swagger:model subscription
type Subscriptions struct {
	// Идентификатор подписки
	ID int `json:"id"`
	// UUID пользователя (формат uuidv4)
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID uuid.UUID `json:"user_id"`
	// Название сервиса подписки
	// example: Netflix
	ServiceName string `json:"service_name"`

	// Стоимость подписки (в копейках/центах)
	// example: 1000
	Price int `json:"price"`
	// Дата начала подписки (RFC3339)
	// example: 2023-01-15T00:00:00Z
	StartDate time.Time `json:"start_date"`
	// Дата окончания подписки (RFC3339)
	// example: 2023-02-15T00:00:00Z
	EndDate time.Time `json:"end_date"`
	// Флаг для запроса всех подписок (используется только в запросах)
	All bool `json:"all"`
}

// ErrorResponse представляет модель ошибки API
// swagger:model errorResponse
type ErrorResponse struct {
	// Сообщение об ошибке
	// example: Ошибка сервера
	Message string `json:"message"`
	// HTTP статус код
	// example: 500
	Code int `json:"code"`
}
