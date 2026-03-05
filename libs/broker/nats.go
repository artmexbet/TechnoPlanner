package broker

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	propagator  propagation.TextMapPropagator
}

// natsHeaderCarrier реализует propagation.TextMapCarrier для NATS headers
type natsHeaderCarrier nats.Header

func (c natsHeaderCarrier) Get(key string) string {
	return nats.Header(c).Get(key)
}

func (c natsHeaderCarrier) Set(key, value string) {
	nats.Header(c).Set(key, value)
}

func (c natsHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

func (n *NATS) Subscribe(subject string, handler MsgHandler) (*nats.Subscription, error) {
	return n.Conn.Subscribe(subject, func(msg *nats.Msg) {
		// Apply middlewares
		for i := len(n.middlewares) - 1; i >= 0; i-- {
			handler = n.middlewares[i](handler)
		}

		// Extract context from NATS headers
		ctx := context.Background()
		if msg.Header != nil && n.propagator != nil {
			ctx = n.propagator.Extract(ctx, natsHeaderCarrier(msg.Header))
		}

		err := handler(&Msg{Msg: msg, ctx: ctx})
		if err != nil {
			return
		}
	})
}

// PublishWithContext публикует сообщение с контекстом (trace context инъектируется в headers)
func (n *NATS) PublishWithContext(ctx context.Context, subject string, data []byte) error {
	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  make(nats.Header),
	}

	if n.propagator != nil {
		n.propagator.Inject(ctx, natsHeaderCarrier(msg.Header))
	}

	return n.PublishMsg(msg)
}

// RequestWithContext отправляет запрос с контекстом и ожидает ответ
func (n *NATS) RequestWithContext(ctx context.Context, subject string, data []byte) (*nats.Msg, error) {
	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  make(nats.Header),
	}

	if n.propagator != nil {
		n.propagator.Inject(ctx, natsHeaderCarrier(msg.Header))
	}

	return n.RequestMsgWithContext(ctx, msg)
}

func Connect(url string, options ...nats.Option) (*NATS, error) {
	nc, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	return &NATS{
		Conn:       nc,
		propagator: otel.GetTextMapPropagator(),
	}, nil
}

// WrapConn оборачивает существующий *nats.Conn в *NATS (удобно в тестах).
func WrapConn(nc *nats.Conn) *NATS {
	return &NATS{
		Conn:       nc,
		propagator: otel.GetTextMapPropagator(),
	}
}

func (n *NATS) Use(mw Middleware) {
	n.middlewares = append(n.middlewares, mw)
}

// SetPropagator позволяет установить кастомный propagator
func (n *NATS) SetPropagator(p propagation.TextMapPropagator) {
	n.propagator = p
}
