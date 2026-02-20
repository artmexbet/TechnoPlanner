package wrapnats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/dto"
)

type statusCode int

const (
	statusOK                  statusCode = 200
	statusBadRequest          statusCode = 400
	statusUnauthorized        statusCode = 401
	statusForbidden           statusCode = 403
	statusNotFound            statusCode = 404
	statusInternalServerError statusCode = 500
)

func respondError(msg *nats.Msg, message string, payload interface{}, code statusCode) error {
	resp := dto.GatewayResponse{
		Success: false,
		Message: message,
	}

	// Если payload - это string, просто используем его как сообщение
	if _, ok := payload.(string); ok {
		resp.Data = json.RawMessage(fmt.Sprintf(`"%s"`, payload))
	} else {
		data, _ := json.Marshal(payload)
		resp.Data = data
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshaling error response: %w", err)
	}
	err = msg.Respond(data)
	if err != nil {
		return fmt.Errorf("respondError: %w", err)
	}
	slog.Info("responded with error", "message", resp.Message, "code", code)
	return nil
}

func respondSuccess(msg *nats.Msg, message string, payload interface{}) error {
	resp := dto.GatewayResponse{
		Success: true,
		Message: message,
	}

	// Если payload уже []byte, используем его напрямую
	if data, ok := payload.([]byte); ok {
		resp.Data = data
	} else {
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshaling success response payload: %w", err)
		}
		resp.Data = data
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshaling success response: %w", err)
	}
	err = msg.Respond(data)
	if err != nil {
		return fmt.Errorf("respondSuccess: %w", err)
	}
	slog.Info("responded with success", "message", resp.Message)
	return nil
}
