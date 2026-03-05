package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type Postgres interface {
	CreateRequest(ctx context.Context, req domain.Request) (*domain.Request, error)
	CreateEquipment(ctx context.Context, technics []domain.Equipment) error
	UpsertEquipment(ctx context.Context, eq domain.Equipment) error
	DeleteEquipment(ctx context.Context, id int) error
	GetRequestsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Request, error)
	UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error
	GetRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.Request, error)
	AssignEquipmentToRequest(ctx context.Context, requestID uuid.UUID, technics []domain.Equipment) []error
	GetEquipmentByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.Equipment, error)
	GetEquipmentByRequestIDs(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Equipment, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (domain.User, error)
	SaveTelegramUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetRequestsByResponsibleID(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error)
	ListRequests(ctx context.Context, offset, limit int32) ([]domain.Request, error)
	// Porters
	ListResponsibles(ctx context.Context) ([]domain.Porter, error)
	GetResponsible(ctx context.Context, id uuid.UUID) (domain.Porter, error)
	DeleteResponsible(ctx context.Context, id uuid.UUID) error
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID) error
	UpdateRequest(ctx context.Context, requestID uuid.UUID, updates domain.RequestUpdate) (*domain.Request, error)
	SaveResponsible(ctx context.Context, id uuid.UUID, username string) error
	// RawRequests
	CreateRawRequest(ctx context.Context, req domain.RawRequest) (*domain.RawRequest, error)
	GetRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error)
	GetRawRequestByID(ctx context.Context, id uuid.UUID) (*domain.RawRequest, error)
	MarkRawRequestProcessed(ctx context.Context, id uuid.UUID, requestID uuid.UUID) (*domain.RawRequest, error)
}

type Publisher interface {
	PublishRequestCreated(req domain.Request) error
	PublishRequestCanceled(req domain.Request) error
	PublishUserAdded(user domain.User) error
}

// Repository struct that interacts with the databases.
// We will be in need of sending data to broker in the future, so having a repository wrapper is a good idea.
type Repository struct {
	pg        Postgres
	publisher Publisher
}

func NewRepository(pg Postgres, publisher Publisher) *Repository {
	return &Repository{pg: pg, publisher: publisher}
}

func (r *Repository) CreateRequest(ctx context.Context, req domain.Request) (*domain.Request, error) {
	newReq, err := r.pg.CreateRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	// Publish event
	err = r.publisher.PublishRequestCreated(*newReq)
	if err != nil {
		return nil, fmt.Errorf("error publishing request created event: %w", err)
	}
	return newReq, nil
}

func (r *Repository) CreateEquipment(ctx context.Context, technics []domain.Equipment) error {
	return r.pg.CreateEquipment(ctx, technics)
}

func (r *Repository) UpsertEquipment(ctx context.Context, eq domain.Equipment) error {
	return r.pg.UpsertEquipment(ctx, eq)
}

func (r *Repository) DeleteEquipment(ctx context.Context, id int) error {
	return r.pg.DeleteEquipment(ctx, id)
}

func (r *Repository) GetRequestsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Request, error) {
	return r.pg.GetRequestsByUserID(ctx, userID, limit, offset)
}

func (r *Repository) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error {
	err := r.pg.UpdateRequestStatus(ctx, requestID, status)
	if err != nil {
		return fmt.Errorf("error updating request status: %w", err)
	}
	if status == domain.StatusCanceled { // only canceled status requires event publishing
		req, err := r.pg.GetRequestByID(ctx, requestID)
		if err != nil {
			return fmt.Errorf("get request for cancelation event: %w", err)
		}
		// Publish event
		err = r.publisher.PublishRequestCanceled(*req)
		if err != nil {
			return fmt.Errorf("error publishing request canceled event: %w", err)
		}
	}
	return nil
}

// GetRequestByID retrieves a request by its ID, including its issuer and associated technics.
func (r *Repository) GetRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.Request, error) {
	req, err := r.pg.GetRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	req.Issuer, err = r.pg.GetUserByID(ctx, req.Issuer.ID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	req.Equipments, err = r.pg.GetEquipmentByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("get technics: %w", err)
	}
	return req, nil
}
func (r *Repository) AssignEquipmentToRequest(ctx context.Context, requestID uuid.UUID, equipment []domain.Equipment) []error {
	return r.pg.AssignEquipmentToRequest(ctx, requestID, equipment)
}
func (r *Repository) GetEquipmentByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.Equipment, error) {
	return r.pg.GetEquipmentByRequestID(ctx, requestID)
}
func (r *Repository) GetEquipmentByRequestIDs(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Equipment, error) {
	return r.pg.GetEquipmentByRequestIDs(ctx, requestIDs)
}

func (r *Repository) SaveTelegramUser(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := r.pg.SaveTelegramUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("save telegram user: %w", err)
	}
	// Publish event
	err = r.publisher.PublishUserAdded(u)
	if err != nil {
		return domain.User{}, fmt.Errorf("error publishing user added event: %w", err)
	}
	return u, nil
}

func (r *Repository) GetRequestsByResponsibleID(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error) {
	requests, err := r.pg.GetRequestsByResponsibleID(ctx, responsibleID)
	if err != nil {
		return nil, fmt.Errorf("get requests by responsible id: %w", err)
	}

	// Enrich requests with equipment and issuer info
	for i := range requests {
		// todo: use batch queries
		requests[i].Issuer, err = r.pg.GetUserByID(ctx, requests[i].Issuer.ID)
		if err != nil {
			return nil, fmt.Errorf("get issuer for request %s: %w", requests[i].ID, err)
		}
		requests[i].Equipments, err = r.pg.GetEquipmentByRequestID(ctx, requests[i].ID)
		if err != nil {
			return nil, fmt.Errorf("get equipment for request %s: %w", requests[i].ID, err)
		}
	}

	return requests, nil
}

func (r *Repository) ListRequests(ctx context.Context, offset, limit int32) ([]domain.Request, error) {
	return r.pg.ListRequests(ctx, offset, limit)
}

// ListResponsibles возвращает список всех портеров
func (r *Repository) ListResponsibles(ctx context.Context) ([]domain.Porter, error) {
	return r.pg.ListResponsibles(ctx)
}

// GetResponsible возвращает портера по ID
func (r *Repository) GetResponsible(ctx context.Context, id uuid.UUID) (domain.Porter, error) {
	return r.pg.GetResponsible(ctx, id)
}

// DeleteResponsible удаляет ответственного по ID
func (r *Repository) DeleteResponsible(ctx context.Context, id uuid.UUID) error {
	return r.pg.DeleteResponsible(ctx, id)
}

// AssignResponsible назначает ответственного за заявку
func (r *Repository) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID) error {
	return r.pg.AssignResponsible(ctx, requestID, responsibleID)
}

// UpdateRequest обновляет заявку
func (r *Repository) UpdateRequest(ctx context.Context, requestID uuid.UUID, updates domain.RequestUpdate) (*domain.Request, error) {
	return r.pg.UpdateRequest(ctx, requestID, updates)
}

// SaveResponsible сохраняет или обновляет ответственного
func (r *Repository) SaveResponsible(ctx context.Context, id uuid.UUID, username string) error {
	return r.pg.SaveResponsible(ctx, id, username)
}

// CreateRawRequest сохраняет сырой запрос от Telegram-бота
func (r *Repository) CreateRawRequest(ctx context.Context, req domain.RawRequest) (*domain.RawRequest, error) {
	return r.pg.CreateRawRequest(ctx, req)
}

// GetRawRequests возвращает список сырых запросов
func (r *Repository) GetRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error) {
	return r.pg.GetRawRequests(ctx, status, limit, offset)
}

// GetRawRequestByID возвращает сырой запрос по ID
func (r *Repository) GetRawRequestByID(ctx context.Context, id uuid.UUID) (*domain.RawRequest, error) {
	return r.pg.GetRawRequestByID(ctx, id)
}

// MarkRawRequestProcessed помечает сырой запрос как обработанный
func (r *Repository) MarkRawRequestProcessed(ctx context.Context, id uuid.UUID, requestID uuid.UUID) (*domain.RawRequest, error) {
	return r.pg.MarkRawRequestProcessed(ctx, id, requestID)
}
