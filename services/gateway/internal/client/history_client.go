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

// HistoryClient клиент для синхронных вызовов History Service через NATS Request-Reply
type HistoryClient struct {
	conn *broker.NATS
}

// NewHistoryClient создает новый клиент для History Service
func NewHistoryClient(conn *broker.NATS) *HistoryClient {
	return &HistoryClient{conn: conn}
}

// Add добавляет запись в историю статусов
func (c *HistoryClient) Add(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error) {
	req := dto.HistoryAddRequest{
		RequestID: entry.RequestID,
		Status:    dto.RequestStatus(entry.Status),
		Comment:   derefString(entry.Comment),
	}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.RequestStatusHistory{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayHistoryAdd, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.RequestStatusHistory{}, domain.ErrNotFound
		}
		return domain.RequestStatusHistory{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.RequestStatusHistory{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.RequestStatusHistory{}, domain.ErrNotFound
		}
		return domain.RequestStatusHistory{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.RequestStatusHistory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.RequestStatusHistory{}, fmt.Errorf("unmarshal history: %w", err)
	}

	return mapHistoryFromDTO(result), nil
}

// List получает историю статусов по ID заявки
func (c *HistoryClient) List(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error) {
	req := map[string]string{"request_id": requestID.String()}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayHistoryList, data)
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

	var result []dto.RequestStatusHistory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal history list: %w", err)
	}

	history := make([]domain.RequestStatusHistory, 0, len(result))
	for _, h := range result {
		history = append(history, mapHistoryFromDTO(h))
	}

	return history, nil
}

func mapHistoryFromDTO(h dto.RequestStatusHistory) domain.RequestStatusHistory {
	return domain.RequestStatusHistory{
		ID:        h.ID,
		RequestID: h.RequestID,
		Status:    domain.RequestStatus(h.Status),
		Comment:   h.Comment,
		ChangedBy: h.ChangedBy,
		ChangedAt: h.ChangedAt,
	}
}
