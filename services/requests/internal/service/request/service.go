package request

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type Repository interface {
	SaveTelegramUser(ctx context.Context, user domain.User) (domain.User, error)
	CreateRequest(ctx context.Context, req domain.Request) (*domain.Request, error)
	GetRequestsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Request, error)
	GetRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.Request, error)
	UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error
	GetRequestsByResponsibleID(ctx context.Context, responsibleID uuid.UUID) ([]domain.Request, error)
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID) error
	ListRequests(ctx context.Context, limit, offset int32) ([]domain.Request, error)
}

// UserProvider интерфейс для получения информации о пользователе по ID
type UserProvider interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
}

type Service struct {
	repository   Repository
	userProvider UserProvider
}

func New(repository Repository, userProvider UserProvider) *Service {
	return &Service{
		repository:   repository,
		userProvider: userProvider,
	}
}

func (s *Service) Add(ctx context.Context, newRequest domain.Request) (*domain.Request, error) {
	issuer, err := s.repository.SaveTelegramUser(ctx, newRequest.Issuer) // Save or retrieve the user from db
	if err != nil {
		return nil, fmt.Errorf("add request: %w", err)
	}
	newRequest.Issuer = issuer

	request, err := s.repository.CreateRequest(ctx, newRequest)
	if err != nil {
		return nil, fmt.Errorf("add request: %w", err)
	}

	return request, nil
}

func (s *Service) List(ctx context.Context, user domain.User, limit, offset int32) ([]domain.Request, error) {
	user, err := s.repository.SaveTelegramUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get requests: %w", err)
	}

	requests, err := s.repository.GetRequestsByUserID(ctx, user.ID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get requests: %w", err)
	}

	return requests, nil
}

func (s *Service) Get(ctx context.Context, requestID uuid.UUID) (*domain.Request, error) {
	req, err := s.repository.GetRequestByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}

	return req, nil
}

func (s *Service) Cancel(ctx context.Context, requestID uuid.UUID) error {
	err := s.repository.UpdateRequestStatus(ctx, requestID, domain.StatusCanceled)
	if err != nil {
		return fmt.Errorf("cancel request: %w", err)
	}
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error {
	err := s.repository.UpdateRequestStatus(ctx, requestID, status)
	if err != nil {
		return fmt.Errorf("update request status: %w", err)
	}
	return nil
}

// ListByResponsible возвращает список заявок по ответственному.
// Если responsibleID nil, возвращает все заявки (для админа).
func (s *Service) ListByResponsible(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error) {
	if responsibleID == nil {
		// Для админа - возвращаем все заявки (без фильтрации по ответственному)
		requests, err := s.repository.ListRequests(ctx, 100, 0)
		if err != nil {
			return nil, fmt.Errorf("list requests: %w", err)
		}
		return requests, nil
	}

	requests, err := s.repository.GetRequestsByResponsibleID(ctx, *responsibleID)
	if err != nil {
		return nil, fmt.Errorf("list requests by responsible: %w", err)
	}

	return requests, nil
}

// AssignResponsible назначает ответственного за заявку.
// Если responsibleID nil, снимает назначение.
func (s *Service) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID) (*domain.Request, error) {
	err := s.repository.AssignResponsible(ctx, requestID, responsibleID)
	if err != nil {
		return nil, fmt.Errorf("assign responsible: %w", err)
	}

	return nil, nil
}
