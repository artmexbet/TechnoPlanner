package broker

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

// Handler represents a handler func that receives T (domain struct)
// and do something on the domain side
type Handler[T any] func(T) error

// Consumer is used to be a couple between NATS and the app.
type Consumer[T any] struct {
	handler Handler[T]
}

func (h *Consumer[T]) unmarshal(msg jetstream.Msg) (T, error) {
	var v T
	err := json.Unmarshal(msg.Data(), &v)
	return v, err
}

func (h *Consumer[T]) Handle(msg jetstream.Msg) error {
	res, err := h.unmarshal(msg)
	if err != nil {
		return fmt.Errorf("error unmarshalling message: %v", err)
	}
	return h.handler(res)
}

func NewConsumer[T any](handler Handler[T]) *Consumer[T] {
	return &Consumer[T]{
		handler: handler,
	}
}
