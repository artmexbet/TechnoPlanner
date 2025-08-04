package main

import (
	"technoBro/internal/app/api"
	"technoBro/internal/app/api/service"

	"technoBro/pkg/broker"
	"technoBro/pkg/config"
)

type Config struct {
	Router api.Config    `yaml:"router" env:"ROUTER"`
	Broker broker.Config `yaml:"broker" env:"BROKER"`
}

func main() {
	cfg := config.MustParseConfig[Config]("config.yaml")

	nats, err := broker.NewNATSBroker(cfg.Broker, nil)
	if err != nil {
		panic(err)
	}

	techSvc := service.NewTechManager(nats)
	taskSvc := service.NewTaskManager(nats)

	r := api.NewRouter(cfg.Router, techSvc, taskSvc)
	r.InitMiddlewares()
	r.Run()
}
