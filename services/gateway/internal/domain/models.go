package domain

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson models.go

type Place struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//easyjson:json
type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	RoleID       int32      `json:"role_id"`
	Role         *Role      `json:"role,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-"`
}

// Equipment представляет оборудование (синхронизировано с Equipment Service)
type Equipment struct {
	ID                        int
	CategoryID                int
	Name                      string
	Description               string
	Quantity                  int32
	AdditionalCharacteristics map[string]string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

// EquipmentCategory представляет категорию оборудования (синхронизировано с Equipment Service)
type EquipmentCategory struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Request struct {
	ID                uuid.UUID
	TelegramUserInfo  []byte
	RequestText       *string
	Status            RequestStatus
	ScheduleTime      string
	EndTime           time.Time
	Address           string
	ResponsibleUserID *uuid.UUID
	Equipment         []RequestEquipment
	Audit             AuditFields
}

type RequestEquipment struct {
	RequestID   uuid.UUID
	EquipmentID int32
	Name        string
	Description *string
	Quantity    int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RequestStatus string

const (
	RequestStatusCanceled   RequestStatus = "canceled"
	RequestStatusPending    RequestStatus = "pending"
	RequestStatusAssigned   RequestStatus = "assigned"
	RequestStatusInProgress RequestStatus = "in_progress"
	RequestStatusCompleted  RequestStatus = "completed"
	RequestStatusRejected   RequestStatus = "rejected"
)

type RequestStatusHistory struct {
	ID        int32
	RequestID uuid.UUID
	Status    RequestStatus
	Comment   *string
	ChangedBy *uuid.UUID
	ChangedAt time.Time
}

type Role struct {
	ID          int32
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AuditFields struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy *uuid.UUID
	UpdatedBy *uuid.UUID
}

func NewAuditFields(createdAt, updatedAt time.Time, createdBy, updatedBy *uuid.UUID, deletedAt *time.Time) AuditFields {
	return AuditFields{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}
}

type Porter struct {
	ID       uuid.UUID
	Username string
}

type RequestUpdate struct {
	RequestText   *string
	Status        *RequestStatus
	ScheduleTime  *string
	Address       *string
	ResponsibleID *uuid.UUID
}

// RawRequest — необработанный запрос от Telegram-бота
type RawRequest struct {
	ID                 uuid.UUID
	TelegramID         int64
	Username           string
	FirstName          string
	LastName           *string
	RawText            string
	Status             string
	ProcessedRequestID *uuid.UUID
	CreatedAt          time.Time
}
