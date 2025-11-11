package domain

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all domain.go

type StatusType string

func (s *StatusType) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

const (
	StatusPending StatusType = "pending"
)

func (s *StatusType) Set(value string) error {
	if value != string(StatusPending) {
		return ErrInvalidStatus
	}
	*s = StatusType(value)
	return nil
}

type Request struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	RequestText string     `json:"request_text"`
	Status      StatusType `json:"status"`
	Technics    []Technic  `json:"technics"`
	Issuer      User       `json:"issuer"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Technic struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
