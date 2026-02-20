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

// ResponsibleClient клиент для работы с ответственными через NATS
type ResponsibleClient struct {
	conn *broker.NATS
}

// NewResponsibleClient создает новый клиент для работы с ответственными
func NewResponsibleClient(conn *broker.NATS) *ResponsibleClient {
	return &ResponsibleClient{conn: conn}
}

// List получает список всех ответственных
func (c *ResponsibleClient) List(ctx context.Context) ([]domain.Responsible, error) {
	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayResponsibleList, nil)
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

	var result []dto.Responsible
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal responsibles: %w", err)
	}

	responsibles := make([]domain.Responsible, 0, len(result))
	for _, r := range result {
		responsibles = append(responsibles, domain.Responsible{
			ID:       r.ID,
			Username: r.Username,
		})
	}

	return responsibles, nil
}

// Create создает нового ответственного
func (c *ResponsibleClient) Create(ctx context.Context, id uuid.UUID, username string) (domain.Responsible, error) {
	req := dto.ResponsibleCreateRequest{
		ID:       id,
		Username: username,
	}

	data, err := req.MarshalJSON()
	if err != nil {
		return domain.Responsible{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayResponsibleCreate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Responsible{}, domain.ErrNotFound
		}
		return domain.Responsible{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Responsible{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		return domain.Responsible{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.Responsible
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Responsible{}, fmt.Errorf("unmarshal responsible: %w", err)
	}

	return domain.Responsible{
		ID:       result.ID,
		Username: result.Username,
	}, nil
}
