package main

import (
	"context"
	"log"
	"net"

	"auth/internal/postgres"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"auth/internal/storeredis"
	"proto"

	"config"
	"telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Config struct {
	Repository repository.Config `yaml:"repository" env:"REPOSITORY"`
	Redis      storeredis.Config `yaml:"redis" env:"REDIS"`
	Postgres   postgres.Config   `yaml:"postgres" env:"POSTGRES"`
	Telemetry  telemetry.Config  `yaml:"telemetry" env:"TELEMETRY"`

	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET"`
}

func main() {
	cfg := config.MustParseConfig[Config]("./cmd/config/cfg.yaml")

	// Initialize telemetry
	tel, err := telemetry.New(cfg.Telemetry)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tel.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down telemetry: %v", err)
		}
	}()

	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	redisClient, err := storeredis.New(cfg.Redis)
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

	// Create gRPC server with OpenTelemetry interceptors
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	proto.RegisterAuthServer(grpcServer, handler)
	grpcListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	if err = grpcServer.Serve(grpcListener); err != nil {
		panic(err) // TODO: handle error properly
	}
}
