package wrapnats

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"requests/internal/domain"
)

const (
	requestCreatedSubject  = "events.request.created"
	requestCanceledSubject = "events.request.canceled"
	userAdded              = "events.user.added"
)

type NatsPublisher struct {
	conn *nats.Conn
}

func NewNatsPublisher(conn *nats.Conn) *NatsPublisher {
	return &NatsPublisher{
		conn: conn,
	}
}

func (n *NatsPublisher) PublishRequestCreated(req domain.Request) error {
	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}
	err = n.conn.Publish(requestCreatedSubject, data)
	if err != nil {
		return fmt.Errorf("publishing request created event: %w", err)
	}
	return nil
}

func (n *NatsPublisher) PublishRequestCanceled(req domain.Request) error {
	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}
	err = n.conn.Publish(requestCanceledSubject, data)
	if err != nil {
		return fmt.Errorf("publishing request canceled event: %w", err)
	}
	return nil
}

func (n *NatsPublisher) PublishUserAdded(user domain.User) error {
	data, err := user.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshaling user: %w", err)
	}
	err = n.conn.Publish(userAdded, data)
	if err != nil {
		return fmt.Errorf("publishing user added event: %w", err)
	}
	return nil
}
