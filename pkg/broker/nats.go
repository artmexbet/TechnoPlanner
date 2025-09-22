package broker

import (
	"context"
	"fmt"
	"time"

	"technoBro/internal/domain"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type StreamConfig struct {
	Name     string        `yaml:"name"`
	Subjects []string      `yaml:"subjects"`
	MaxAge   time.Duration `yaml:"max_age"`
}

type Config struct {
	URL     string         `yaml:"url" env:"URL"`
	Streams []StreamConfig `yaml:"streams" env:"STREAMS"`
}

type NATSBroker struct {
	js              jetstream.JetStream
	cfg             Config
	consumeContexts map[string]jetstream.ConsumeContext
}

func NewNATSBroker(cfg Config) (*NATSBroker, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to init JetStream: %w", err)
	}
	n := &NATSBroker{
		js:  js,
		cfg: cfg,
	}
	return n, nil
}

func (n *NATSBroker) Start(ctx context.Context) error {
	for _, stream := range n.cfg.Streams {
		_, err := n.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
			Name:     stream.Name,
			Subjects: stream.Subjects,
			MaxAge:   stream.MaxAge,
		})
		if err != nil {
			return fmt.Errorf("failed to create stream %s: %w", stream.Name, err)
		}
	}
	return nil
}

func (n *NATSBroker) Stop(_ context.Context) error {
	for _, consumeContext := range n.consumeContexts {
		consumeContext.Stop()
	}
	return nil
}

func (n *NATSBroker) BuildConsumer(ctx context.Context, stream, name string) (jetstream.Consumer, error) {
	cons, err := n.js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{
		Durable: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer %s.%s: %w", stream, name, err)
	}
	return cons, nil
}

func (n *NATSBroker) WithSmthConsumer(ctx context.Context, stream string, consumer *Consumer[domain.Something]) *NATSBroker {
	cons, err := n.BuildConsumer(ctx, stream, "smth")
	if err != nil {
		panic(err)
	}
	consumeCtx, _ := cons.Consume(func(msg jetstream.Msg) {
		err = consumer.Handle(msg)
		if err != nil {
			fmt.Println(err)
			msg.Nack()
			// TODO: log
		}
		msg.Ack()
	})

	n.consumeContexts["smth"] = consumeCtx
	return n
}
