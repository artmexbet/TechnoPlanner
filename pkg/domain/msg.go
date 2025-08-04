package domain

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Msg struct {
	*nats.Msg
}

// Unmarshal returns the message payload as a Go value.
func (m *Msg) Unmarshal(v any) error {
	return json.Unmarshal(m.Data, v)
}
