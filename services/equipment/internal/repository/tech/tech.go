package techrepo

import (
	"context"

	"tech/internal/domain"
)

type Postgres interface {
	AddEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error)
	DeleteEquipment(ctx context.Context, equipmentID int) error
	UpdateEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error)
	GetEquipmentByID(ctx context.Context, equipmentID int) (*domain.Equipment, error)
	GetAllEquipment(ctx context.Context) ([]domain.Equipment, error)
	GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error)
	ReserveEquipment(ctx context.Context, items []domain.ReserveItem) error
	ReleaseEquipment(ctx context.Context, items []domain.ReserveItem) error
	CheckAvailability(ctx context.Context, items []domain.ReserveItem) (bool, []int, error)
}

type Repository struct {
	pg Postgres
}

func NewRepository(pg Postgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) AddEquipment(ctx context.Context, eq domain.Equipment) (*domain.Equipment, error) {
	return r.pg.AddEquipment(ctx, eq)
}

func (r *Repository) DeleteEquipment(ctx context.Context, equipmentID int) error {
	return r.pg.DeleteEquipment(ctx, equipmentID)
}

func (r *Repository) UpdateEquipment(ctx context.Context, eq domain.Equipment) (*domain.Equipment, error) {
	return r.pg.UpdateEquipment(ctx, eq)
}

func (r *Repository) GetEquipmentByID(ctx context.Context, equipmentID int) (*domain.Equipment, error) {
	return r.pg.GetEquipmentByID(ctx, equipmentID)
}

func (r *Repository) GetAllEquipment(ctx context.Context) ([]domain.Equipment, error) {
	return r.pg.GetAllEquipment(ctx)
}

func (r *Repository) GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error) {
	return r.pg.GetEquipmentByCategory(ctx, categoryID)
}

func (r *Repository) ReserveEquipment(ctx context.Context, items []domain.ReserveItem) error {
	return r.pg.ReserveEquipment(ctx, items)
}

func (r *Repository) ReleaseEquipment(ctx context.Context, items []domain.ReserveItem) error {
	return r.pg.ReleaseEquipment(ctx, items)
}

func (r *Repository) CheckAvailability(ctx context.Context, items []domain.ReserveItem) (bool, []int, error) {
	return r.pg.CheckAvailability(ctx, items)
}
