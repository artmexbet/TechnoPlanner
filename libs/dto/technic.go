package dto

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all technic.go

// Technic DTO модель техники
//
//easyjson:json
type Technic struct {
	ID                        uuid.UUID         `json:"id"`
	CategoryID                uuid.UUID         `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics,omitempty"`
	CreatedAt                 time.Time         `json:"created_at,omitempty"`
	UpdatedAt                 time.Time         `json:"updated_at,omitempty"`
}

// TechnicCategory DTO модель категории техники
//
//easyjson:json
type TechnicCategory struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Audit       AuditFields `json:"audit"`
}

// TechnicCreateRequest DTO запроса на создание техники
//
//easyjson:json
type TechnicCreateRequest struct {
	CategoryID                uuid.UUID         `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description,omitempty"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics,omitempty"`
}

// TechnicUpdateRequest DTO запроса на обновление техники
//
//easyjson:json
type TechnicUpdateRequest struct {
	ID                        uuid.UUID         `json:"id"`
	CategoryID                uuid.UUID         `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description,omitempty"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics,omitempty"`
}

// TechnicDeleteRequest DTO запроса на удаление техники
//
//easyjson:json
type TechnicDeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

// TechnicGetByIDRequest DTO запроса на получение техники по ID
//
//easyjson:json
type TechnicGetByIDRequest struct {
	ID uuid.UUID `json:"id"`
}

// TechnicGetByCategoryRequest DTO запроса на получение техники по категории
//
//easyjson:json
type TechnicGetByCategoryRequest struct {
	CategoryID uuid.UUID `json:"category_id"`
}

// TechnicCategoryCreateRequest DTO запроса на создание категории техники
//
//easyjson:json
type TechnicCategoryCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// TechnicCategoryUpdateRequest DTO запроса на обновление категории техники
//
//easyjson:json
type TechnicCategoryUpdateRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

// TechnicCategoryDeleteRequest DTO запроса на удаление категории техники
//
//easyjson:json
type TechnicCategoryDeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

// TechnicCategoryGetByIDRequest DTO запроса на получение категории по ID
//
//easyjson:json
type TechnicCategoryGetByIDRequest struct {
	ID uuid.UUID `json:"id"`
}
