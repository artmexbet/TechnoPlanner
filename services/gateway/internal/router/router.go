package router

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	slogfiber "github.com/samber/slog-fiber"
	"go.opentelemetry.io/otel/trace"

	"gateway/internal/models"
)

type iUserService interface {
	// Define methods for user service operations
}

type iAuthSvcConnector interface {
	Login(ctx context.Context, req models.LoginRequest) (models.TokenPair, error)
	Register(ctx context.Context, username, password, email string) (string, error)
	ValidateToken(ctx context.Context, token string) (models.TokenValidationResponse, error)
	Refresh(ctx context.Context, req models.TokenRefreshRequest) (models.TokenPair, error)
	Logout(ctx context.Context, token string) error
	LogoutAll(ctx context.Context, token string) error
}

type Config struct {
	Address string `yaml:"address" env:"ADDRESS"`
	Port    string `yaml:"port" env:"PORT"`
}

type Router struct {
	r         *fiber.App
	validator *validator.Validate

	cfg     Config
	userSvc iUserService
	authSvc iAuthSvcConnector
}

func NewRouter(cfg Config, userSvc iUserService, authSvc iAuthSvcConnector) *Router {
	return &Router{
		r:         fiber.New(),
		validator: validator.New(validator.WithRequiredStructEnabled()),
		cfg:       cfg,
		userSvc:   userSvc,
		authSvc:   authSvc,
	}
}

func (r *Router) InitMiddlewares(provider trace.TracerProvider) *Router {
	r.r.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
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
	r.r.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Простая проверка доступности API
	r.r.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	return r
}

func (r *Router) Run() {
	if err := r.r.Listen(fmt.Sprintf("%s:%s", r.cfg.Address, r.cfg.Port)); err != nil {
		panic(err)
	}
}
