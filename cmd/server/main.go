package main

import (
	"context"

	"technoBro/internal/app/api"
	"technoBro/internal/app/api/service"
	"technoBro/internal/domain"
	"technoBro/pkg/broker"
	"technoBro/pkg/config"
)

type Config struct {
	Router api.Config    `yaml:"router" env:"ROUTER"`
	Broker broker.Config `yaml:"broker" env:"BROKER"`
}

func main() {
	cfg := config.MustParseConfig[Config]("config/cfg.yaml")

	nats, err := broker.NewNATSBroker(cfg.Broker)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	nats.WithSmthConsumer(ctx, "SMTHING",
		broker.NewConsumer(
			func(d domain.Something) error {
				return nil // Base example, don't use
			}),
	)

	techSvc := service.NewTechManager(nats)
	taskSvc := service.NewTaskManager(nats)

	r := api.NewRouter(cfg.Router, techSvc, taskSvc)
	r.InitMiddlewares()
	r.InitRoutes()
	r.Run()
}
