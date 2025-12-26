package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
)

func NewLoggingMiddleware(nackMsgOnFailure bool) broker.Middleware {
	return func(next broker.MsgHandler) broker.MsgHandler {
		return func(msg *broker.Msg) error {
			ctx := msg.Context()
			// Log before processing
			slog.InfoContext(ctx,
				"Processing message",
				"subject", msg.Subject,
				"data", string(msg.Data))

			err := next(msg)

			if err != nil {
				slog.ErrorContext(ctx, "Error processing message", "subject", msg.Subject, "error", err)

				if nackMsgOnFailure {
					errNack := msg.Nak()
					if errNack != nil {
						slog.ErrorContext(ctx, "Error sending NAK", "subject", msg.Subject, "error", errNack)
					}
				}
				return err
			}

			// Log after successful processing
			slog.InfoContext(ctx, "Successfully processed message", "subject", msg.Subject)
			return nil
		}
	}
}

func NewRecoveryMiddleware() broker.Middleware {
	return func(next broker.MsgHandler) broker.MsgHandler {
		return func(msg *broker.Msg) (err error) {
			defer func() {
				if r := recover(); r != nil {
					slog.ErrorContext(msg.Context(), "Recovered from panic", "subject", msg.Subject, "panic", r)
					err = msg.Nak()
					if err != nil {
						slog.ErrorContext(msg.Context(), "Error sending Term", "subject", msg.Subject, "error", err)
					}
				}
			}()
			return next(msg)
		}
	}
}

type RequestIDKeyType string

type RequestIDConfig struct {
	Generator    func(msg *broker.Msg) string
	RequestIDKey RequestIDKeyType
}

var defaultRequestIDConfig = RequestIDConfig{
	Generator: func(_ *broker.Msg) string {
		return uuid.NewString()
	},
	RequestIDKey: "request_id",
}

func NewRequestIDMiddleware(cfg ...RequestIDConfig) broker.Middleware {
	_cfg := defaultRequestIDConfig
	if len(cfg) != 0 {
		_cfg = cfg[0]
	}
	return func(next broker.MsgHandler) broker.MsgHandler {
		return func(msg *broker.Msg) error {
			id := _cfg.Generator(msg)
			msg.SetContext(context.WithValue(msg.Context(), _cfg.RequestIDKey, id))
			slog.Debug("Adding request ID to context", "subject", msg.Subject)
			return next(msg)
		}
	}
}

func NewTimeoutMiddleware(timeout time.Duration) broker.Middleware {
	return func(next broker.MsgHandler) broker.MsgHandler {
		return func(msg *broker.Msg) error {
			ctx, cancel := context.WithTimeout(msg.Context(), timeout)
			defer cancel()
			msg.SetContext(ctx)
			return next(msg)
		}
	}
}
