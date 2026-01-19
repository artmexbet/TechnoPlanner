package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/broker/middleware"
	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/artmexbet/TechnoPlanner/libs/observability/opentelemetry"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/postgres"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/repository"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/service/equipment"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/service/request"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/wrapnats"
)

const (
	defaultConfigPath = "./cmd/config/config.yaml"
	configPathKey     = "CONFIG_PATH"
)

type Config struct {
	Nats     *wrapnats.Config `yaml:"nats" env-prefix:"NATS_"`
	Postgres *config.Postgres `yaml:"postgres" env-prefix:"POSTGRES_"`
	Trace    config.Trace     `yaml:"trace" env-prefix:"TRACE_"`
}

func main() {
	cfgPath := os.Getenv(configPathKey)
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}
	cfg := config.MustParseConfig[Config](cfgPath)
	ctx := context.Background()

	traceExp, err := opentelemetry.NewOTLPHTTPExporter(ctx, cfg.Trace.Endpoint, cfg.Trace.Insecure)
	if err != nil {
		panic(err)
	}
	tracer, shutdownFn := opentelemetry.NewTracerProvider(traceExp, "auth-service")
	defer func() {
		if err := shutdownFn(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetTracerProvider(tracer)

	// Настраиваем propagator для передачи trace context между сервисами
	otel.SetTextMapPropagator(opentelemetry.NewPropagator()) //todo: use traces later

	conn, err := broker.Connect(cfg.Nats.URL, nats.Name("Requests Service"))
	if err != nil {
		panic(err)
	}

	// Apply middlewares
	conn.Use(middleware.NewLoggingMiddleware(true))
	conn.Use(middleware.NewRecoveryMiddleware())
	conn.Use(middleware.NewRequestIDMiddleware())
	conn.Use(middleware.NewTimeoutMiddleware(cfg.Nats.RequestTimeout))

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		panic(err)
	}

	publisher := wrapnats.NewNatsPublisher(conn)

	pg := postgres.New(pool)
	repo := repository.NewRepository(pg, publisher)

	requestService := request.New(repo)
	equipmentService := equipment.New(repo)

	wrapper, err := wrapnats.New(cfg.Nats, conn, requestService, equipmentService)
	if err != nil {
		panic(err)
	}
	wrapper.HandleMsgs()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	conn.Close()
	pool.Close()
}
