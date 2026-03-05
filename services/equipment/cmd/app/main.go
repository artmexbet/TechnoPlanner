package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/artmexbet/TechnoPlanner/libs/observability/opentelemetry"

	"tech/internal/postgres"
	categoryrepo "tech/internal/repository/category"
	techrepo "tech/internal/repository/tech"
	categoryservice "tech/internal/services/category"
	techservice "tech/internal/services/tech"
	wrapnats "tech/internal/wrapNats"
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

	// Трейсинг
	traceExp, err := opentelemetry.NewOTLPHTTPExporter(ctx, cfg.Trace.Endpoint, cfg.Trace.Insecure)
	if err != nil {
		panic(err)
	}
	tracer, shutdownFn := opentelemetry.NewTracerProvider(traceExp, "equipment-service")
	defer func() {
		if err := shutdownFn(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetTracerProvider(tracer)
	otel.SetTextMapPropagator(opentelemetry.NewPropagator())

	// PostgreSQL
	pgCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN())
	if err != nil {
		panic(err)
	}
	pgCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	// NATS через libs/broker
	conn, err := broker.Connect(cfg.Nats.URL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Dependency injection
	pg := postgres.New(pool)

	techRepo := techrepo.NewRepository(pg)
	catRepo := categoryrepo.NewRepository(pg)

	eqSvc := techservice.New(techRepo)
	catSvc := categoryservice.New(catRepo)

	wrapper, err := wrapnats.New(cfg.Nats, conn, eqSvc, catSvc)
	if err != nil {
		panic(err)
	}
	wrapper.HandleMsgs()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}
