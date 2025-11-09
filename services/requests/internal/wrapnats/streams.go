package wrapnats

import (
	"log/slog"

	"github.com/nats-io/nats.go"
)

func (w *NatsWrapper) SetupStreams() *NatsWrapper {
	js, err := w.conn.JetStream()
	if err != nil {
		slog.Error("error getting jetstream context", "error", err)
	}
	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     "SOMETHING",
		Subjects: []string{"something.*"},
	})
	return w
}
