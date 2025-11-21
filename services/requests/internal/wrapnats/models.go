package wrapnats

import (
	"github.com/google/uuid"

	"requests/internal/domain"
)

//go:generate easyjson -all models.go

type statusCode int

const (
	statusOK                  statusCode = 200
	statusBadRequest          statusCode = 400
	statusUnauthorized        statusCode = 401
	statusForbidden           statusCode = 403
	statusNotFound            statusCode = 404
	statusInternalServerError statusCode = 500
)

type response struct {
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload"`
	StatusCode statusCode  `json:"statusCode"`
	IsError    bool        `json:"isError"`
}

func newResponse(message string, statusCode statusCode, payload interface{}, isError ...bool) *response {
	var isErr bool
	if len(isError) > 0 {
		isErr = isError[0]
	}
	return &response{
		Message:    message,
		Payload:    payload,
		StatusCode: statusCode,
		IsError:    isErr,
	}
}

type equipmentInfo struct {
	ID       int `json:"id" validate:"required,gt=0"`
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

type requestCreate struct {
	Text         *string         `json:"text" validate:"omitempty,max=1000"`
	ScheduleTime string          `json:"schedule_time" validate:"required"`
	TelegramID   int64           `json:"telegram_id" validate:"required"`
	Equipments   []equipmentInfo `json:"equipments" validate:"required,dive"`
	Address      string          `json:"address" validate:"required,max=1000"`
}

func (r *requestCreate) ToDomain() domain.Request {
	req := domain.Request{
		RequestText:  r.Text,
		ScheduleTime: r.ScheduleTime,
		Equipments:   make([]domain.Equipment, len(r.Equipments)),
		Issuer: domain.User{
			TelegramID: r.TelegramID,
		},
		Address: r.Address,
	}
	for i, eq := range r.Equipments {
		req.Equipments[i] = domain.Equipment{ID: eq.ID, Quantity: eq.Quantity}
	}
	return req
}

type requestUpdateStatus struct {
	RequestID uuid.UUID `json:"request_id" validate:"required,uuid4"`
	Status    string    `json:"status" validate:"required,oneof=canceled pending assigned in_progress completed rejected"`
}

type requestByID struct {
	RequestID uuid.UUID `json:"request_id" validate:"required,uuid4"`
}

type requestList struct {
	TelegramID int64 `json:"telegram_id" validate:"required"`
	Limit      int32 `json:"limit" validate:"required,min=1"`
	Offset     int32 `json:"offset" validate:"min=0"`
}

type addEquipment struct {
	Equipments []domain.Equipment `json:"equipments" validate:"required,dive"`
}
