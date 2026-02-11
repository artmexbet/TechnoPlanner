package wrapnats

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

func respondError(msg *nats.Msg, message string, payload interface{}, statusCode statusCode) error {
	resp := newResponse(message, statusCode, payload, true)
	data, err := resp.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshaling error response: %w", err)
	}
	err = msg.Respond(data)
	if err != nil {
		return fmt.Errorf("respondError: %w", err)
	}
	slog.Info("responded with error", "message", resp.Message)
	return nil
}

func respondSuccess(msg *nats.Msg, message string, payload interface{}) error {
	resp := newResponse(message, statusOK, payload, false)
	data, err := resp.MarshalJSON()
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
