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

func respondError(msg *nats.Msg, message string, payload interface{}, statusCode statusCode) error { //nolint:revive //todo: удалить статус код, он не нужен, так как в GatewayResponse уже есть поле Message для описания ошибки
	resp := dto.GatewayResponse{
		Success: false,
		Message: message,
	}

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
	if err = msg.Respond(data); err != nil {
		return fmt.Errorf("respondError: %w", err)
	}
	slog.Info("responded with error", "message", resp.Message)
	return nil
}

func respondSuccess(msg *nats.Msg, message string, payload interface{}) error {
	resp := dto.GatewayResponse{
		Success: true,
		Message: message,
	}

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
	if err = msg.Respond(data); err != nil {
		return fmt.Errorf("respondSuccess: %w", err)
	}
	slog.Info("responded with success", "message", resp.Message)
	return nil
}
