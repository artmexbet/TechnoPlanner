package wrapnats

import (
	"encoding/json"
	"log/slog"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"
	"github.com/artmexbet/TechnoPlanner/libs/dto"
)

// NatsPublisher публикует события об изменениях оборудования.
type NatsPublisher struct {
	conn *broker.NATS
}

func NewNatsPublisher(conn *broker.NATS) *NatsPublisher {
	return &NatsPublisher{conn: conn}
}

func (p *NatsPublisher) PublishEquipmentCreated(eq dto.EquipmentSyncEvent) {
	data, err := json.Marshal(eq)
	if err != nil {
		slog.Error("publish EquipmentCreated: marshal error", "error", err)
		return
	}
	if err := p.conn.Publish(subjects.EquipmentCreated, data); err != nil {
		slog.Error("publish EquipmentCreated: nats error", "error", err)
	}
}

func (p *NatsPublisher) PublishEquipmentUpdated(eq dto.EquipmentSyncEvent) {
	data, err := json.Marshal(eq)
	if err != nil {
		slog.Error("publish EquipmentUpdated: marshal error", "error", err)
		return
	}
	if err := p.conn.Publish(subjects.EquipmentUpdated, data); err != nil {
		slog.Error("publish EquipmentUpdated: nats error", "error", err)
	}
}

func (p *NatsPublisher) PublishEquipmentDeleted(id int) {
	ev := dto.EquipmentDeletedEvent{ID: id}
	data, err := json.Marshal(ev)
	if err != nil {
		slog.Error("publish EquipmentDeleted: marshal error", "error", err)
		return
	}
	if err := p.conn.Publish(subjects.EquipmentDeleted, data); err != nil {
		slog.Error("publish EquipmentDeleted: nats error", "error", err)
	}
}
