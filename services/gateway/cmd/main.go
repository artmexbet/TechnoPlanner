package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/artmexbet/TechnoPlanner/libs/observability/opentelemetry"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/client"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/service"
)

const (
	defaultConfigPath = "./cmd/config/cfg.yaml"
	configPathKey     = "CONFIG_PATH"
)

type Config struct {
	Router router.Config             `yaml:"router" env:"ROUTER"`
	Broker config.NATSConfig         `yaml:"broker" env:"BROKER"`
	GRPC   service.AuthServiceConfig `yaml:"grpc" env:"GRPC"`
	Trace  config.Trace              `yaml:"trace" env:"TRACE"`
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

	// Подключаемся к NATS
	natsConn, err := broker.Connect(cfg.Broker.URL(), nats.Name("Gateway Service"))
	if err != nil {
		panic(err)
	}
	defer natsConn.Close()
	slog.Info("Connected to NATS")

	// Создаем клиенты для сервисов через NATS Request-Reply
	requestClient := client.NewRequestClient(natsConn)

	authSvc, err := service.NewGRPCWrapper(cfg.GRPC)
	if err != nil {
		panic(err)
	}
	defer authSvc.Close() //nolint:errcheck

	// Создаем клиенты для сервисов через NATS Request-Reply
	userClient := client.NewUserClient(natsConn)
	porterClient := client.NewPorterClient(natsConn)
	porterSvc := service.NewPorterService(porterClient, userClient, authSvc)
	equipmentClient := client.NewEquipmentClient(natsConn)
	equipmentSvc := service.NewEquipmentService(equipmentClient)
	categoryClient := client.NewCategoryClient(natsConn)
	categorySvc := service.NewCategoryService(categoryClient)
	requestSvc := service.NewRequestService(requestClient)
	historyClient := client.NewHistoryClient(natsConn)
	historySvc := service.NewRequestHistoryService(historyClient)
	rawRequestSvc := service.NewRawRequestService(requestClient)

	r := router.NewRouter(cfg.Router, nil, authSvc, porterSvc, equipmentSvc, categorySvc, requestSvc, historySvc, rawRequestSvc).
		InitMiddlewares(tracer).
		InitBaseRoutes().
		InitUserRoutes().
		InitProtectedUserRoutes().
		InitPorterRoutes().
		InitEquipmentRoutes().
		InitRequestRoutes().
		InitRawRequestRoutes()
	slog.Info("Starting HTTP server")
	r.Run()
}
