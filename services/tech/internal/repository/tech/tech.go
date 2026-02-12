package techrepo

import (
	"context"

	"github.com/google/uuid"
	"tech/internal/domain"
)

type Postgres interface {
	AddEquipment(ctx context.Context, technic domain.Equipment) (*domain.Equipment, error)
	DeleteEquipment(ctx context.Context, techID uuid.UUID) error
	UpdateEquipment(ctx context.Context, technic domain.Equipment) (*domain.Equipment, error)
	GetEquipmentByID(ctx context.Context, techID uuid.UUID) (*domain.Equipment, error)
}

type Repository struct {
	pg Postgres
}

func NewRepository(pg Postgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) AddEquipment(ctx context.Context, technic domain.Equipment) (*domain.Equipment, error) {
	return r.pg.AddEquipment(ctx, technic)
}

func (r *Repository) DeleteEquipment(ctx context.Context, techID uuid.UUID) error {
	return r.pg.DeleteEquipment(ctx, techID)
}

func (r *Repository) UpdateEquipment(ctx context.Context, technic domain.Equipment) (*domain.Equipment, error) {
	return r.pg.UpdateEquipment(ctx, technic)
}

func (r *Repository) GetEquipmentByID(ctx context.Context, techID uuid.UUID) (*domain.Equipment, error) {
	return r.pg.GetEquipmentByID(ctx, techID)
}
