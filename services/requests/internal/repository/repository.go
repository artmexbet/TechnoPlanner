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
	GetRequestsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Request, error)
	UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error
	GetRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.Request, error)
	AssignEquipmentToRequest(ctx context.Context, requestID uuid.UUID, technics []domain.Equipment) []error
	GetEquipmentByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.Equipment, error)
	GetEquipmentByRequestIDs(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Equipment, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (domain.User, error)
	SaveTelegramUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetRequestsByResponsibleID(ctx context.Context, responsibleID uuid.UUID) ([]domain.Request, error)
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleInfo *domain.ResponsibleInfo) (*domain.Request, error)
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

func (r *Repository) GetRequestsByResponsibleID(ctx context.Context, responsibleID uuid.UUID) ([]domain.Request, error) {
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

func (r *Repository) AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleInfo *domain.ResponsibleInfo) (*domain.Request, error) {
	updatedReq, err := r.pg.AssignResponsible(ctx, requestID, responsibleInfo)
	if err != nil {
		return nil, fmt.Errorf("assign responsible: %w", err)
	}

	// Enrich with equipment and issuer info
	updatedReq.Issuer, err = r.pg.GetUserByID(ctx, updatedReq.Issuer.ID)
	if err != nil {
		return nil, fmt.Errorf("get issuer: %w", err)
	}
	updatedReq.Equipments, err = r.pg.GetEquipmentByRequestID(ctx, updatedReq.ID)
	if err != nil {
		return nil, fmt.Errorf("get equipment: %w", err)
	}

	return updatedReq, nil
}
