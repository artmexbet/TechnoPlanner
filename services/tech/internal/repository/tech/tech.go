package techrepo

import (
	"context"
	"tech/internal/domain"

	"github.com/google/uuid"
)

type IPostgres interface {
	AddTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error)
	DeleteTechnic(ctx context.Context, techID uuid.UUID) error
	UpdateTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error)
	GetTechnicByID(ctx context.Context, techID uuid.UUID) (*domain.Technic, error)
}

type Repository struct {
	pg IPostgres
}

func NewRepository(pg IPostgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) AddTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error) {
	return r.pg.AddTechnic(ctx, technic)
}

func (r *Repository) DeleteTechnic(ctx context.Context, techID uuid.UUID) error {
	return r.pg.DeleteTechnic(ctx, techID)
}

func (r *Repository) UpdateTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error) {
	return r.pg.UpdateTechnic(ctx, technic)
}

func (r *Repository) GetTechnicByID(ctx context.Context, techID uuid.UUID) (*domain.Technic, error) {
	return r.pg.GetTechnicByID(ctx, techID)
}
