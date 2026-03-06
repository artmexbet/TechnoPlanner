package domain

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all domain.go

type StatusType string

func (s *StatusType) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

const (
	StatusCanceled   StatusType = "canceled"
	StatusPending    StatusType = "pending"
	StatusAssigned   StatusType = "assigned"
	StatusInProgress StatusType = "in_progress"
	StatusCompleted  StatusType = "completed"
	StatusRejected   StatusType = "rejected"
)

func (s *StatusType) Set(value interface{}) error {
	sVal, ok := value.(string)
	if !ok {
		return ErrInvalidStatus
	}
	*s = StatusType(sVal)
	return nil
}

type Request struct {
	ID              uuid.UUID   `json:"id"`
	RequestText     *string     `json:"request_text"`
	Status          StatusType  `json:"status"`
	Equipments      []Equipment `json:"equipments"`
	EquipmentString *string     `json:"equipment_string"`
	Issuer          User        `json:"issuer"`
	ScheduleTime    string      `json:"schedule_time"`
	EndTime         time.Time   `json:"end_time"`
	Address         string      `json:"address"`
	PorterInfo      *PorterInfo `json:"porter_info,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type PorterInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}

type User struct {
	ID         uuid.UUID `json:"id,omitempty"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   *string   `json:"last_name"`
	Email      string    `json:"email,omitempty"`
	RoleID     int32     `json:"role_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type Equipment struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      *string   `json:"description"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

type RequestUpdate struct {
	RequestText  *string     `json:"request_text,omitempty"`
	Status       *StatusType `json:"status,omitempty"`
	ScheduleTime *string     `json:"schedule_time,omitempty"`
	Address      *string     `json:"address,omitempty"`
	PorterID     *uuid.UUID  `json:"porter_id,omitempty"`
}

type Porter struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// EquipmentReserveItem — позиция для резервации/освобождения оборудования
type EquipmentReserveItem struct {
	EquipmentID int `json:"equipment_id"`
	Quantity    int `json:"quantity"`
}

// RawRequestStatus статус сырого запроса от бота
type RawRequestStatus string

const (
	RawRequestStatusNew       RawRequestStatus = "new"
	RawRequestStatusProcessed RawRequestStatus = "processed"
)

// RawRequest — необработанный запрос от Telegram-бота
type RawRequest struct {
	ID                 uuid.UUID        `json:"id"`
	TelegramID         int64            `json:"telegram_id"`
	Username           string           `json:"username"`
	FirstName          string           `json:"first_name"`
	LastName           *string          `json:"last_name,omitempty"`
	RawText            string           `json:"raw_text"`
	Status             RawRequestStatus `json:"status"`
	ProcessedRequestID *uuid.UUID       `json:"processed_request_id,omitempty"`
	CreatedAt          time.Time        `json:"created_at"`
}
