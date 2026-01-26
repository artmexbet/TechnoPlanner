package publisher

import (
	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/broker/middleware"
	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc *broker.NATS
}

func New(cfg config.NATSConfig) (*Publisher, error) {
	nc, err := broker.Connect(cfg.URL(), nats.Name("Gateway Service"))
	if err != nil {
		return nil, err
	}

	// Apply middlewares
	nc.Use(middleware.NewLoggingMiddleware(true))
	nc.Use(middleware.NewRecoveryMiddleware())
	nc.Use(middleware.NewRequestIDMiddleware())
	return &Publisher{nc: nc}, nil
}
