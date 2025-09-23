package var1

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go/jetstream"

	"technoBro/internal/domain"
)

type Envelope[T any] struct {
	jetstream.Msg
}

func (e *Envelope[T]) Unmarshal() T {
	var t T
	_ = json.Unmarshal(e.Data(), &t)
	return t
}

type NATSBroker struct {
	js jetstream.JetStream

	Somethings <-chan Envelope[domain.Something]
	/*
		Здесь мы по сути вообще разделяем логику вычитывания и обработки
	*/
}

type Option func(*NATSBroker)

func NewNATSBroker(opts ...Option) *NATSBroker {
	b := &NATSBroker{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func WithSomething() Option {
	return func(b *NATSBroker) {
		// b.Somethings = GenEnvelope[domain.Something](...)
	}
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

// Usage

func Process[T any](ctx context.Context, ch <-chan Envelope[T], processorFn func(T) error) {
	go func() {
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
	}()
}

/*
func main() {
	br := NewNATSBroker()
	ctx := context.Background()
	br.Something = GenEnvelope[domain.Something](ctx, ..., "something")

	Process(ctx, br.Something, func(d domain.Something) error {})  // А ещё можно вот эту тему вынести в пул воркеров

}
*/
