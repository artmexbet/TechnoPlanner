package main

import (
	"context"

	"technoBro/internal/app/api"
	"technoBro/internal/app/api/service"
	"technoBro/internal/broker"
	"technoBro/internal/config"
	"technoBro/internal/postgres"
	"technoBro/internal/storage"
)

type Config struct {
	Router   api.Config      `yaml:"router" env:"ROUTER"`
	Broker   broker.Config   `yaml:"broker" env:"BROKER"`
	Postgres postgres.Config `yaml:"postgres" env:"POSTGRES"`
}

func main() {
	cfg := config.MustParseConfig[Config]("config/cfg.yaml")

	ctx := context.Background()

	nats := broker.NewNATSBroker(cfg.Broker).
		WithSomething(ctx)

	_postgres, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer _postgres.Close()

	store := storage.NewStorage(_postgres) // инициализация хранилища
	_ = store                              // заглушка, чтобы не было ошибки о неиспользуемой переменной

	techSvc := service.NewTechManager(nats)
	taskSvc := service.NewTaskManager(nats)

	r := api.NewRouter(cfg.Router, techSvc, taskSvc)
	r.InitMiddlewares()
	r.InitRoutes()
	r.Run()
}
