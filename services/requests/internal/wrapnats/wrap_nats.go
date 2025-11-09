package wrapnats

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"URL"`
}

type NatsWrapper struct {
	conn *nats.Conn
}

func New(cfg Config) (*NatsWrapper, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats server: %v", err)
	}
	return &NatsWrapper{conn: conn}, nil
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]nats.MsgHandler{
		"base": w.baseHandler,
	}

	for subject, handler := range handlers {
		_, err := w.conn.Subscribe(subject, handler)
		if err != nil {
			slog.Error("cannot subscribe to subject", "subject", subject, "error", err)
		}
	}
	return w
}

func (w *NatsWrapper) baseHandler(msg *nats.Msg) {

}
