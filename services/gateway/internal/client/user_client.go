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

// UserClient клиент для синхронных вызовов User Service через NATS Request-Reply
type UserClient struct {
	conn *broker.NATS
}

// NewUserClient создает новый клиент для User Service
func NewUserClient(conn *broker.NATS) *UserClient {
	return &UserClient{conn: conn}
}

// Get получает пользователя по ID (для интерфейса PorterStorage)
func (c *UserClient) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return c.GetUserByID(ctx, id)
}

// List получает список портеров по roleID (для интерфейса PorterStorage)
func (c *UserClient) List(ctx context.Context, roleID int32) ([]domain.User, error) {
	req := dto.RoleIDRequest{RoleID: roleID}

	data, err := req.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayUserList, data)
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

	var result []dto.User
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal user list: %w", err)
	}

	users := make([]domain.User, 0, len(result))
	for _, u := range result {
		users = append(users, mapUserFromDTO(u))
	}

	return users, nil
}

// GetUserByID получает пользователя по ID
func (c *UserClient) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	req := dto.UUIDRequest{ID: id}

	data, err := req.MarshalJSON()
	if err != nil {
		return domain.User{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayUserGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.User{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.User
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.User{}, fmt.Errorf("unmarshal user: %w", err)
	}

	return mapUserFromDTO(result), nil
}

// CreateUser создает нового пользователя
func (c *UserClient) CreateUser(ctx context.Context, user domain.User) error {
	req := dto.UserCreateRequest{
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
	}

	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayUserCreate, data)
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

// Update обновляет пользователя по ID
func (c *UserClient) Update(ctx context.Context, id uuid.UUID, username, email string) (domain.User, error) {
	req := dto.UserUpdateRequest{
		ID:       id,
		Username: username,
		Email:    email,
	}

	data, err := req.MarshalJSON()
	if err != nil {
		return domain.User{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayUserUpdate, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.User{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message == "not found" {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.User
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.User{}, fmt.Errorf("unmarshal user: %w", err)
	}

	return mapUserFromDTO(result), nil
}

// Delete удаляет пользователя по ID
func (c *UserClient) Delete(ctx context.Context, id uuid.UUID) error {
	req := dto.UserDeleteRequest{ID: id}

	data, err := req.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayUserDelete, data)
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

func mapUserFromDTO(u dto.User) domain.User {
	return domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}
