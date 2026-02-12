package wrapnats

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type RequestService interface {
	Add(ctx context.Context, newRequest domain.Request) (*domain.Request, error)
	UpdateStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error
	Get(ctx context.Context, requestID uuid.UUID) (*domain.Request, error)
	List(ctx context.Context, user domain.User, limit, offset int32) ([]domain.Request, error)
	Cancel(ctx context.Context, requestID uuid.UUID) error
	// Gateway methods

	ListByResponsible(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error)
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID) (*domain.Request, error)
}

type EquipmentService interface {
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

	reqService RequestService
	eqService  EquipmentService
}

func New(cfg *Config, conn *broker.NATS, reqService RequestService, eqService EquipmentService) (*NatsWrapper, error) {
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
		// Bot handlers
		subjects.SubjectBotRequestCreated: w.handleCreateRequest,
		subjects.RequestStatusChanged:     w.handleUpdateStatus,
		subjects.SubjectBotRequestGet:     w.handleGetRequest,
		subjects.SubjectBotRequestList:    w.handleListRequests,
		subjects.SubjectBotRequestCancel:  w.handleCancelRequest,
		subjects.GatewayRequestCanceled:   w.handleCancelRequestFromService,
		subjects.ServiceEquipmentAdd:      w.handleAddEquipment,
		// Gateway handlers
		subjects.GatewayRequestList:              w.handleGatewayListRequests,
		subjects.GatewayRequestGet:               w.handleGatewayGetRequest,
		subjects.GatewayRequestAssignResponsible: w.handleGatewayAssignResponsible,
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

// Bot handlers

func (w *NatsWrapper) handleCreateRequest(msg *broker.Msg) error {
	ctx := msg.Context()
	var req dto.RequestCreateRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	domainReq := mapRequestCreateToDomain(req)
	createdReq, err := w.reqService.Add(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapRequestToDTO(*createdReq)
	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", respData)
}

func (w *NatsWrapper) handleUpdateStatus(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestStatusUpdateRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update status request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	err = w.reqService.UpdateStatus(ctx, req.RequestID, domain.StatusType(req.Status))
	if err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", "status updated")
}

func (w *NatsWrapper) handleGetRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	reqDomain, err := w.reqService.Get(ctx, req.RequestID)
	if err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapRequestToDTO(*reqDomain)
	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleListRequests(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestListByTelegramRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling list request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	user := domain.User{TelegramID: req.TelegramID}
	requests, err := w.reqService.List(ctx, user, req.Limit, req.Offset)
	if err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make([]dto.Request, len(requests))
	for i, r := range requests {
		respDTO[i] = mapRequestToDTO(r)
	}

	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleCancelRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling cancel request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	err = w.reqService.Cancel(ctx, req.RequestID)
	if err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", "request canceled")
}

func (w *NatsWrapper) handleCancelRequestFromService(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling cancel request from service", "error", err)
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

	var req dto.EquipmentAddRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling add equipment request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	domainEquipment := make([]domain.Equipment, len(req.Equipments))
	for i, eq := range req.Equipments {
		domainEquipment[i] = domain.Equipment{
			ID:       eq.ID,
			Name:     eq.Name,
			Quantity: eq.Quantity,
		}
	}

	err = w.eqService.Add(ctx, domainEquipment)
	if err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", "equipment added")
}

// Gateway handlers

func (w *NatsWrapper) handleGatewayListRequests(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestListRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling gateway list request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	requests, err := w.reqService.ListByResponsible(ctx, req.ResponsibleID)
	if err != nil {
		slog.ErrorContext(ctx, "error listing requests", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make([]dto.Request, len(requests))
	for i, r := range requests {
		respDTO[i] = mapRequestToDTO(r)
	}

	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleGatewayGetRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling gateway get request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	reqDomain, err := w.reqService.Get(ctx, req.RequestID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting request", "error", err)
		return respondError(msg.Msg, "not found", "request not found", statusNotFound)
	}

	respDTO := mapRequestToDTO(*reqDomain)
	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleGatewayAssignResponsible(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.AssignResponsibleRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling gateway assign request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	updatedReq, err := w.reqService.AssignResponsible(ctx, req.RequestID, req.ResponsibleID)
	if err != nil {
		slog.ErrorContext(ctx, "error assigning responsible", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapRequestToDTO(*updatedReq)
	return respondSuccess(msg.Msg, "success", respDTO)
}
