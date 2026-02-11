package wrapnats

import (
	"tech/internal/domain"

	"github.com/google/uuid"
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

type requestUpdateTech struct {
	ID                        uuid.UUID         `json:"id"`
	CategoryID                uuid.UUID         `json:"category_id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	AdditionalCharacteristics map[string]string `json:"additional_characteristics"`
}

func (t requestUpdateTech) ToDomain() domain.Technic {
	return domain.Technic{
		ID:                        t.ID,
		CategoryID:                t.CategoryID,
		Name:                      t.Name,
		Description:               t.Description,
		AdditionalCharacteristics: t.AdditionalCharacteristics,
	}
}

type requestSpecificTech struct {
	ID uuid.UUID `json:"id"`
}

type requestUpdateCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type requestSpecificCategory struct {
	ID uuid.UUID `json:"id"`
}
