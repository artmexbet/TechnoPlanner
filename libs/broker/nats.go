package broker

import (
	"context"

	"github.com/nats-io/nats.go"
)

type Msg struct {
	*nats.Msg

	ctx context.Context
}

func (m *Msg) Context() context.Context {
	return m.ctx
}

func (m *Msg) SetContext(ctx context.Context) {
	m.ctx = ctx
}

type MsgHandler func(msg *Msg) error

type Middleware func(MsgHandler) MsgHandler

type NATS struct {
	*nats.Conn
	middlewares []Middleware
}

func (n *NATS) Subscribe(subject string, handler MsgHandler) (*nats.Subscription, error) {
	return n.Conn.Subscribe(subject, func(msg *nats.Msg) {
		// Apply middlewares
		for i := len(n.middlewares) - 1; i >= 0; i-- {
			handler = n.middlewares[i](handler)
		}

		err := handler(&Msg{Msg: msg, ctx: context.Background()}) // можно использовать любую сигнатуру функции
		if err != nil {
			return
		}
	})
}

func Connect(url string, options ...nats.Option) (*NATS, error) {
	nc, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	return &NATS{Conn: nc}, nil
}

func (n *NATS) Use(mw Middleware) {
	n.middlewares = append(n.middlewares, mw)
}
