package main

import (
	"context"

	"technoBro/internal/app/api"
	"technoBro/internal/app/api/service"
	"technoBro/internal/broker"
	"technoBro/internal/config"
)

type Config struct {
	Router api.Config    `yaml:"router" env:"ROUTER"`
	Broker broker.Config `yaml:"broker" env:"BROKER"`
}

func main() {
	cfg := config.MustParseConfig[Config]("config/cfg.yaml")

	ctx := context.Background()

	nats := broker.NewNATSBroker(cfg.Broker).
		WithSomething(ctx)

	techSvc := service.NewTechManager(nats)
	taskSvc := service.NewTaskManager(nats)

	r := api.NewRouter(cfg.Router, techSvc, taskSvc)
	r.InitMiddlewares()
	r.InitRoutes()
	r.Run()
}
