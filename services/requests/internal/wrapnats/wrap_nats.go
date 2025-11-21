package wrapnats

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"requests/internal/domain"
)

type iRequestService interface {
	Add(ctx context.Context, newRequest domain.Request) (*domain.Request, error)
	UpdateStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error
	Get(ctx context.Context, requestID uuid.UUID) (*domain.Request, error)
	List(ctx context.Context, user domain.User, limit, offset int32) ([]domain.Request, error)
	Cancel(ctx context.Context, requestID uuid.UUID) error
}

type iEquipmentService interface {
	Add(ctx context.Context, technics []domain.Equipment) error
}

type Config struct {
	URL string `yaml:"url" env:"URL"`

	RequestTimeout time.Duration `yaml:"request_timeout" env:"REQUEST_TIMEOUT" env-default:"5s"`
}

type NatsWrapper struct {
	conn      *nats.Conn
	cfg       *Config
	validator *validator.Validate

	reqService iRequestService
	eqService  iEquipmentService
}

func New(cfg *Config, conn *nats.Conn, reqService iRequestService, eqService iEquipmentService) (*NatsWrapper, error) {
	return &NatsWrapper{
		conn:       conn,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
		cfg:        cfg,
		reqService: reqService,
		eqService:  eqService,
	}, nil
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]nats.MsgHandler{
		"requests.create":        w.handleCreateRequest,
		"requests.status.update": w.handleUpdateStatus,
		"requests.get":           w.handleGetRequest,
		"requests.list":          w.handleListRequests,
		"requests.cancel":        w.handleCancelRequest,
		"equipment.add":          w.handleAddEquipment,
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
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestCreate
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	domainReq := req.ToDomain()
	createdReq, err := w.reqService.Add(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respData, err := createdReq.MarshalJSON()
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleUpdateStatus(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestUpdateStatus
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update status request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(
			ctx,
			"validation error",
			"error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err = w.reqService.UpdateStatus(ctx, req.RequestID, domain.StatusType(req.Status))
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(
			ctx,
			"error updating request status",
			"error", err,
		)
		return
	}

	err = respondSuccess(msg, "success", "status updated")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetRequest(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestByID
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	reqDomain, err := w.reqService.Get(ctx, req.RequestID)
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.Error("error getting request", "error", err)
		return
	}

	err = respondSuccess(msg, "success", reqDomain)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleListRequests(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestList
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling list request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	user := domain.User{TelegramID: req.TelegramID}
	requests, err := w.reqService.List(ctx, user, req.Limit, req.Offset)
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(ctx, "error listing requests", "error", err)
		return
	}

	err = respondSuccess(msg, "success", requests)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleCancelRequest(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestByID
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling cancel request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err = w.reqService.Cancel(ctx, req.RequestID)
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(ctx, "error canceling request", "error", err)
		return
	}

	err = respondSuccess(msg, "success", "request canceled")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleAddEquipment(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req addEquipment
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling add equipment request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err = w.eqService.Add(ctx, req.Equipments)
	if err != nil {
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(ctx, "error adding equipment", "error", err)
		return
	}

	err = respondSuccess(msg, "success", "equipment added")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}
