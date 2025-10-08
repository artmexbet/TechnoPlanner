package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go/jetstream"

	"technoBro/internal/domain"
)

type StreamConfig struct {
	Name     string        `yaml:"name"`
	Subjects []string      `yaml:"subjects"`
	MaxAge   time.Duration `yaml:"max_age"`
}

type Streams struct {
	Something StreamConfig `yaml:"something"`
}

type Config struct {
	URL     string  `yaml:"url" env:"URL"`
	Streams Streams `yaml:"streams" env:"STREAMS"`
}

type Envelope[T any] struct {
	jetstream.Msg
}

func (e *Envelope[T]) Unmarshal() T {
	var t T
	_ = json.Unmarshal(e.Data(), &t)
	return t
}

type NATSBroker struct {
	js  jetstream.JetStream
	cfg Config

	Somethings <-chan Envelope[domain.Something]
}

func NewNATSBroker(cfg Config) *NATSBroker {
	b := &NATSBroker{
		cfg: cfg,
	}
	return b
}

func (b *NATSBroker) WithSomething(ctx context.Context) *NATSBroker {
	b.Somethings = GenEnvelope[domain.Something](ctx, b.js, b.cfg.Streams.Something.Name)
	return b
}

func GenEnvelope[T any](ctx context.Context, js jetstream.JetStream, stream string) <-chan Envelope[T] {
	c, err := js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{Durable: stream})
	if err != nil {
		panic(err)
	}
	res := make(chan Envelope[T])
	_, _ = c.Consume(func(msg jetstream.Msg) {
		res <- Envelope[T]{msg}
	})
	return res
}

// Process apply proccessorFn to value T from ch chanel
func Process[T any](ctx context.Context, ch <-chan Envelope[T], processorFn func(T) error) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-ch:
			err := processorFn(e.Unmarshal())
			if err != nil {
				return // log | aggregate error
			}
		}
	}
}
