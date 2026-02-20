package categoryrepo

import (
	"context"

	"github.com/google/uuid"
	"tech/internal/domain"
)

type Postgres interface {
	GetEquipmentByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Equipment, error)
	AddCategory(ctx context.Context, categoryName string) (*domain.EquipmentCategory, error)
	UpdateCategoryName(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error)
}

type Repository struct {
	pg Postgres
}

func NewRepository(pg Postgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) GetEquipmentByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Equipment, error) {
	return r.pg.GetEquipmentByCategory(ctx, categoryID)
}

func (r *Repository) AddCategory(ctx context.Context, categoryName string) (*domain.EquipmentCategory, error) {
	return r.pg.AddCategory(ctx, categoryName)
}

func (r *Repository) UpdateCategoryName(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error) {
	return r.pg.UpdateCategoryName(ctx, category)
}

func (r *Repository) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return r.pg.DeleteCategory(ctx, categoryID)
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error) {
	return r.pg.GetAllCategories(ctx)
}
