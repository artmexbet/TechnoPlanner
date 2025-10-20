package main

import (
	"context"
	"log"

	"broker"
	"config"
	"telemetry"

	"gateway/internal/app"
	"gateway/internal/postgres"
	"gateway/internal/service"
	"gateway/internal/storage"
)

type Config struct {
	Router    app.Config                `yaml:"router" env:"ROUTER"`
	Broker    broker.Config             `yaml:"broker" env:"BROKER"`
	Postgres  postgres.Config           `yaml:"postgres" env:"POSTGRES"`
	GRPC      service.AuthServiceConfig `yaml:"grpc" env:"GRPC"`
	Telemetry telemetry.Config          `yaml:"telemetry" env:"TELEMETRY"`
}

func main() {
	cfg := config.MustParseConfig[Config]("config/cfg.yaml")

	ctx := context.Background()

	// Initialize telemetry
	tel, err := telemetry.New(cfg.Telemetry)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tel.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down telemetry: %v", err)
		}
	}()

	//nats := broker.NewNATSBroker(cfg.Broker) // понадобится позже

	_postgres, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer _postgres.Close()

	store := storage.NewStorage(_postgres) // инициализация хранилища

	userSvc := service.NewUserService(store)
	authSvc, err := service.NewGRPCWrapper(cfg.GRPC)
	if err != nil {
		panic(err)
	}
	defer authSvc.Close()

	// Pass service name to router for telemetry middleware
	cfg.Router.ServiceName = cfg.Telemetry.ServiceName

	r := app.NewRouter(cfg.Router, userSvc, authSvc).
		InitMiddlewares().
		InitBaseRoutes().
		InitUserRoutes()
	r.Run()
}
