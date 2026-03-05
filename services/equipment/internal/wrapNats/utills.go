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
	statusBadRequest          statusCode = 400
	statusInternalServerError statusCode = 500
)

// respondError отвечает с ошибкой на нативный *nats.Msg (доступен через broker.Msg.Msg).
func respondError(msg *nats.Msg, message string, payload interface{}, code statusCode) error {
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
	if err := msg.Respond(data); err != nil {
		return fmt.Errorf("respondError: %w", err)
	}
	slog.Info("responded with error", "message", resp.Message, "code", code)
	return nil
}

// respondSuccess отвечает с успехом на нативный *nats.Msg (доступен через broker.Msg.Msg).
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
	if err := msg.Respond(data); err != nil {
		return fmt.Errorf("respondSuccess: %w", err)
	}
	slog.Info("responded with success", "message", resp.Message)
	return nil
}
