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

// PorterClient клиент для работы с портерами через NATS (Requests сервис)
type PorterClient struct {
	conn *broker.NATS
}

// NewPorterClient создает новый клиент для работы с портерами
func NewPorterClient(conn *broker.NATS) *PorterClient {
	return &PorterClient{conn: conn}
}

// List получает список всех портеров
func (c *PorterClient) List(ctx context.Context) ([]domain.Porter, error) {
	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayPorterList, nil)
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

	var result []dto.Porter
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal porters: %w", err)
	}

	porters := make([]domain.Porter, 0, len(result))
	for _, p := range result {
		porters = append(porters, domain.Porter{
			ID:       p.ID,
			Username: p.Username,
		})
	}

	return porters, nil
}

// Get получает портера по ID
func (c *PorterClient) Get(ctx context.Context, id uuid.UUID) (domain.Porter, error) {
	req := dto.UUIDRequest{ID: id}
	data, err := req.MarshalJSON()
	if err != nil {
		return domain.Porter{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayPorterGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Porter{}, domain.ErrNotFound
		}
		return domain.Porter{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Porter{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Porter{}, domain.ErrNotFound
		}
		return domain.Porter{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.Porter
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Porter{}, fmt.Errorf("unmarshal porter: %w", err)
	}

	return domain.Porter{
		ID:       result.ID,
		Username: result.Username,
	}, nil
}

// Delete удаляет портера по ID
func (c *PorterClient) Delete(ctx context.Context, id uuid.UUID) error {
	req := dto.UUIDRequest{ID: id}
	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayPorterDelete, data)
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
		return fmt.Errorf("service error: %s", resp.Message)
	}

	return nil
}

// Save сохраняет (upsert) портера в Requests сервисе
func (c *PorterClient) Save(ctx context.Context, id uuid.UUID, username string) error {
	req := dto.PorterSaveRequest{ID: id, Username: username}
	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayPorterSave, data)
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
		return fmt.Errorf("service error: %s", resp.Message)
	}

	return nil
}
