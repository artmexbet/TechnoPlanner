package domain

import "github.com/google/uuid"

//go:generate easyjson -all domain.go

type Technic struct {
	ID                        uuid.UUID         `json:"id"`
	CategoryID                uuid.UUID         `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics"`
}

type TechnicCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
