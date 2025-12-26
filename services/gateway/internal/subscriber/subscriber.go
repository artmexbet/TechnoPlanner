package subscriber

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/broker/middleware"
	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type iRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type Subscriber struct {
	nc   *broker.NATS
	subs []*nats.Subscription

	repo iRepository
}

func New(cfg config.NATSConfig, repo iRepository) (*Subscriber, error) {
	nc, err := broker.Connect(cfg.URL(), nats.Name("Gateway Service"))
	if err != nil {
		return nil, err
	}

	// Apply middlewares
	nc.Use(middleware.NewLoggingMiddleware(true))
	nc.Use(middleware.NewRecoveryMiddleware())
	nc.Use(middleware.NewRequestIDMiddleware())

	return &Subscriber{nc: nc, repo: repo}, nil
}

func (s *Subscriber) Close() {
	for _, sub := range s.subs {
		if err := sub.Unsubscribe(); err != nil {
			slog.Error("error unsubscribing from subject", "subject", sub.Subject, "error", err)
		}
	}
	s.nc.Close()
}

func (s *Subscriber) Init() error {
	if err := s.subscribeUserCreated(); err != nil {
		return fmt.Errorf("error initializing subscriber: %w", err)
	}
	return nil
}

func (s *Subscriber) subscribeUserCreated() error {
	sub, err := s.nc.Subscribe(subjects.UserCreated, func(msg *broker.Msg) error {
		var user domain.User
		if err := user.UnmarshalJSON(msg.Data); err != nil {
			return fmt.Errorf("error unmarshalling user data: %w", err)
		}

		if err := s.repo.Create(msg.Context(), user); err != nil {
			return fmt.Errorf("error creating user from user created event: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error subscribing to %s: %w", subjects.UserCreated, err)
	}
	s.subs = append(s.subs, sub)
	return nil
}
