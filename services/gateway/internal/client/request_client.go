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

// RequestClient клиент для синхронных вызовов Request Service через NATS Request-Reply
type RequestClient struct {
	conn *broker.NATS
}

// NewRequestClient создает новый клиент для Request Service
func NewRequestClient(conn *broker.NATS) *RequestClient {
	return &RequestClient{conn: conn}
}

// List получает список заявок
func (c *RequestClient) List(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error) {
	req := dto.RequestListRequest{
		ResponsibleID: responsibleID,
	}

	data, err := req.MarshalJSON()
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

	var result []dto.Request
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal requests: %w", err)
	}

	requests := make([]domain.Request, 0, len(result))
	for _, r := range result {
		requests = append(requests, mapRequestFromDTO(r))
	}

	return requests, nil
}

// Get получает заявку по ID
func (c *RequestClient) Get(ctx context.Context, id uuid.UUID) (domain.Request, error) {
	req := dto.RequestByIDRequest{RequestID: id}

	data, err := req.MarshalJSON()
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

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Request{}, domain.ErrNotFound
		}
		return domain.Request{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.Request
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal request: %w", err)
	}

	return mapRequestFromDTO(result), nil
}

// AssignResponsible назначает ответственного за заявку
func (c *RequestClient) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID, userID *uuid.UUID) (domain.Request, error) {
	req := dto.AssignResponsibleRequest{
		RequestID: requestID,
	}
	if responsibleID != nil {
		req.ResponsibleID = responsibleID
	}

	data, err := req.MarshalJSON()
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

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.Request{}, domain.ErrNotFound
		}
		return domain.Request{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.Request
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal request: %w", err)
	}

	return mapRequestFromDTO(result), nil
}

func mapRequestFromDTO(r dto.Request) domain.Request {
	equip := make([]domain.RequestEquipment, 0, len(r.Equipment))
	for _, eq := range r.Equipment {
		equip = append(equip, domain.RequestEquipment{
			RequestID:   eq.RequestID,
			EquipmentID: eq.EquipmentID,
			Quantity:    eq.Quantity,
			CreatedAt:   eq.CreatedAt,
			UpdatedAt:   eq.UpdatedAt,
		})
	}

	return domain.Request{
		ID:                r.ID,
		TelegramUserInfo:  r.TelegramUserInfo,
		RequestText:       r.RequestText,
		Status:            domain.RequestStatus(r.Status),
		ScheduleTime:      r.ScheduleTime,
		EndTime:           r.EndTime,
		Address:           r.Address,
		ResponsibleUserID: r.ResponsibleUserID,
		Equipment:         equip,
		Audit: domain.AuditFields{
			CreatedAt: r.Audit.CreatedAt,
			UpdatedAt: r.Audit.UpdatedAt,
			DeletedAt: r.Audit.DeletedAt,
			CreatedBy: r.Audit.CreatedBy,
			UpdatedBy: r.Audit.UpdatedBy,
		},
	}
}
