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

// CategoryClient клиент для синхронных вызовов Category Service через NATS Request-Reply
type CategoryClient struct {
	conn *broker.NATS
}

// NewCategoryClient создает новый клиент для Category Service
func NewCategoryClient(conn *broker.NATS) *CategoryClient {
	return &CategoryClient{conn: conn}
}

// Create создает новую категорию
func (c *CategoryClient) Create(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	req := dto.CategoryCreateRequest{
		Name:        cat.Name,
		Description: derefString(cat.Description),
	}

	data, err := req.MarshalJSON()
	if err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayCategoryCreate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.EquipmentCategory{}, domain.ErrNotFound
		}
		return domain.EquipmentCategory{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.EquipmentCategory{}, domain.ErrNotFound
		}
		return domain.EquipmentCategory{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.EquipmentCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("unmarshal category: %w", err)
	}

	return mapCategoryFromDTO(result), nil
}

// Update обновляет категорию
func (c *CategoryClient) Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	req := dto.CategoryUpdateRequest{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: derefString(cat.Description),
	}

	data, err := req.MarshalJSON()
	if err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayCategoryUpdate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.EquipmentCategory{}, domain.ErrNotFound
		}
		return domain.EquipmentCategory{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.EquipmentCategory{}, domain.ErrNotFound
		}
		return domain.EquipmentCategory{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.EquipmentCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("unmarshal category: %w", err)
	}

	return mapCategoryFromDTO(result), nil
}

// List получает список категорий
func (c *CategoryClient) List(ctx context.Context) ([]domain.EquipmentCategory, error) {
	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayCategoryList, []byte("{}"))
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

	var result []dto.EquipmentCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal category list: %w", err)
	}

	categories := make([]domain.EquipmentCategory, 0, len(result))
	for _, cat := range result {
		categories = append(categories, mapCategoryFromDTO(cat))
	}

	return categories, nil
}

// SoftDelete мягко удаляет категорию
func (c *CategoryClient) SoftDelete(ctx context.Context, id int32, userID *uuid.UUID) error {
	req := dto.SoftDeleteRequest{
		ID:     id,
		UserID: userID,
	}

	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayCategoryDelete, data)
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

func mapCategoryFromDTO(cat dto.EquipmentCategory) domain.EquipmentCategory {
	return domain.EquipmentCategory{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		Audit: domain.AuditFields{
			CreatedAt: cat.Audit.CreatedAt,
			UpdatedAt: cat.Audit.UpdatedAt,
			DeletedAt: cat.Audit.DeletedAt,
			CreatedBy: cat.Audit.CreatedBy,
			UpdatedBy: cat.Audit.UpdatedBy,
		},
	}
}
