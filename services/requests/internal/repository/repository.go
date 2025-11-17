package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"requests/internal/domain"
)

type iPostgres interface {
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
}

// Repository struct that interacts with the databases.
// We will be in need of sending data to broker in the future, so having a repository wrapper is a good idea.
type Repository struct {
	pg iPostgres
}

func NewRepository(pg iPostgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) CreateRequest(ctx context.Context, req domain.Request) (*domain.Request, error) {
	return r.pg.CreateRequest(ctx, req)
}

func (r *Repository) CreateEquipment(ctx context.Context, technics []domain.Equipment) error {
	return r.pg.CreateEquipment(ctx, technics)
}

func (r *Repository) GetRequestsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Request, error) {
	return r.pg.GetRequestsByUserID(ctx, userID, limit, offset)
}

func (r *Repository) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error {
	return r.pg.UpdateRequestStatus(ctx, requestID, status)
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
	req.Technics, err = r.pg.GetEquipmentByRequestID(ctx, requestID)
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
	return r.pg.SaveTelegramUser(ctx, user)
}
