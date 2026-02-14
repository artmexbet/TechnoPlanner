package dto

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all events.go

// UserCreatedEvent событие создания пользователя от auth сервиса
type UserCreatedEvent struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RoleID    int32     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
