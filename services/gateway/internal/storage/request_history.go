package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type requestHistoryRepository interface {
	AddRequestStatusHistory(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error)
	ListRequestStatusHistory(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error)
}

type RequestHistoryStorage struct {
	repo      requestHistoryRepository
	publisher EventPublisher
}

func NewRequestHistoryStorage(repo requestHistoryRepository, publisher EventPublisher) *RequestHistoryStorage {
	return &RequestHistoryStorage{repo: repo, publisher: publisher}
}

func (s *RequestHistoryStorage) Add(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error) {
	created, err := s.repo.AddRequestStatusHistory(ctx, entry)
	if err != nil {
		return domain.RequestStatusHistory{}, err
	}
	if err := s.publish(ctx, subjects.RequestStatusChanged, created); err != nil {
		return domain.RequestStatusHistory{}, err
	}
	return created, nil
}

func (s *RequestHistoryStorage) List(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error) {
	return s.repo.ListRequestStatusHistory(ctx, requestID)
}

func (s *RequestHistoryStorage) publish(ctx context.Context, subject string, payload interface{}) error {
	if s.publisher == nil || subject == "" {
		return nil
	}
	return s.publisher.Publish(ctx, subject, payload)
}
