package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"observability/opentelemetry"
	"proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	natsPublisher "auth/internal/nats-publisher"
	"auth/internal/postgres"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"auth/internal/storeredis"

	"config"
)

const (
	defaultConfigPath = "./cmd/config/cfg.yaml"
	configPathKey     = "CONFIG_PATH"
)

type Config struct {
	Repository repository.Config `yaml:"repository" env-prefix:"REPOSITORY_"`
	Redis      storeredis.Config `yaml:"redis" env-prefix:"REDIS_"`
	Postgres   config.Postgres   `yaml:"postgres" env-prefix:"POSTGRES_"`
	Publisher  config.NATSConfig `yaml:"publisher" env-prefix:"PUBLISHER_"`

	Traces config.Trace `yaml:"trace" env-prefix:"TRACE_"`

	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET"`
	TokenCost int    `yaml:"token_cost" env:"TOKEN_COST" env-default:"8"`
	Port      string `yaml:"port" env:"PORT"`
}

func main() {
	cfgPath := os.Getenv(configPathKey)
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}
	cfg := config.MustParseConfig[Config](cfgPath)
	slog.Info("loaded config", "config", cfg)

	ctx := context.Background()

	traceExp, err := opentelemetry.NewOTLPHTTPExporter(ctx, cfg.Traces.Endpoint, cfg.Traces.Insecure)
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
	otel.SetTextMapPropagator(opentelemetry.NewPropagator())

	pg, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	redisClient, err := storeredis.New(cfg.Redis, tracer)
	if err != nil {
		panic(err)
	}
	defer redisClient.Close() //nolint:errcheck

	publisher, err := natsPublisher.NewPublisher(cfg.Publisher.Host, cfg.Publisher.Port)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	repo, err := repository.New(cfg.Repository, redisClient, pg, publisher)
	if err != nil {
		panic(err)
	}

	gen := service.NewTokenizer(cfg.Repository.AccessTokenTTL, cfg.Repository.RefreshTokenTTL, cfg.JWTSecret)

	svc := service.NewAuth(gen, repo, cfg.TokenCost)
	handler := server.NewHandler(svc)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(server.LoggingInterceptor),
	)
	proto.RegisterAuthServer(grpcServer, handler)
	slog.Info("Starting gRPC server", "port", cfg.Port)
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}
	if err = grpcServer.Serve(grpcListener); err != nil {
		panic(err) // TODO: handle error properly
	}
}
