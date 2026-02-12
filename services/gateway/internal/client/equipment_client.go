package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
	req := dto.EquipmentCreateRequest{
		Name:        eq.Name,
		Description: derefString(eq.Description),
		Quantity:    eq.Quantity,
	}
	for _, cat := range eq.Categories {
		req.CategoryIDs = append(req.CategoryIDs, cat.ID)
	}

	data, err := req.MarshalJSON()
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

	var result dto.Equipment
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal equipment: %w", err)
	}

	return mapEquipmentFromDTO(result), nil
}

// Update обновляет оборудование
func (c *EquipmentClient) Update(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	req := dto.EquipmentUpdateRequest{
		ID:          eq.ID,
		Name:        eq.Name,
		Description: derefString(eq.Description),
		Quantity:    eq.Quantity,
	}
	for _, cat := range eq.Categories {
		req.CategoryIDs = append(req.CategoryIDs, cat.ID)
	}

	data, err := req.MarshalJSON()
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

	var result dto.Equipment
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Equipment{}, fmt.Errorf("unmarshal equipment: %w", err)
	}

	return mapEquipmentFromDTO(result), nil
}

// Get получает оборудование по ID
func (c *EquipmentClient) Get(ctx context.Context, id int32) (domain.Equipment, error) {
	req := dto.IDRequest{ID: id}

	data, err := req.MarshalJSON()
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

	var result dto.Equipment
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

	var result []dto.Equipment
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal equipment list: %w", err)
	}

	equipment := make([]domain.Equipment, 0, len(result))
	for _, eq := range result {
		equipment = append(equipment, mapEquipmentFromDTO(eq))
	}

	return equipment, nil
}

// SoftDelete мягко удаляет оборудование
func (c *EquipmentClient) SoftDelete(ctx context.Context, id int32, userID *uuid.UUID) error {
	req := dto.SoftDeleteRequest{
		ID:     id,
		UserID: userID,
	}

	data, err := req.MarshalJSON()
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

func mapEquipmentFromDTO(eq dto.Equipment) domain.Equipment {
	categories := make([]domain.EquipmentCategory, 0, len(eq.Categories))
	for _, cat := range eq.Categories {
		categories = append(categories, mapCategoryFromDTO(cat))
	}

	return domain.Equipment{
		ID:          eq.ID,
		Name:        eq.Name,
		Description: eq.Description,
		Quantity:    eq.Quantity,
		Categories:  categories,
		Audit: domain.AuditFields{
			CreatedAt: eq.Audit.CreatedAt,
			UpdatedAt: eq.Audit.UpdatedAt,
			DeletedAt: eq.Audit.DeletedAt,
			CreatedBy: eq.Audit.CreatedBy,
			UpdatedBy: eq.Audit.UpdatedBy,
		},
	}
}
