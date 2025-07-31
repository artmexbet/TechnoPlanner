package broker

import (
	"github.com/nats-io/nats.go"
)

type NATSBroker struct {
	conn        nats.Conn
	connections map[string][]nats.Subscription
}

type Tmp interface {
	Add()
}

func (n *NATSBroker) RegisterConsumer(topic string) error {
	//sub, err := n.conn.Subscribe(topic, )
	//if err != nil {
	//	return errors.Join(err, errors.New("failed to subscribe to topic"))
	//}
	return nil
}
