package natssubscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/auth/internal/models"
)

// UserRepository интерфейс для операций с пользователями
type UserRepository interface {
	FindUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, username, email string) (models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// Subscriber NATS subscriber для auth service
type Subscriber struct {
	conn          *nats.Conn
	userRepo      UserRepository
	subscriptions []*nats.Subscription
}

// NewSubscriber создает новый NATS subscriber
func NewSubscriber(cfg config.NATSConfig, userRepo UserRepository) (*Subscriber, error) {
	nc, err := nats.Connect(
		cfg.URL(),
		nats.Name("auth-service-subscriber"),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats: %w", err)
	}
	return &Subscriber{
		conn:     nc,
		userRepo: userRepo,
	}, nil
}

// Close закрывает соединение
func (s *Subscriber) Close() {
	for _, sub := range s.subscriptions {
		_ = sub.Unsubscribe()
	}
	s.conn.Close()
}

// HandleMsgs регистрирует обработчики NATS сообщений
func (s *Subscriber) HandleMsgs() *Subscriber {
	handlers := map[string]nats.MsgHandler{
		subjects.GatewayUserGet:    s.handleGetUser,
		subjects.GatewayUserList:   s.handleListUsers,
		subjects.GatewayUserUpdate: s.handleUpdateUser,
		subjects.GatewayUserDelete: s.handleDeleteUser,
	}

	for subject, handler := range handlers {
		sub, err := s.conn.Subscribe(subject, handler)
		if err != nil {
			slog.Error("cannot subscribe to subject", "subject", subject, "error", err)
		} else {
			s.subscriptions = append(s.subscriptions, sub)
		}
	}
	return s
}

func respondSuccess(msg *nats.Msg, data interface{}) {
	raw, _ := json.Marshal(data)
	resp := dto.GatewayResponse{
		Success: true,
		Message: "success",
		Data:    json.RawMessage(raw),
	}
	respData, _ := json.Marshal(resp)
	_ = msg.Respond(respData)
}

func respondError(msg *nats.Msg, message string) {
	resp := dto.GatewayResponse{
		Success: false,
		Message: message,
	}
	respData, _ := json.Marshal(resp)
	_ = msg.Respond(respData)
}

func (s *Subscriber) handleGetUser(msg *nats.Msg) {
	ctx := context.Background()

	var req dto.UUIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get user request", "error", err)
		respondError(msg, "invalid request format")
		return
	}

	user, err := s.userRepo.FindUserByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error finding user", "error", err, "id", req.ID)
		respondError(msg, "not found")
		return
	}

	respondSuccess(msg, dto.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (s *Subscriber) handleListUsers(msg *nats.Msg) {
	// For now, this is not fully implemented on auth side.
	// Returns empty list - actual list comes from user-specific queries.
	respondSuccess(msg, []dto.User{})
}

func (s *Subscriber) handleUpdateUser(msg *nats.Msg) {
	ctx := context.Background()

	var req dto.UserUpdateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update user request", "error", err)
		respondError(msg, "invalid request format")
		return
	}

	user, err := s.userRepo.UpdateUser(ctx, req.ID, req.Username, req.Email)
	if err != nil {
		slog.ErrorContext(ctx, "error updating user", "error", err, "id", req.ID)
		respondError(msg, "internal server error")
		return
	}

	respondSuccess(msg, dto.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (s *Subscriber) handleDeleteUser(msg *nats.Msg) {
	ctx := context.Background()

	var req dto.UserDeleteRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling delete user request", "error", err)
		respondError(msg, "invalid request format")
		return
	}

	if err := s.userRepo.DeleteUser(ctx, req.ID); err != nil {
		slog.ErrorContext(ctx, "error deleting user", "error", err, "id", req.ID)
		respondError(msg, "internal server error")
		return
	}

	respondSuccess(msg, nil)
}
