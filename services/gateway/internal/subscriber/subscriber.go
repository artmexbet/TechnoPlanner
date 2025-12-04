package subscriber

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	"gateway/internal/domain"

	"broker"

	"config"
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
	sub, err := s.nc.Subscribe(config.SubjectUserCreated, func(msg *broker.Msg) {
		var user domain.User
		if err := user.UnmarshalJSON(msg.Data); err != nil {
			slog.Error("error unmarshalling user created event", "error", err)
			return
		}

		if err := s.repo.Create(msg.Context(), user); err != nil {
			slog.Error("error creating user from user created event", "error", err)
			return
		}
	})

	if err != nil {
		return fmt.Errorf("error subscribing to %s: %w", config.SubjectUserCreated, err)
	}
	s.subs = append(s.subs, sub)
	return nil
}
