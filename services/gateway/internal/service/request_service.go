package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

// ErrForbidden deprecated: use domain.ErrForbidden
var ErrForbidden = domain.ErrForbidden

// ErrNotFound deprecated: use domain.ErrNotFound
var ErrNotFound = domain.ErrNotFound

type RequestStorage interface {
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID, userID *uuid.UUID) (domain.Request, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Request, error)
	List(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error)
}

type RequestService struct {
	storage RequestStorage
}

func NewRequestService(storage RequestStorage) *RequestService {
	return &RequestService{storage: storage}
}

func (s *RequestService) List(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error) {
	if responsibleID != nil && roleFromCtx(ctx) != RoleAdmin {
		userID := userIDFromCtx(ctx)
		if userID == nil || userID.String() != responsibleID.String() {
			return nil, ErrForbidden
		}
	}
	return s.storage.List(ctx, responsibleID)
}

func (s *RequestService) Get(ctx context.Context, id uuid.UUID) (domain.Request, error) {
	req, err := s.storage.Get(ctx, id)
	if err != nil {
		return domain.Request{}, err
	}
	if roleFromCtx(ctx) == RoleAdmin {
		return req, nil
	}
	if req.ResponsibleUserID != nil {
		userID := userIDFromCtx(ctx)
		if userID != nil && req.ResponsibleUserID.String() == userID.String() {
			return req, nil
		}
	}
	return domain.Request{}, ErrForbidden
}

func (s *RequestService) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID uuid.UUID) (domain.Request, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.Request{}, err
	}
	return s.storage.AssignResponsible(ctx, requestID, &responsibleID, userIDFromCtx(ctx))
}
