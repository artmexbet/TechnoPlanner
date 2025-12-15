package wrapnats

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
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
	conn          *broker.NATS
	cfg           *Config
	validator     *validator.Validate
	subscriptions []*nats.Subscription

	reqService iRequestService
	eqService  iEquipmentService
}

func New(cfg *Config, conn *broker.NATS, reqService iRequestService, eqService iEquipmentService) (*NatsWrapper, error) {
	return &NatsWrapper{
		conn:       conn,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
		cfg:        cfg,
		reqService: reqService,
		eqService:  eqService,
	}, nil
}

func (w *NatsWrapper) Close() {
	for _, sub := range w.subscriptions {
		_ = sub.Unsubscribe()
	}
	w.conn.Close()
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]broker.MsgHandler{
		subjects.SubjectBotRequestCreated: w.handleCreateRequest,
		subjects.RequestStatusChanged:     w.handleUpdateStatus,
		subjects.SubjectBotRequestGet:     w.handleGetRequest,
		subjects.SubjectBotRequestList:    w.handleListRequests,
		subjects.SubjectBotRequestCancel:  w.handleCancelRequest,
		subjects.GatewayRequestCanceled:   w.handleCancelRequestFromService,
		subjects.ServiceEquipmentAdd:      w.handleAddEquipment,
	}

	for subject, handler := range handlers {
		sub, err := w.conn.Subscribe(subject, handler)
		if err != nil {
			slog.Error("cannot subscribe to subject", "subject", subject, "error", err)
		}
		w.subscriptions = append(w.subscriptions, sub)
	}
	return w
}

func (w *NatsWrapper) handleCreateRequest(msg *broker.Msg) error {
	ctx := msg.Context()
	var req requestCreate
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	domainReq := req.ToDomain()
	createdReq, err := w.reqService.Add(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respData, err := createdReq.MarshalJSON()
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	err = respondSuccess(msg.Msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
		return err
	}
	return nil
}

func (w *NatsWrapper) handleUpdateStatus(msg *broker.Msg) error {
	ctx := msg.Context()

	var req requestUpdateStatus
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update status request", "error", err)
		_ = respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
		return err
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(
			ctx,
			"validation error",
			"error", err)
		_ = respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
		return err
	}

	err = w.reqService.UpdateStatus(ctx, req.RequestID, domain.StatusType(req.Status))
	if err != nil {
		_ = respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(
			ctx,
			"error updating request status",
			"error", err,
		)
		return err
	}

	err = respondSuccess(msg.Msg, "success", "status updated")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
		return err
	}
	return nil
}

func (w *NatsWrapper) handleGetRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req requestByID
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get request", "error", err)
		_ = respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
		return err
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
		return err
	}

	reqDomain, err := w.reqService.Get(ctx, req.RequestID)
	if err != nil {
		_ = respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
		slog.Error("error getting request", "error", err)
		return err
	}

	err = respondSuccess(msg.Msg, "success", reqDomain)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
		return err
	}
	return nil
}

func (w *NatsWrapper) handleListRequests(msg *broker.Msg) error {
	ctx := msg.Context()

	var req requestList
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling list request", "error", err)
		_ = respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
		return err
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
		return err
	}

	user := domain.User{TelegramID: req.TelegramID}
	requests, err := w.reqService.List(ctx, user, req.Limit, req.Offset)
	if err != nil {
		_ = respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(ctx, "error listing requests", "error", err)
		return err
	}

	err = respondSuccess(msg.Msg, "success", requests)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
		return err
	}
	return nil
}

func (w *NatsWrapper) handleCancelRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req requestByID
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling cancel request", "error", err)
		_ = respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
		return err
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
		return err
	}

	err = w.reqService.Cancel(ctx, req.RequestID)
	if err != nil {
		_ = respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(ctx, "error canceling request", "error", err)
		return err
	}

	err = respondSuccess(msg.Msg, "success", "request canceled")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
		return err
	}
	return nil
}

func (w *NatsWrapper) handleCancelRequestFromService(msg *broker.Msg) error {
	ctx := msg.Context()

	var req requestByID
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling cancel request from service", "error", err)
		return err
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return err
	}

	err = w.reqService.Cancel(ctx, req.RequestID)
	if err != nil {
		slog.ErrorContext(ctx, "error canceling request from service", "error", err)
		return err
	}

	slog.InfoContext(ctx, "request canceled from service", "request_id", req.RequestID)
	return nil
}

func (w *NatsWrapper) handleAddEquipment(msg *broker.Msg) error {
	ctx := msg.Context()

	var req addEquipment
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling add equipment request", "error", err)
		_ = respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
		return err
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
		return err
	}

	err = w.eqService.Add(ctx, req.Equipments)
	if err != nil {
		_ = respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
		slog.ErrorContext(ctx, "error adding equipment", "error", err)
		return err
	}

	err = respondSuccess(msg.Msg, "success", "equipment added")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
		return err
	}
	return nil
}
