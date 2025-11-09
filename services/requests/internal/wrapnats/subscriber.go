package wrapnats

import "github.com/nats-io/nats.go"

type Subscriber struct {
	conn          *nats.Conn
	subscriptions []*nats.Subscription
}

func NewSubscriber(conn *nats.Conn) *Subscriber {
	return &Subscriber{
		conn:          conn,
		subscriptions: []*nats.Subscription{},
	}
}
