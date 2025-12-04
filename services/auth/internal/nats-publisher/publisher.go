package natspublisher

import (
	"auth/internal/models"
	"config"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(host string, port int) (*Publisher, error) {
	nc, err := nats.Connect(
		fmt.Sprintf("nats://%s:%d", host, port),
		nats.Name("auth-service-publisher"),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats: %s", err)
	}
	return &Publisher{conn: nc}, nil
}

func (p *Publisher) Close() {
	p.conn.Close()
}

func (p *Publisher) PublishUserCreated(user models.User) error {
	data, err := user.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error marshalling user created event: %w", err)
	}
	return p.conn.Publish(config.SubjectUserCreated, data)
}
