package broker

import (
	"encoding/json"
	"errors"
	"fmt"

	"technoBro/pkg/domain"

	"github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"URL"`
}

type Handler func(msg domain.Msg)
type Consumers map[string]Handler

type NATSBroker struct {
	conn *nats.Conn
}

func NewNATSBroker(cfg Config, consumers Consumers) (*NATSBroker, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	n := &NATSBroker{conn: conn}
	for topic, handler := range consumers {
		err = n.RegisterConsumer(topic, handler)
		if err != nil {
			return nil, fmt.Errorf("failed to register consumer for topic %s: %w", topic, err)
		}
	}
	return n, nil
}

func (n *NATSBroker) RegisterConsumer(topic string, handler Handler) error {
	_, err := n.conn.Subscribe(topic, func(msg *nats.Msg) {
		handler(domain.Msg{Msg: msg})
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}
	return nil
}

func (n *NATSBroker) Publish(topic string, msg any) error {
	data, err := json.Marshal(msg) // Верим, что везде будем использовать json
	if err != nil {
		return errors.Join(err, errors.New("failed to marshal message"))
	}
	return n.conn.Publish(topic, data)
}
