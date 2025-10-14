package main

import (
	"broker"
	"config"
	"context"
	"gateway/internal/app"
	service2 "gateway/internal/app/service"

	"gateway/internal/postgres"
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

	nats := broker.NewNATSBroker(cfg.Broker)

	_postgres, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer _postgres.Close()

	store := storage.NewStorage(_postgres) // инициализация хранилища
	_ = store                              // заглушка, чтобы не было ошибки о неиспользуемой переменной

	techSvc := service2.NewTechManager(nats)
	taskSvc := service2.NewTaskManager(nats)

	r := app.NewRouter(cfg.Router, techSvc, taskSvc)
	r.InitMiddlewares()
	r.InitRoutes()
	r.Run()
}
