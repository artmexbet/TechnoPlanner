package wrapnats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

// EquipmentClient — NATS клиент для резервации оборудования в equipment сервисе.
type EquipmentClient struct {
	conn *broker.NATS
}

// NewEquipmentClient создаёт новый клиент для Equipment Service.
func NewEquipmentClient(conn *broker.NATS) *EquipmentClient {
	return &EquipmentClient{conn: conn}
}

// ReserveEquipment резервирует оборудование для заявки.
func (c *EquipmentClient) ReserveEquipment(ctx context.Context, items []domain.EquipmentReserveItem) error {
	req := dto.TechEquipmentReserveRequest{Items: toDTOItems(items)}
	data, err := easyjson.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal reserve request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentReserve, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return fmt.Errorf("equipment service not available: %w", err)
		}
		return fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("reserve failed: %s", resp.Message)
	}

	return nil
}

// ReleaseEquipment освобождает зарезервированное оборудование.
func (c *EquipmentClient) ReleaseEquipment(ctx context.Context, items []domain.EquipmentReserveItem) error {
	req := dto.TechEquipmentReleaseRequest{Items: toDTOItems(items)}
	data, err := easyjson.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal release request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentRelease, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return fmt.Errorf("equipment service not available: %w", err)
		}
		return fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("release failed: %s", resp.Message)
	}

	return nil
}

func toDTOItems(items []domain.EquipmentReserveItem) []dto.EquipmentReserveItem {
	result := make([]dto.EquipmentReserveItem, len(items))
	for i, item := range items {
		result[i] = dto.EquipmentReserveItem{
			EquipmentID: item.EquipmentID,
			Quantity:    item.Quantity,
		}
	}
	return result
}
