package domain

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all domain.go

type Technic struct {
	ID                        uuid.UUID         `json:"id"`
	CategoryID                uuid.UUID         `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics"`
	CreatedAt                 time.Time         `json:"created_at,omitempty"`
	UpdatedAt                 time.Time         `json:"updated_at,omitempty"`
}

type TechnicCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
