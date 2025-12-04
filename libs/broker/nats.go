package broker

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type Msg struct {
	*nats.Msg

	ctx context.Context
}

func (m *Msg) Context() context.Context {
	return m.ctx
}

type MsgHandler func(msg *Msg)

type NATS struct {
	*nats.Conn
}

func (n *NATS) Subscribe(subject string, handler MsgHandler) (*nats.Subscription, error) {
	return n.Conn.Subscribe(subject, func(msg *nats.Msg) {
		slog.Debug("Received message", "subject", subject, "data", string(msg.Data))
		handler(&Msg{Msg: msg, ctx: context.Background()}) // можно использовать любую сигнатуру функции
	})
}

func Connect(url string, options ...nats.Option) (*NATS, error) {
	nc, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	return &NATS{Conn: nc}, nil
}
