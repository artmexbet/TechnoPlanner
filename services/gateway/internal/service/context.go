package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type contextKey string

const (
	contextUserRole contextKey = "user_role"
	contextUserID   contextKey = "user_id"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RolePorter Role = "porter"
)

func roleFromCtx(ctx context.Context) Role {
	if v, ok := ctx.Value(contextUserRole).(string); ok {
		return Role(v)
	}
	return ""
}

func userIDFromCtx(ctx context.Context) *uuid.UUID {
	if v, ok := ctx.Value(contextUserID).(string); ok {
		id, err := uuid.Parse(v)
		if err == nil {
			return &id
		}
	}
	return nil
}

func requireAdmin(ctx context.Context) error {
	if roleFromCtx(ctx) != RoleAdmin {
		return domain.ErrForbidden
	}
	return nil
}

func WithUserContext(ctx context.Context, userID, role string) context.Context {
	if userID != "" {
		ctx = context.WithValue(ctx, contextUserID, userID)
	}
	if role != "" {
		ctx = context.WithValue(ctx, contextUserRole, role)
	}
	return ctx
}
