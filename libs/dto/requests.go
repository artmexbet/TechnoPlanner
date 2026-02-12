package dto

import (
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
