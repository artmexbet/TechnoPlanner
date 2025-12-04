package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"

	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/artmexbet/TechnoPlanner/libs/observability/opentelemetry"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/postgres"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/service"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/storage"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/subscriber"
)

const (
	defaultConfigPath = "./cmd/config/cfg.yaml"
	configPathKey     = "CONFIG_PATH"
)

type Config struct {
	Router   router.Config             `yaml:"router" env:"ROUTER"`
	Broker   config.NATSConfig         `yaml:"broker" env:"BROKER"`
	Postgres postgres.Config           `yaml:"postgres" env:"POSTGRES"`
	GRPC     service.AuthServiceConfig `yaml:"grpc" env:"GRPC"`
	Trace    config.Trace              `yaml:"trace" env:"TRACE"`
}

func main() {
	cfgPath := os.Getenv(configPathKey)
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}
	cfg := config.MustParseConfig[Config](cfgPath)
	slog.Info("config", "config", cfg)

	ctx := context.Background()

	// Трейсы мои трейсы
	exporter, err := opentelemetry.NewOTLPHTTPExporter(ctx, cfg.Trace.Endpoint, cfg.Trace.Insecure)
	if err != nil {
		panic(err)
	}
	tracer, shutdownFn := opentelemetry.NewTracerProvider(exporter, "gateway-service")
	defer func() {
		if err := shutdownFn(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetTracerProvider(tracer)

	// Настраиваем propagator для передачи trace context между сервисами
	otel.SetTextMapPropagator(opentelemetry.NewPropagator())
	slog.Info("Starting otel connection")

	//nats := broker.NewNATSBroker(cfg.Broker) // понадобится позже

	_postgres, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer _postgres.Close()
	slog.Info("Starting postgres connection")

	store := storage.NewStorage(_postgres, nil) // TODO: inject publisher

	authSvc, err := service.NewGRPCWrapper(cfg.GRPC)
	if err != nil {
		panic(err)
	}
	defer authSvc.Close() //nolint:errcheck

	userSvc := service.NewUserService(store)
	porterSvc := service.NewPorterService(store.Porters, authSvc)
	equipmentSvc := service.NewEquipmentService(store.Equipment)
	categorySvc := service.NewCategoryService(store.Categories)
	requestSvc := service.NewRequestService(store.Requests)
	historySvc := service.NewRequestHistoryService(store.StatusHistory)

	slog.Info("Starting subscriber service")

	subs, err := subscriber.New(cfg.Broker, store.Porters)
	if err != nil {
		panic(err)
	}
	defer subs.Close()
	slog.Info("Subscriber connected to NATS")

	if err := subs.Init(); err != nil {
		panic(err)
	}

	r := router.NewRouter(cfg.Router, userSvc, authSvc, porterSvc, equipmentSvc, categorySvc, requestSvc, historySvc).
		InitMiddlewares(tracer).
		InitBaseRoutes().
		InitUserRoutes().
		InitProtectedUserRoutes().
		InitPorterRoutes().
		InitEquipmentRoutes().
		InitRequestRoutes()
	slog.Info("Starting HTTP server")
	r.Run()
}
