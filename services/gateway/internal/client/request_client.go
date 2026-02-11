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
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

// RequestClient клиент для синхронных вызовов Request Service через NATS Request-Reply
type RequestClient struct {
	conn *broker.NATS
}

// NewRequestClient создает новый клиент для Request Service
func NewRequestClient(conn *broker.NATS) *RequestClient {
	return &RequestClient{conn: conn}
}

// requestListRequest запрос на получение списка заявок
type requestListRequest struct {
	ResponsibleID *string `json:"responsible_id,omitempty"`
}

// requestByIDRequest запрос на получение заявки по ID
type requestByIDRequest struct {
	RequestID string `json:"request_id"`
}

// assignResponsibleRequest запрос на назначение ответственного
type assignResponsibleRequest struct {
	RequestID     string `json:"request_id"`
	ResponsibleID string `json:"responsible_id"`
}

// gatewayRequestResponse стандартный ответ от сервисов
type gatewayRequestResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// List получает список заявок
func (c *RequestClient) List(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error) {
	req := requestListRequest{}
	if responsibleID != nil {
		id := responsibleID.String()
		req.ResponsibleID = &id
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayRequestList, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("nats request: %w", err)
	}

	var resp gatewayRequestResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("service error: %s", resp.Message)
	}

	var requests []domain.Request
	if err := json.Unmarshal(resp.Data, &requests); err != nil {
		return nil, fmt.Errorf("unmarshal requests: %w", err)
	}

	return requests, nil
}

// Get получает заявку по ID
func (c *RequestClient) Get(ctx context.Context, id uuid.UUID) (domain.Request, error) {
	req := requestByIDRequest{RequestID: id.String()}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.Request{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayRequestGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Request{}, domain.ErrNotFound
		}
		return domain.Request{}, fmt.Errorf("nats request: %w", err)
	}

	var resp gatewayRequestResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Request{}, domain.ErrNotFound
		}
		return domain.Request{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var request domain.Request
	if err := json.Unmarshal(resp.Data, &request); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal request: %w", err)
	}

	return request, nil
}

// AssignResponsible назначает ответственного за заявку
func (c *RequestClient) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID, userID *uuid.UUID) (domain.Request, error) {
	req := assignResponsibleRequest{
		RequestID:     requestID.String(),
		ResponsibleID: responsibleID.String(),
	}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.Request{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayRequestAssignResponsible, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Request{}, domain.ErrNotFound
		}
		return domain.Request{}, fmt.Errorf("nats request: %w", err)
	}

	var resp gatewayRequestResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Request{}, domain.ErrNotFound
		}
		return domain.Request{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var request domain.Request
	if err := json.Unmarshal(resp.Data, &request); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal request: %w", err)
	}

	return request, nil
}
