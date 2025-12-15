package natspublisher

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"

	"github.com/artmexbet/TechnoPlanner/services/auth/internal/models"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(cfg config.NATSConfig) (*Publisher, error) {
	nc, err := nats.Connect(
		cfg.URL(),
		nats.Name("auth-service-publisher"),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats: %w", err)
	}
	return &Publisher{conn: nc}, nil
}

func (p *Publisher) Close() {
	p.conn.Close()
}

func (p *Publisher) PublishUserCreated(user models.User) error {
	data, err := user.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error marshaling user created event: %w", err)
	}
	if err = p.conn.Publish(subjects.UserCreated, data); err != nil {
		return fmt.Errorf("error publishing user created event: %w", err)
	}
	return nil
}
