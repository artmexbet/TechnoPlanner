package server

import (
	"context"
	"log/slog"

	"proto"
)

type authService interface {
	Login(username, password string) (string, error)
	Register(username, password string) error
	ValidateToken(token string) (string, error)
}

type Handler struct {
	proto.UnimplementedAuthServer

	svc authService
}

func NewHandler(svc authService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	slog.Info("Login called", "username", in.Username)
	return nil, nil
}

func (h *Handler) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	slog.Info("Register called", "username", in.Username)
	return nil, nil
}
