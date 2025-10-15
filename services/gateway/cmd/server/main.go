package main

import (
	"context"

	"broker"
	"config"

	"gateway/internal/app"
	"gateway/internal/postgres"
	"gateway/internal/service"
	"gateway/internal/storage"
)

type Config struct {
	Router   app.Config      `yaml:"router" env:"ROUTER"`
	Broker   broker.Config   `yaml:"broker" env:"BROKER"`
	Postgres postgres.Config `yaml:"postgres" env:"POSTGRES"`
}

func main() {
	cfg := config.MustParseConfig[Config]("config/cfg.yaml")

	ctx := context.Background()

	//nats := broker.NewNATSBroker(cfg.Broker) // понадобится позже

	_postgres, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer _postgres.Close()

	store := storage.NewStorage(_postgres) // инициализация хранилища

	userSvc := service.NewUserService(store)
	authSvc := service.NewAuthService(userSvc)

	r := app.NewRouter(cfg.Router, userSvc, authSvc).
		InitMiddlewares().
		InitBaseRoutes()
	r.Run()
}
