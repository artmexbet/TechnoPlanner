package domain

import (
	"time"

	"github.com/google/uuid"
)

type Place struct {
	ID          uuid.UUID
	Name        string
	Description string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
