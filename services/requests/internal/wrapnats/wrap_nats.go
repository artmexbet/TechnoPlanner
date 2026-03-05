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
	UpdateRequest(ctx context.Context, requestID uuid.UUID, updates domain.RequestUpdate) (*domain.Request, error)
	ListResponsibles(ctx context.Context) ([]domain.Porter, error)
	SaveResponsible(ctx context.Context, id uuid.UUID, username string) error
	GetResponsible(ctx context.Context, id uuid.UUID) (domain.Porter, error)
	DeleteResponsible(ctx context.Context, id uuid.UUID) error
	// Raw request methods
	CreateRawRequest(ctx context.Context, req domain.RawRequest) (*domain.RawRequest, error)
	ListRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error)
	GetRawRequest(ctx context.Context, id uuid.UUID) (*domain.RawRequest, error)
	ProcessRawRequest(ctx context.Context, rawID uuid.UUID, newRequest domain.Request) (*domain.Request, *domain.RawRequest, error)
}

type EquipmentService interface {
	Add(ctx context.Context, technics []domain.Equipment) error
	SyncCreate(ctx context.Context, eq domain.Equipment) error
	SyncUpdate(ctx context.Context, eq domain.Equipment) error
	SyncDelete(ctx context.Context, id int) error
}

// ResponsibleStorage интерфейс для работы с портерами при событии UserCreated
type PorterStorage interface {
	SaveResponsible(ctx context.Context, id uuid.UUID, username string) error
}

type Config struct {
	URL            string        `yaml:"url" env:"URL"`
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
		// Bot handlers (все запросы от бота теперь сырые)
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
		subjects.GatewayRequestUpdate:            w.handleGatewayUpdateRequest,
		// Porter handlers (хранение портеров в Requests)
		subjects.GatewayPorterList:   w.handleGatewayListResponsibles,
		subjects.GatewayPorterGet:    w.handleGatewayGetResponsible,
		subjects.GatewayPorterDelete: w.handleGatewayDeleteResponsible,
		subjects.GatewayPorterSave:   w.handleGatewayCreateResponsible,
		// Event handlers
		subjects.UserCreated: w.handleUserCreated,
		// Equipment sync events (pub/sub от Equipment Service)
		subjects.EquipmentCreated: w.handleEquipmentCreatedEvent,
		subjects.EquipmentUpdated: w.handleEquipmentUpdatedEvent,
		subjects.EquipmentDeleted: w.handleEquipmentDeletedEvent,
		// Raw request handlers (для gateway)
		subjects.GatewayRawRequestList:    w.handleGatewayListRawRequests,
		subjects.GatewayRawRequestGet:     w.handleGatewayGetRawRequest,
		subjects.GatewayRawRequestProcess: w.handleGatewayProcessRawRequest,
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

// handleUserCreated обрабатывает событие создания пользователя от auth сервиса
// Сохраняет пользователя как портера только если у него роль porter (roleID=2)
func (w *NatsWrapper) handleUserCreated(msg *broker.Msg) error {
	ctx := msg.Context()

	var event dto.UserCreatedEvent
	if err := event.UnmarshalJSON(msg.Data); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling user created event", "error", err)
		return err
	}

	// Сохраняем только портеров (role_id = 2)
	if event.RoleID != 2 {
		slog.InfoContext(ctx, "skipping non-porter user", "user_id", event.ID, "role_id", event.RoleID)
		return nil
	}

	if err := w.reqService.SaveResponsible(ctx, event.ID, event.Username); err != nil {
		slog.ErrorContext(ctx, "error saving porter", "error", err, "user_id", event.ID, "username", event.Username)
		return err
	}

	slog.InfoContext(ctx, "porter saved", "user_id", event.ID, "username", event.Username)
	return nil
}
