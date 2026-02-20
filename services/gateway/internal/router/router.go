package router

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	otelfiber "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
	slogfiber "github.com/samber/slog-fiber"
	"go.opentelemetry.io/otel/trace"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/service"
)

type UserService interface {
	// Define methods for user service operations
}

type AuthSvcConnector interface {
	Login(ctx context.Context, req models.LoginRequest) (models.TokenPair, error)
	Register(ctx context.Context, username, password, email string) (string, error)
	ValidateToken(ctx context.Context, token string) (models.TokenValidationResponse, error)
	Refresh(ctx context.Context, req models.TokenRefreshRequest) (models.TokenPair, error)
	Logout(ctx context.Context, token string) error
	LogoutAll(ctx context.Context, token string) error
}

type PorterService interface {
	List(ctx context.Context) ([]domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetCurrentUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	Create(ctx context.Context, username, email, password string) (string, error)
}

type EquipmentService interface {
	Create(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	Update(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	Get(ctx context.Context, id int) (domain.Equipment, error)
	List(ctx context.Context) ([]domain.Equipment, error)
	Delete(ctx context.Context, id int) error
}

type CategoryService interface {
	Create(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	List(ctx context.Context) ([]domain.EquipmentCategory, error)
	Delete(ctx context.Context, id int) error
}

type RequestService interface {
	List(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Request, error)
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID uuid.UUID) (domain.Request, error)
	UpdateRequest(ctx context.Context, requestID uuid.UUID, updates domain.RequestUpdate) (domain.Request, error)
}

type HistoryService interface {
	List(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error)
	Add(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error)
}

type ResponsibleService interface {
	List(ctx context.Context) ([]domain.Responsible, error)
	Create(ctx context.Context, id uuid.UUID, username string) (domain.Responsible, error)
}

type RawRequestService interface {
	List(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error)
	Get(ctx context.Context, id uuid.UUID) (domain.RawRequest, error)
	Process(ctx context.Context, rawID uuid.UUID, req models.RawRequestProcessRequest) (domain.Request, domain.RawRequest, error)
}

type Config struct {
	Address string `yaml:"address" env:"ADDRESS"`
	Port    string `yaml:"port" env:"PORT"`
}

type Router struct {
	r         *fiber.App
	validator *validator.Validate

	cfg            Config
	userSvc        UserService
	authSvc        AuthSvcConnector
	porterSvc      PorterService
	equipmentSvc   EquipmentService
	categorySvc    CategoryService
	requestSvc     RequestService
	historySvc     HistoryService
	responsibleSvc ResponsibleService
	rawRequestSvc  RawRequestService
}

func NewRouter(
	cfg Config,
	userSvc UserService,
	authSvc AuthSvcConnector,
	porterSvc PorterService,
	equipmentSvc EquipmentService,
	categorySvc CategoryService,
	requestSvc RequestService,
	historySvc HistoryService,
	responsibleSvc ResponsibleService,
	rawRequestSvc RawRequestService,
) *Router {
	return &Router{
		r:              fiber.New(),
		validator:      validator.New(validator.WithRequiredStructEnabled()),
		cfg:            cfg,
		userSvc:        userSvc,
		authSvc:        authSvc,
		porterSvc:      porterSvc,
		equipmentSvc:   equipmentSvc,
		categorySvc:    categorySvc,
		requestSvc:     requestSvc,
		historySvc:     historySvc,
		responsibleSvc: responsibleSvc,
		rawRequestSvc:  rawRequestSvc,
	}
}

func (r *Router) InitMiddlewares(provider trace.TracerProvider) *Router {
	r.r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
		},
	))
	r.r.Use(recover.New())
	r.r.Use(requestid.New()) // Trace
	r.r.Use(
		otelfiber.Middleware(
			otelfiber.WithTracerProvider(provider),
		),
	)
	r.r.Use(slogfiber.New(slog.Default()))
	return r
}

// InitBaseRoutes регистрирует HTTP-маршруты
func (r *Router) InitBaseRoutes() *Router {
	// Healthcheck
	r.r.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Простая проверка доступности API
	r.r.Get("/ping", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	return r
}

func (r *Router) Run() {
	if err := r.r.Listen(fmt.Sprintf("%s:%s", r.cfg.Address, r.cfg.Port)); err != nil {
		panic(err)
	}
}

func (r *Router) userContext(c fiber.Ctx) context.Context {
	userID, _ := c.Locals(middlwares.ContextUserIDKey).(string)
	role, _ := c.Locals(middlwares.ContextUserRoleKey).(string)
	return service.WithUserContext(c.Context(), userID, role)
}
