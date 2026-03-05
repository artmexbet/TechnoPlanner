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
	GetRequestsByResponsibleID(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error)
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID) error
	ListRequests(ctx context.Context, limit, offset int32) ([]domain.Request, error)
	UpdateRequest(ctx context.Context, requestID uuid.UUID, updates domain.RequestUpdate) (*domain.Request, error)
	ListResponsibles(ctx context.Context) ([]domain.Porter, error)
	GetResponsible(ctx context.Context, id uuid.UUID) (domain.Porter, error)
	DeleteResponsible(ctx context.Context, id uuid.UUID) error
	SaveResponsible(ctx context.Context, id uuid.UUID, username string) error
	// RawRequests
	CreateRawRequest(ctx context.Context, req domain.RawRequest) (*domain.RawRequest, error)
	GetRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error)
	GetRawRequestByID(ctx context.Context, id uuid.UUID) (*domain.RawRequest, error)
	MarkRawRequestProcessed(ctx context.Context, id uuid.UUID, requestID uuid.UUID) (*domain.RawRequest, error)
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

	requests, err := s.repository.GetRequestsByResponsibleID(ctx, responsibleID)
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

// UpdateRequest обновляет заявку
func (s *Service) UpdateRequest(ctx context.Context, requestID uuid.UUID, updates domain.RequestUpdate) (*domain.Request, error) {
	return s.repository.UpdateRequest(ctx, requestID, updates)
}

// ListResponsibles возвращает список всех портеров
func (s *Service) ListResponsibles(ctx context.Context) ([]domain.Porter, error) {
	responsibles, err := s.repository.ListResponsibles(ctx)
	if err != nil {
		return nil, fmt.Errorf("list porters: %w", err)
	}

	result := make([]domain.Porter, len(responsibles))
	for i, r := range responsibles {
		result[i] = domain.Porter{
			ID:       r.ID,
			Username: r.Username,
		}
	}

	return result, nil
}

// SaveResponsible сохраняет или обновляет ответственного
func (s *Service) SaveResponsible(ctx context.Context, id uuid.UUID, username string) error {
	return s.repository.SaveResponsible(ctx, id, username)
}

// GetResponsible возвращает портера по ID
func (s *Service) GetResponsible(ctx context.Context, id uuid.UUID) (domain.Porter, error) {
	resp, err := s.repository.GetResponsible(ctx, id)
	if err != nil {
		return domain.Porter{}, fmt.Errorf("get porter: %w", err)
	}
	return resp, nil
}

// DeleteResponsible удаляет ответственного по ID
func (s *Service) DeleteResponsible(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteResponsible(ctx, id)
}

// CreateRawRequest сохраняет сырой запрос от Telegram-бота
func (s *Service) CreateRawRequest(ctx context.Context, req domain.RawRequest) (*domain.RawRequest, error) {
	created, err := s.repository.CreateRawRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create raw request: %w", err)
	}
	return created, nil
}

// ListRawRequests возвращает список сырых запросов с фильтрацией по статусу
func (s *Service) ListRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error) {
	requests, err := s.repository.GetRawRequests(ctx, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list raw requests: %w", err)
	}
	return requests, nil
}

// GetRawRequest возвращает сырой запрос по ID
func (s *Service) GetRawRequest(ctx context.Context, id uuid.UUID) (*domain.RawRequest, error) {
	req, err := s.repository.GetRawRequestByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get raw request: %w", err)
	}
	return req, nil
}

// ProcessRawRequest создаёт нормальную заявку на основе сырого запроса и помечает его как обработанный
func (s *Service) ProcessRawRequest(ctx context.Context, rawID uuid.UUID, newRequest domain.Request) (*domain.Request, *domain.RawRequest, error) {
	raw, err := s.repository.GetRawRequestByID(ctx, rawID)
	if err != nil {
		return nil, nil, fmt.Errorf("get raw request: %w", err)
	}
	if raw.Status == domain.RawRequestStatusProcessed {
		return nil, nil, fmt.Errorf("raw request already processed")
	}

	// Заполняем информацию об отправителе из сырого запроса, если не указана
	if newRequest.Issuer.TelegramID == 0 {
		newRequest.Issuer.TelegramID = raw.TelegramID
		newRequest.Issuer.Username = raw.Username
		newRequest.Issuer.FirstName = raw.FirstName
		newRequest.Issuer.LastName = raw.LastName
	}

	createdReq, err := s.Add(ctx, newRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("create request from raw: %w", err)
	}

	updatedRaw, err := s.repository.MarkRawRequestProcessed(ctx, rawID, createdReq.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("mark raw request processed: %w", err)
	}

	return createdReq, updatedRaw, nil
}
