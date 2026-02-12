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

// TechClient клиент для синхронных вызовов Tech Service через NATS Request-Reply
type TechClient struct {
	conn *broker.NATS
}

// NewTechClient создает новый клиент для Tech Service
func NewTechClient(conn *broker.NATS) *TechClient {
	return &TechClient{conn: conn}
}

// CreateTechnic создает новую технику
func (c *TechClient) CreateTechnic(ctx context.Context, req dto.TechnicCreateRequest) (dto.Technic, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return dto.Technic{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicCreate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return dto.Technic{}, domain.ErrNotFound
		}
		return dto.Technic{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return dto.Technic{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return dto.Technic{}, domain.ErrNotFound
		}
		return dto.Technic{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.Technic
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return dto.Technic{}, fmt.Errorf("unmarshal technic: %w", err)
	}

	return result, nil
}

// UpdateTechnic обновляет технику
func (c *TechClient) UpdateTechnic(ctx context.Context, req dto.TechnicUpdateRequest) (dto.Technic, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return dto.Technic{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicUpdate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return dto.Technic{}, domain.ErrNotFound
		}
		return dto.Technic{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return dto.Technic{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return dto.Technic{}, domain.ErrNotFound
		}
		return dto.Technic{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.Technic
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return dto.Technic{}, fmt.Errorf("unmarshal technic: %w", err)
	}

	return result, nil
}

// GetTechnic получает технику по ID
func (c *TechClient) GetTechnic(ctx context.Context, id uuid.UUID) (dto.Technic, error) {
	req := dto.TechnicGetByIDRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return dto.Technic{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return dto.Technic{}, domain.ErrNotFound
		}
		return dto.Technic{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return dto.Technic{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return dto.Technic{}, domain.ErrNotFound
		}
		return dto.Technic{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.Technic
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return dto.Technic{}, fmt.Errorf("unmarshal technic: %w", err)
	}

	return result, nil
}

// DeleteTechnic удаляет технику
func (c *TechClient) DeleteTechnic(ctx context.Context, id uuid.UUID) error {
	req := dto.TechnicDeleteRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicDelete, data)
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

// ListTechnics получает список техники
func (c *TechClient) ListTechnics(ctx context.Context) ([]dto.Technic, error) {
	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicList, []byte("{}"))
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

	var result []dto.Technic
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal technic list: %w", err)
	}

	return result, nil
}

// GetTechnicsByCategory получает технику по категории
func (c *TechClient) GetTechnicsByCategory(ctx context.Context, categoryID uuid.UUID) ([]dto.Technic, error) {
	req := dto.TechnicGetByCategoryRequest{CategoryID: categoryID}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicGetByCategory, data)
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
		if resp.Message == "not found" {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("service error: %s", resp.Message)
	}

	var result []dto.Technic
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal technic list: %w", err)
	}

	return result, nil
}

// CreateCategory создает категорию техники
func (c *TechClient) CreateCategory(ctx context.Context, req dto.TechnicCategoryCreateRequest) (dto.TechnicCategory, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicCategoryCreate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return dto.TechnicCategory{}, domain.ErrNotFound
		}
		return dto.TechnicCategory{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return dto.TechnicCategory{}, domain.ErrNotFound
		}
		return dto.TechnicCategory{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.TechnicCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("unmarshal category: %w", err)
	}

	return result, nil
}

// UpdateCategory обновляет категорию техники
func (c *TechClient) UpdateCategory(ctx context.Context, req dto.TechnicCategoryUpdateRequest) (dto.TechnicCategory, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicCategoryUpdate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return dto.TechnicCategory{}, domain.ErrNotFound
		}
		return dto.TechnicCategory{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return dto.TechnicCategory{}, domain.ErrNotFound
		}
		return dto.TechnicCategory{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.TechnicCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("unmarshal category: %w", err)
	}

	return result, nil
}

// GetCategory получает категорию по ID
func (c *TechClient) GetCategory(ctx context.Context, id uuid.UUID) (dto.TechnicCategory, error) {
	req := dto.TechnicCategoryGetByIDRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicCategoryGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return dto.TechnicCategory{}, domain.ErrNotFound
		}
		return dto.TechnicCategory{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return dto.TechnicCategory{}, domain.ErrNotFound
		}
		return dto.TechnicCategory{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.TechnicCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return dto.TechnicCategory{}, fmt.Errorf("unmarshal category: %w", err)
	}

	return result, nil
}

// DeleteCategory удаляет категорию техники
func (c *TechClient) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	req := dto.TechnicCategoryDeleteRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicCategoryDelete, data)
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

// ListCategories получает список категорий техники
func (c *TechClient) ListCategories(ctx context.Context) ([]dto.TechnicCategory, error) {
	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayTechnicCategoryList, []byte("{}"))
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

	var result []dto.TechnicCategory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal category list: %w", err)
	}

	return result, nil
}
