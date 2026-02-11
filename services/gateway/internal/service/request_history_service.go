package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type RequestHistoryStorage interface {
	Add(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error)
	List(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error)
}

type RequestHistoryService struct {
	storage RequestHistoryStorage
}

func NewRequestHistoryService(storage RequestHistoryStorage) *RequestHistoryService {
	return &RequestHistoryService{storage: storage}
}

func (s *RequestHistoryService) Add(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.RequestStatusHistory{}, err
	}
	if s.storage == nil {
		return domain.RequestStatusHistory{}, errors.New("request history storage not implemented")
	}
	entry.ChangedBy = userIDFromCtx(ctx)
	return s.storage.Add(ctx, entry)
}

func (s *RequestHistoryService) List(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error) {
	if s.storage == nil {
		return nil, errors.New("request history storage not implemented")
	}
	return s.storage.List(ctx, requestID)
}
