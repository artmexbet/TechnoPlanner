package service

import (
	"context"

	"github.com/google/uuid"

	"gateway/internal/domain"
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
	entry.ChangedBy = userIDFromCtx(ctx)
	return s.storage.Add(ctx, entry)
}

func (s *RequestHistoryService) List(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error) {
	return s.storage.List(ctx, requestID)
}
