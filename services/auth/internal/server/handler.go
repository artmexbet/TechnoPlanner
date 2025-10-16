package server

import (
	"auth/internal/models"
	"proto"

	"context"
	"fmt"
)

type authService interface {
	Login(ctx context.Context, loginRequest models.LoginRequest) (models.TokenPair, error)
	Register(ctx context.Context, req models.RegisterRequest) (models.User, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type Handler struct {
	proto.UnimplementedAuthServer

	svc authService
}

func NewHandler(svc authService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	pair, err := h.svc.Login(ctx, *models.UserLoginFromProto(in))
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	resp := &proto.LoginResponse{
		Token:        pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
	return resp, nil
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
