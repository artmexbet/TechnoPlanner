package dto

import "time"

//go:generate easyjson -all technic.go

// TechEquipment DTO модель оборудования (тех сервис)
//
//easyjson:json
type TechEquipment struct {
	ID                        int               `json:"id"`
	CategoryID                int               `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics,omitempty"`
	CreatedAt                 time.Time         `json:"created_at,omitempty"`
	UpdatedAt                 time.Time         `json:"updated_at,omitempty"`
}

// TechEquipmentList DTO модель для списка оборудования
//
//easyjson:json
type TechEquipmentList []TechEquipment

// TechEquipmentCategoryList DTO
//
//easyjson:json
type TechEquipmentCategoryList []TechEquipmentCategory

// TechEquipmentCategory DTO модель категории оборудования (тех сервис)
//
//easyjson:json
type TechEquipmentCategory struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Audit       AuditFields `json:"audit"`
}

// TechEquipmentCreateRequest DTO запроса на создание оборудования
//
//easyjson:json
type TechEquipmentCreateRequest struct {
	CategoryID                int               `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description,omitempty"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics,omitempty"`
}

// TechEquipmentUpdateRequest DTO запроса на обновление оборудования
//
//easyjson:json
type TechEquipmentUpdateRequest struct {
	ID                        int               `json:"id"`
	CategoryID                int               `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description,omitempty"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics,omitempty"`
}

// TechEquipmentDeleteRequest DTO запроса на удаление оборудования
//
//easyjson:json
type TechEquipmentDeleteRequest struct {
	ID int `json:"id"`
}

// TechEquipmentGetByIDRequest DTO запроса на получение оборудования по ID
//
//easyjson:json
type TechEquipmentGetByIDRequest struct {
	ID int `json:"id"`
}

// TechEquipmentGetByCategoryRequest DTO запроса на получение оборудования по категории
//
//easyjson:json
type TechEquipmentGetByCategoryRequest struct {
	CategoryID int `json:"category_id"`
}

// TechEquipmentCategoryCreateRequest DTO запроса на создание категории оборудования
//
//easyjson:json
type TechEquipmentCategoryCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// TechEquipmentCategoryUpdateRequest DTO запроса на обновление категории оборудования
//
//easyjson:json
type TechEquipmentCategoryUpdateRequest struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// TechEquipmentCategoryDeleteRequest DTO запроса на удаление категории оборудования
//
//easyjson:json
type TechEquipmentCategoryDeleteRequest struct {
	ID int `json:"id"`
}

// TechEquipmentCategoryGetByIDRequest DTO запроса на получение категории по ID
//
//easyjson:json
type TechEquipmentCategoryGetByIDRequest struct {
	ID int `json:"id"`
}

// Deprecated: используйте TechEquipment
type Technic = TechEquipment

// Deprecated: используйте TechEquipmentCategory
type TechnicCategory = TechEquipmentCategory

// Deprecated: используйте TechEquipmentCreateRequest
type TechnicCreateRequest = TechEquipmentCreateRequest

// Deprecated: используйте TechEquipmentUpdateRequest
type TechnicUpdateRequest = TechEquipmentUpdateRequest

// Deprecated: используйте TechEquipmentDeleteRequest
type TechnicDeleteRequest = TechEquipmentDeleteRequest

// Deprecated: используйте TechEquipmentGetByIDRequest
type TechnicGetByIDRequest = TechEquipmentGetByIDRequest

// Deprecated: используйте TechEquipmentGetByCategoryRequest
type TechnicGetByCategoryRequest = TechEquipmentGetByCategoryRequest

// Deprecated: используйте TechEquipmentCategoryCreateRequest
type TechnicCategoryCreateRequest = TechEquipmentCategoryCreateRequest

// Deprecated: используйте TechEquipmentCategoryUpdateRequest
type TechnicCategoryUpdateRequest = TechEquipmentCategoryUpdateRequest

// Deprecated: используйте TechEquipmentCategoryDeleteRequest
type TechnicCategoryDeleteRequest = TechEquipmentCategoryDeleteRequest

// Deprecated: используйте TechEquipmentCategoryGetByIDRequest
type TechnicCategoryGetByIDRequest = TechEquipmentCategoryGetByIDRequest
