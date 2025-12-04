package server

import (
	"context"
	"fmt"

	"github.com/artmexbet/TechnoPlanner/libs/proto"

	"github.com/artmexbet/TechnoPlanner/services/auth/internal/models"
)

type authService interface {
	Login(ctx context.Context, loginRequest models.LoginRequest) (models.TokenPair, error)
	Register(ctx context.Context, req models.RegisterRequest) (models.User, error)
	RegisterPorter(ctx context.Context, req models.RegisterRequest) (models.User, error)
	ValidateToken(ctx context.Context, token string) (models.TokenValidateResult, error)
	Refresh(ctx context.Context, req models.TokenRefreshRequest) (models.TokenPair, error)
	Logout(ctx context.Context, token string) error
	LogoutAll(ctx context.Context, token string) error
}

type Handler struct {
	proto.UnimplementedAuthServer

	svc authService
}

func NewHandler(svc authService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(ctx context.Context, in *proto.LoginRequest) (*proto.TokenPair, error) {
	pair, err := h.svc.Login(ctx, *models.UserLoginFromProto(in))
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	return pair.ToProto(), nil
}

func (h *Handler) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	u, err := h.svc.Register(ctx, *models.UserRegisterFromProto(in))
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	resp := &proto.RegisterResponse{
		UserId: u.ID.String(),
	}
	return resp, nil
}

func (h *Handler) RegisterPorter(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	u, err := h.svc.RegisterPorter(ctx, *models.UserRegisterFromProto(in))
	if err != nil {
		return nil, fmt.Errorf("register porter: %w", err)
	}
	resp := &proto.RegisterResponse{
		UserId: u.ID.String(),
	}
	return resp, nil
}

func (h *Handler) Validate(ctx context.Context, in *proto.TokenRequest) (*proto.ValidateResponse, error) {
	r, err := h.svc.ValidateToken(ctx, in.Token)
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}
	resp := &proto.ValidateResponse{
		State:  proto.TokenState(r.State),
		UserId: r.UserID,
		Role:   r.Role,
	}
	return resp, nil
}

func (h *Handler) Refresh(ctx context.Context, in *proto.RefreshRequest) (*proto.TokenPair, error) {
	pair, err := h.svc.Refresh(ctx, models.TokenRefreshRequest{
		Pair:      models.TokenPairFromProto(in.GetPair()),
		DeviceID:  in.GetDeviceId(),
		UserAgent: in.GetUserAgent(),
		IP:        in.GetIpAddress(),
	})
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}
	return pair.ToProto(), nil
}

func (h *Handler) Logout(ctx context.Context, in *proto.TokenRequest) (*proto.Empty, error) {
	err := h.svc.Logout(ctx, in.Token)
	if err != nil {
		return nil, fmt.Errorf("logout: %w", err)
	}
	return &proto.Empty{}, nil
}

func (h *Handler) LogoutAll(ctx context.Context, in *proto.TokenRequest) (*proto.Empty, error) {
	err := h.svc.LogoutAll(ctx, in.Token)
	if err != nil {
		return nil, fmt.Errorf("logout all: %w", err)
	}
	return &proto.Empty{}, nil
}
