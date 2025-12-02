package storage

import (
	"context"

	"github.com/google/uuid"

	"gateway/internal/domain"
)

const (
	requestAssignedSubject = "requests.responsible.assigned"
)

type requestRepository interface {
	AssignRequestResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID, userID *uuid.UUID) (domain.Request, error)
	GetRequestByID(ctx context.Context, id uuid.UUID) (domain.Request, error)
	ListRequests(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error)
}

type RequestStorage struct {
	repo      requestRepository
	publisher EventPublisher
}

func NewRequestStorage(repo requestRepository, publisher EventPublisher) *RequestStorage {
	return &RequestStorage{repo: repo, publisher: publisher}
}

func (s *RequestStorage) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID, userID *uuid.UUID) (domain.Request, error) {
	req, err := s.repo.AssignRequestResponsible(ctx, requestID, responsibleID, userID)
	if err != nil {
		return domain.Request{}, err
	}
	payload := struct {
		RequestID     uuid.UUID  `json:"request_id"`
		ResponsibleID *uuid.UUID `json:"responsible_id,omitempty"`
		UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
	}{RequestID: requestID, ResponsibleID: responsibleID, UpdatedBy: userID}
	if err := s.publish(ctx, requestAssignedSubject, payload); err != nil {
		return domain.Request{}, err
	}
	return req, nil
}

func (s *RequestStorage) Get(ctx context.Context, id uuid.UUID) (domain.Request, error) {
	return s.repo.GetRequestByID(ctx, id)
}

func (s *RequestStorage) List(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error) {
	return s.repo.ListRequests(ctx, responsibleID)
}

func (s *RequestStorage) publish(ctx context.Context, subject string, payload interface{}) error {
	if s.publisher == nil || subject == "" {
		return nil
	}
	return s.publisher.Publish(ctx, subject, payload)
}
