package wrapnats

import (
	"fmt"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type NatsPublisher struct {
	conn *broker.NATS
}

func NewNatsPublisher(conn *broker.NATS) *NatsPublisher {
	return &NatsPublisher{
		conn: conn,
	}
}

func (n *NatsPublisher) PublishRequestCreated(req domain.Request) error {
	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}
	err = n.conn.Publish(subjects.ServiceRequestCreated, data)
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
	err = n.conn.Publish(subjects.ServiceRequestCanceled, data)
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
	err = n.conn.Publish(subjects.ServiceUserAdded, data)
	if err != nil {
		return fmt.Errorf("publishing user added event: %w", err)
	}
	return nil
}
