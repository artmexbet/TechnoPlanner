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
	req := dto.TechEquipmentCategoryCreateRequest{
		Name:        cat.Name,
		Description: cat.Description,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentCategoryCreate, data)
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

	var result dto.TechEquipmentCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("unmarshal category: %w", err)
	}

	return mapEquipmentCategoryFromDTO(result), nil
}

// Update обновляет категорию
func (c *CategoryClient) Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	req := dto.TechEquipmentCategoryUpdateRequest{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentCategoryUpdate, data)
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

	var result dto.TechEquipmentCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("unmarshal category: %w", err)
	}

	return mapEquipmentCategoryFromDTO(result), nil
}

// List получает список категорий
func (c *CategoryClient) List(ctx context.Context) ([]domain.EquipmentCategory, error) {
	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentCategoryList, []byte("{}"))
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

	var result []dto.TechEquipmentCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal category list: %w", err)
	}

	categories := make([]domain.EquipmentCategory, 0, len(result))
	for _, cat := range result {
		categories = append(categories, mapEquipmentCategoryFromDTO(cat))
	}

	return categories, nil
}

// Delete удаляет категорию
func (c *CategoryClient) Delete(ctx context.Context, id int) error {
	req := dto.TechEquipmentCategoryDeleteRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayEquipmentCategoryDelete, data)
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

func mapEquipmentCategoryFromDTO(cat dto.TechEquipmentCategory) domain.EquipmentCategory {
	desc := ""
	if cat.Description != nil {
		desc = *cat.Description
	}
	return domain.EquipmentCategory{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: desc,
		CreatedAt:   cat.Audit.CreatedAt,
		UpdatedAt:   cat.Audit.UpdatedAt,
	}
}
