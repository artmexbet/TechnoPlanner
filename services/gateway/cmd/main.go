package main

import (
	"context"
	"log/slog"

	"observability/opentelemetry"

	"broker"
	"config"
	"gateway/internal/postgres"
	"gateway/internal/router"
	"gateway/internal/service"
	"gateway/internal/storage"

	"go.opentelemetry.io/otel"
)

type Config struct {
	Router   router.Config             `yaml:"router" env:"ROUTER"`
	Broker   broker.Config             `yaml:"broker" env:"BROKER"`
	Postgres postgres.Config           `yaml:"postgres" env:"POSTGRES"`
	GRPC     service.AuthServiceConfig `yaml:"grpc" env:"GRPC"`
}

func main() {
	cfg := config.MustParseConfig[Config]("./cmd/config/cfg.yaml")

	ctx := context.Background()

	// Трейсы мои трейсы
	exporter, err := opentelemetry.NewOTLPHTTPExporter(ctx, "", true)
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

	store := storage.NewStorage(_postgres) // инициализация хранилища

	userSvc := service.NewUserService(store)
	authSvc, err := service.NewGRPCWrapper(cfg.GRPC)
	if err != nil {
		panic(err)
	}
	defer authSvc.Close() //nolint:errcheck

	r := router.NewRouter(cfg.Router, userSvc, authSvc).
		InitMiddlewares(tracer).
		InitBaseRoutes().
		InitUserRoutes().
		InitProtectedUserRoutes()
	slog.Info("Starting HTTP server")
	r.Run()
}
