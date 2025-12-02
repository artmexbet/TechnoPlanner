package domain

import (
	"time"

	"github.com/google/uuid"
)

type Place struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	RoleID       int32
	Role         *Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type Equipment struct {
	ID          int32
	Name        string
	Description *string
	Quantity    int32
	Categories  []EquipmentCategory
	Audit       AuditFields
}

type EquipmentCategory struct {
	ID          int32
	Name        string
	Description *string
	Audit       AuditFields
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
