package dto

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all requests.go

// IDRequest запрос по ID (int32)
//
//easyjson:json
type IDRequest struct {
	ID int32 `json:"id"`
}

// UUIDRequest запрос по UUID
//
//easyjson:json
type UUIDRequest struct {
	ID uuid.UUID `json:"id"`
}

// RoleIDRequest запрос по RoleID
//
//easyjson:json
type RoleIDRequest struct {
	RoleID int32 `json:"role_id"`
}

// SoftDeleteRequest запрос на мягкое удаление
//
//easyjson:json
type SoftDeleteRequest struct {
	ID     int32      `json:"id"`
	UserID *uuid.UUID `json:"user_id,omitempty"`
}

// RequestIDRequest запрос по RequestID
//
//easyjson:json
type RequestIDRequest struct {
	RequestID uuid.UUID `json:"request_id"`
}

// RequestCreateRequest DTO запроса на создание заявки (от бота)
//
//easyjson:json
type RequestCreateRequest struct {
	Text            *string         `json:"text,omitempty"`
	ScheduleTime    string          `json:"schedule_time"`
	TelegramUserID  int64           `json:"user_id"`
	Username        *string         `json:"user_name,omitempty"`
	Equipments      []EquipmentInfo `json:"equipments,omitempty"`
	EquipmentString *string         `json:"equipment_string,omitempty"`
	Address         string          `json:"address"`
}

// EquipmentInfo информация об оборудовании в запросе
//
//easyjson:json
type EquipmentInfo struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

// RequestStatusUpdateRequest DTO запроса на обновление статуса
//
//easyjson:json
type RequestStatusUpdateRequest struct {
	RequestID uuid.UUID `json:"request_id"`
	Status    string    `json:"status"`
}

// RequestListByTelegramRequest DTO запроса списка заявок по Telegram ID
//
//easyjson:json
type RequestListByTelegramRequest struct {
	TelegramID int64 `json:"telegram_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

// EquipmentAddRequest DTO запроса на добавление оборудования
//
//easyjson:json
type EquipmentAddRequest struct {
	Equipments []EquipmentItem `json:"equipments"`
}

// EquipmentItem элемент оборудования
//
//easyjson:json
type EquipmentItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// RawRequestCreateRequest DTO сырого запроса от Telegram-бота
//
//easyjson:json
type RawRequestCreateRequest struct {
	TelegramID int64   `json:"telegram_id" validate:"required"`
	Username   string  `json:"username"`
	FirstName  string  `json:"first_name"`
	LastName   *string `json:"last_name,omitempty"`
	RawText    string  `json:"raw_text" validate:"required"`
}

// RawRequest DTO сырого запроса (ответ)
//
//easyjson:json
type RawRequest struct {
	ID                 string  `json:"id"`
	TelegramID         int64   `json:"telegram_id"`
	Username           string  `json:"username"`
	FirstName          string  `json:"first_name"`
	LastName           *string `json:"last_name,omitempty"`
	RawText            string  `json:"raw_text"`
	Status             string  `json:"status"`
	ProcessedRequestID *string `json:"processed_request_id,omitempty"`
	CreatedAt          string  `json:"created_at"`
}

// RawRequestListRequest DTO запроса списка сырых запросов
//
//easyjson:json
type RawRequestListRequest struct {
	Status string `json:"status"` // "new", "processed" или "" (все)
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

// RawRequestProcessRequest DTO запроса на обработку сырого запроса (создание нормальной заявки)
//
//easyjson:json
type RawRequestProcessRequest struct {
	RawRequestID    string          `json:"raw_request_id" validate:"required,uuid4"`
	RequestText     *string         `json:"request_text,omitempty"`
	ScheduleTime    string          `json:"schedule_time" validate:"required"`
	EndTime         *time.Time      `json:"end_time,omitempty"`
	Address         string          `json:"address" validate:"required"`
	Equipments      []EquipmentInfo `json:"equipments,omitempty"`
	EquipmentString *string         `json:"equipment_string,omitempty"`
}
