package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"config"
	"observability/opentelemetry"
	"proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	"auth/internal/postgres"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"auth/internal/storeredis"
)

type Config struct {
	Repository repository.Config `yaml:"repository" env:"REPOSITORY"`
	Redis      storeredis.Config `yaml:"redis" env:"REDIS"`
	Postgres   postgres.Config   `yaml:"postgres" env:"POSTGRES"`

	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET"`
	Port      string `yaml:"port" env:"PORT"`
}

func main() {
	cfg := config.MustParseConfig[Config]("./cmd/config/cfg.yaml")

	ctx := context.Background()

	traceExp, err := opentelemetry.NewOTLPHTTPExporter(ctx, "", true)
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
	redisClient, err := storeredis.New(cfg.Redis, tracer)
	if err != nil {
		panic(err)
	}

	repo, err := repository.New(cfg.Repository, redisClient, pg)
	if err != nil {
		panic(err)
	}

	gen := service.NewTokenizer(cfg.Repository.AccessTokenTTL, cfg.Repository.RefreshTokenTTL, cfg.JWTSecret)

	svc := service.NewAuth(gen, repo)
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
