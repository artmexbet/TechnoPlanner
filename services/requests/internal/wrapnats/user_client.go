package wrapnats

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

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

// UserClient клиент для получения информации о пользователях через NATS
type UserClient struct {
	conn *broker.NATS
}

// NewUserClient создает новый клиент для User Service
func NewUserClient(conn *broker.NATS) *UserClient {
	return &UserClient{conn: conn}
}

// GetUserByID получает пользователя по ID
func (c *UserClient) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	req := dto.UUIDRequest{ID: userID}

	data, err := req.MarshalJSON()
	if err != nil {
		return domain.User{}, fmt.Errorf("marshal request: %w", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, subjects.GatewayUserGet, data)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return domain.User{}, fmt.Errorf("user service not available: %w", err)
		}
		return domain.User{}, fmt.Errorf("nats request: %w", err)
	}

	var resp dto.GatewayResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.User{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		return domain.User{}, fmt.Errorf("service error: %s", resp.Message)
	}

	var result dto.User
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return domain.User{}, fmt.Errorf("unmarshal user: %w", err)
	}

	return mapUserFromDTO(result), nil
}

func mapUserFromDTO(u dto.User) domain.User {
	return domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
