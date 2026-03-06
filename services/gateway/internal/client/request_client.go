package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
func (c *RequestClient) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID, _ *uuid.UUID) (domain.Request, error) {
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

// UpdateRequest обновляет заявку
func (c *RequestClient) UpdateRequest(ctx context.Context, requestID uuid.UUID, updates domain.RequestUpdate) (domain.Request, error) {
	var status *dto.RequestStatus
	if updates.Status != nil {
		s := dto.RequestStatus(*updates.Status)
		status = &s
	}

	req := dto.RequestUpdateRequest{
		RequestID:     requestID,
		RequestText:   updates.RequestText,
		Status:        status,
		ScheduleTime:  updates.ScheduleTime,
		EndTime:       updates.EndTime,
		Address:       updates.Address,
		ResponsibleID: updates.ResponsibleID,
	}

	data, err := req.MarshalJSON()
	if err != nil {
		return domain.Request{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayRequestUpdate, data)
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

// ListRawRequests возвращает список сырых запросов от бота
func (c *RequestClient) ListRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error) {
	req := dto.RawRequestListRequest{
		Status: status,
		Limit:  limit,
		Offset: offset,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayRawRequestList, data)
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

	var result []dto.RawRequest
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal raw requests: %w", err)
	}

	requests := make([]domain.RawRequest, 0, len(result))
	for _, r := range result {
		requests = append(requests, mapRawRequestFromDTO(r))
	}
	return requests, nil
}

// GetRawRequest возвращает сырой запрос по ID
func (c *RequestClient) GetRawRequest(ctx context.Context, id uuid.UUID) (domain.RawRequest, error) {
	req := dto.UUIDRequest{ID: id}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.RawRequest{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayRawRequestGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.RawRequest{}, domain.ErrNotFound
		}
		return domain.RawRequest{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.RawRequest{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.RawRequest{}, domain.ErrNotFound
		}
		return domain.RawRequest{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.RawRequest
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.RawRequest{}, fmt.Errorf("unmarshal raw request: %w", err)
	}

	return mapRawRequestFromDTO(result), nil
}

// ProcessRawRequest создаёт нормальную заявку из сырого запроса
func (c *RequestClient) ProcessRawRequest(ctx context.Context, req dto.RawRequestProcessRequest) (domain.Request, domain.RawRequest, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return domain.Request{}, domain.RawRequest{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayRawRequestProcess, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.Request{}, domain.RawRequest{}, domain.ErrNotFound
		}
		return domain.Request{}, domain.RawRequest{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.Request{}, domain.RawRequest{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		return domain.Request{}, domain.RawRequest{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result struct {
		Request    dto.Request    `json:"request"`
		RawRequest dto.RawRequest `json:"raw_request"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.Request{}, domain.RawRequest{}, fmt.Errorf("unmarshal process response: %w", err)
	}

	return mapRequestFromDTO(result.Request), mapRawRequestFromDTO(result.RawRequest), nil
}

func mapRequestFromDTO(r dto.Request) domain.Request {
	equip := make([]domain.RequestEquipment, 0, len(r.Equipment))
	for _, eq := range r.Equipment {
		equip = append(equip, domain.RequestEquipment{
			RequestID:   eq.RequestID,
			EquipmentID: eq.EquipmentID,
			Name:        eq.Name,
			Description: eq.Description,
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

func mapRawRequestFromDTO(r dto.RawRequest) domain.RawRequest {
	result := domain.RawRequest{
		TelegramID: r.TelegramID,
		Username:   r.Username,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		RawText:    r.RawText,
		Status:     r.Status,
		CreatedAt:  parseTime(r.CreatedAt),
	}
	if id, err := uuid.Parse(r.ID); err == nil {
		result.ID = id
	}
	if r.ProcessedRequestID != nil {
		if id, err := uuid.Parse(*r.ProcessedRequestID); err == nil {
			result.ProcessedRequestID = &id
		}
	}
	return result
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
