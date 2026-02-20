package domain

import "time"

//go:generate easyjson -all domain.go

// Equipment представляет оборудование/технику
type Equipment struct {
	ID                        int               `json:"id"`
	CategoryID                int               `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics"`
	CreatedAt                 time.Time         `json:"created_at,omitempty"`
	UpdatedAt                 time.Time         `json:"updated_at,omitempty"`
}

// EquipmentCategory представляет категорию оборудования
type EquipmentCategory struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
