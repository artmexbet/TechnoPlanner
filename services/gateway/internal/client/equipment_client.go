package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

// EquipmentClient клиент для синхронных вызовов Equipment Service через NATS Request-Reply
type EquipmentClient struct {
	conn *broker.NATS
}

// NewEquipmentClient создает новый клиент для Equipment Service
func NewEquipmentClient(conn *broker.NATS) *EquipmentClient {
	return &EquipmentClient{conn: conn}
}

// Create создает новое оборудование
func (c *EquipmentClient) Create(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	req := dto.TechEquipmentCreateRequest{
		CategoryID:                eq.CategoryID,
		Name:                      eq.Name,
		Description:               eq.Description,
		AdditionalCharacteristics: eq.AdditionalCharacteristics,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.Equipment{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentCreate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Equipment{}, domain.ErrNotFound
		}
		return domain.Equipment{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Equipment{}, domain.ErrNotFound
		}
		return domain.Equipment{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.TechEquipment
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal equipment: %w", err)
	}

	return mapEquipmentFromDTO(result), nil
}

// Update обновляет оборудование
func (c *EquipmentClient) Update(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	req := dto.TechEquipmentUpdateRequest{
		ID:                        eq.ID,
		CategoryID:                eq.CategoryID,
		Name:                      eq.Name,
		Description:               eq.Description,
		AdditionalCharacteristics: eq.AdditionalCharacteristics,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.Equipment{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentUpdate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Equipment{}, domain.ErrNotFound
		}
		return domain.Equipment{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Equipment{}, domain.ErrNotFound
		}
		return domain.Equipment{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.TechEquipment
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal equipment: %w", err)
	}

	return mapEquipmentFromDTO(result), nil
}

// Get получает оборудование по ID
func (c *EquipmentClient) Get(ctx context.Context, id int) (domain.Equipment, error) {
	req := dto.TechEquipmentGetByIDRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.Equipment{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Equipment{}, domain.ErrNotFound
		}
		return domain.Equipment{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Equipment{}, domain.ErrNotFound
		}
		return domain.Equipment{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.TechEquipment
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal equipment: %w", err)
	}

	return mapEquipmentFromDTO(result), nil
}

// List получает список оборудования
func (c *EquipmentClient) List(ctx context.Context) ([]domain.Equipment, error) {
	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentList, []byte("{}"))
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("service error: %s", resp.Message)
	}

	var result []dto.TechEquipment
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal equipment list: %w", err)
	}

	equipment := make([]domain.Equipment, 0, len(result))
	for _, eq := range result {
		equipment = append(equipment, mapEquipmentFromDTO(eq))
	}

	return equipment, nil
}

// Delete удаляет оборудование
func (c *EquipmentClient) Delete(ctx context.Context, id int) error {
	req := dto.TechEquipmentDeleteRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentDelete, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.ErrNotFound
		}
		return fmt.Errorf("service error: %s", resp.Message)
	}

	return nil
}

func mapEquipmentFromDTO(eq dto.TechEquipment) domain.Equipment {
	return domain.Equipment{
		ID:                        eq.ID,
		CategoryID:                eq.CategoryID,
		Name:                      eq.Name,
		Description:               eq.Description,
		AdditionalCharacteristics: eq.AdditionalCharacteristics,
		CreatedAt:                 eq.CreatedAt,
		UpdatedAt:                 eq.UpdatedAt,
	}
}
