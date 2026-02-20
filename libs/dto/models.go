// Package dto содержит общие DTO модели для взаимодействия между сервисами
package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

// RequestStatus статус заявки
type RequestStatus string

const (
	RequestStatusCanceled   RequestStatus = "canceled"
	RequestStatusPending    RequestStatus = "pending"
	RequestStatusAssigned   RequestStatus = "assigned"
	RequestStatusInProgress RequestStatus = "in_progress"
	RequestStatusCompleted  RequestStatus = "completed"
	RequestStatusRejected   RequestStatus = "rejected"
)

// Request DTO модель заявки
//
//easyjson:json
type Request struct {
	ID                uuid.UUID          `json:"id"`
	TelegramUserInfo  []byte             `json:"telegram_user_info,omitempty"`
	RequestText       *string            `json:"request_text,omitempty"`
	Status            RequestStatus      `json:"status"`
	ScheduleTime      string             `json:"schedule_time"`
	EndTime           time.Time          `json:"end_time"`
	Address           string             `json:"address"`
	ResponsibleUserID *uuid.UUID         `json:"responsible_user_id,omitempty"`
	Equipment         []RequestEquipment `json:"equipment,omitempty"`
	Audit             AuditFields        `json:"audit"`
}

// RequestEquipment DTO оборудования в заявке
//
//easyjson:json
type RequestEquipment struct {
	RequestID   uuid.UUID `json:"request_id"`
	EquipmentID int32     `json:"equipment_id"`
	Quantity    int32     `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RequestStatusHistory DTO истории статусов заявки
//
//easyjson:json
type RequestStatusHistory struct {
	ID        int32         `json:"id"`
	RequestID uuid.UUID     `json:"request_id"`
	Status    RequestStatus `json:"status"`
	Comment   *string       `json:"comment,omitempty"`
	ChangedBy *uuid.UUID    `json:"changed_by,omitempty"`
	ChangedAt time.Time     `json:"changed_at"`
}

// User DTO модель пользователя
//
//easyjson:json
type User struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	RoleID    int32      `json:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// UserWithRole DTO пользователя с ролью
//
//easyjson:json
type UserWithRole struct {
	User
	Role *Role `json:"role,omitempty"`
}

// Role DTO роли
//
//easyjson:json
type Role struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Equipment DTO оборудования
//
//easyjson:json
type Equipment struct {
	ID          int32               `json:"id"`
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	Quantity    int32               `json:"quantity"`
	Categories  []EquipmentCategory `json:"categories,omitempty"`
	Audit       AuditFields         `json:"audit"`
}

// EquipmentCategory DTO категории оборудования
//
//easyjson:json
type EquipmentCategory struct {
	ID          int32       `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Audit       AuditFields `json:"audit"`
}

// AuditFields DTO аудит полей
//
//easyjson:json
type AuditFields struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

// Place DTO места
//
//easyjson:json
type Place struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RequestListRequest DTO запроса на получение списка заявок
//
//easyjson:json
type RequestListRequest struct {
	ResponsibleID *uuid.UUID `json:"responsible_id,omitempty"`
}

// RequestByIDRequest DTO запроса на получение заявки по ID
//
//easyjson:json
type RequestByIDRequest struct {
	RequestID uuid.UUID `json:"request_id"`
}

// AssignResponsibleRequest DTO запроса на назначение ответственного
//
//easyjson:json
type AssignResponsibleRequest struct {
	RequestID     uuid.UUID  `json:"request_id"`
	ResponsibleID *uuid.UUID `json:"responsible_id"`
}

// GatewayResponse стандартный ответ от сервисов
//
//easyjson:json
type GatewayResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// CategoryCreateRequest DTO запроса на создание категории
//
//easyjson:json
type CategoryCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CategoryUpdateRequest DTO запроса на обновление категории
//
//easyjson:json
type CategoryUpdateRequest struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// EquipmentCreateRequest DTO запроса на создание оборудования
//
//easyjson:json
type EquipmentCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Quantity    int32   `json:"quantity"`
	CategoryIDs []int32 `json:"category_ids,omitempty"`
}

// EquipmentUpdateRequest DTO запроса на обновление оборудования
//
//easyjson:json
type EquipmentUpdateRequest struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Quantity    int32   `json:"quantity"`
	CategoryIDs []int32 `json:"category_ids,omitempty"`
}

// UserUpdateRequest DTO запроса на обновление пользователя
//
//easyjson:json
type UserUpdateRequest struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
}

// UserDeleteRequest DTO запроса на удаление пользователя
//
//easyjson:json
type UserDeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

//easyjson:json
type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int32  `json:"role_id"`
}

// HistoryAddRequest DTO запроса на добавление записи в историю
//
//easyjson:json
type HistoryAddRequest struct {
	RequestID uuid.UUID     `json:"request_id"`
	Status    RequestStatus `json:"status"`
	Comment   string        `json:"comment,omitempty"`
}

// Responsible DTO ответственного
//
//easyjson:json
type Responsible struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// ResponsibleCreateRequest DTO запроса на создание ответственного
//
//easyjson:json
type ResponsibleCreateRequest struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// RequestUpdateRequest DTO запроса на обновление заявки
//
//easyjson:json
type RequestUpdateRequest struct {
	RequestID     uuid.UUID      `json:"request_id"`
	RequestText   *string        `json:"request_text,omitempty"`
	Status        *RequestStatus `json:"status,omitempty"`
	ScheduleTime  *string        `json:"schedule_time,omitempty"`
	Address       *string        `json:"address,omitempty"`
	ResponsibleID *uuid.UUID     `json:"responsible_id,omitempty"`
}
