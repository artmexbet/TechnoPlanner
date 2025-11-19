package wrapnats

import (
	"context"
	"fmt"
	"log/slog"

	"requests/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type iRequestService interface {
	Add(ctx context.Context, newRequest domain.Request) (*domain.Request, error)
	UpdateStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error
	Get(ctx context.Context, requestID uuid.UUID) (*domain.Request, error)
	List(ctx context.Context, user domain.User, limit, offset int32) ([]domain.Request, error)
	Cancel(ctx context.Context, requestID uuid.UUID) error
}

type Config struct {
	URL string `yaml:"url" env:"URL"`
}

type NatsWrapper struct {
	conn      *nats.Conn
	validator *validator.Validate

	reqService iRequestService
}

func New(cfg Config, reqService iRequestService) (*NatsWrapper, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats server: %v", err)
	}
	return &NatsWrapper{
		conn:       conn,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
		reqService: reqService,
	}, nil
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]nats.MsgHandler{
		"requests.create":        w.handleCreateRequest,
		"requests.status.update": w.handleUpdateStatus,
		"requests.get":           w.handleGetRequest,
		"requests.list":          w.handleListRequests,
		"requests.cancel":        w.handleCancelRequest,
	}

	for subject, handler := range handlers {
		_, err := w.conn.Subscribe(subject, handler)
		if err != nil {
			slog.Error("cannot subscribe to subject", "subject", subject, "error", err)
		}
	}
	return w
}

func (w *NatsWrapper) handleCreateRequest(msg *nats.Msg) {
	var req requestCreate
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.Error("error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.Error("validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	domainReq := req.ToDomain()
	createdReq, err := w.reqService.Add(context.Background(), domainReq)
	if err != nil {
		slog.Error("error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respData, err := createdReq.MarshalJSON()
	if err != nil {
		slog.Error("error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.Error("error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleUpdateStatus(msg *nats.Msg) {
	var req requestUpdateStatus
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.Error("error unmarshaling update status request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.Error(
			"validation error",
			"error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err = w.reqService.UpdateStatus(context.Background(), req.RequestID, domain.StatusType(req.Status))
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.Error(
			"error updating request status",
			"error", err,
		)
		return
	}

	err = respondSuccess(msg, "success", "status updated")
	if err != nil {
		slog.Error("error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetRequest(msg *nats.Msg) {
	var req requestByID
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.Error("error unmarshaling get request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.Error("validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	reqDomain, err := w.reqService.Get(context.Background(), req.RequestID)
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.Error("error getting request", "error", err)
		return
	}

	err = respondSuccess(msg, "success", reqDomain)
	if err != nil {
		slog.Error("error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleListRequests(msg *nats.Msg) {
	var req requestList
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.Error("error unmarshaling list request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.Error("validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	user := domain.User{TelegramID: req.TelegramID}
	requests, err := w.reqService.List(context.Background(), user, req.Limit, req.Offset)
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.Error("error listing requests", "error", err)
		return
	}

	err = respondSuccess(msg, "success", requests)
	if err != nil {
		slog.Error("error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleCancelRequest(msg *nats.Msg) {
	var req requestByID
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.Error("error unmarshaling cancel request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.Error("validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err = w.reqService.Cancel(context.Background(), req.RequestID)
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.Error("error canceling request", "error", err)
		return
	}

	err = respondSuccess(msg, "success", "request canceled")
	if err != nil {
		slog.Error("error sending response", "error", err)
	}
}
