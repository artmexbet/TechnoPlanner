package categoryrepo

import (
	"context"

	"tech/internal/domain"
)

type Postgres interface {
	GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error)
	AddCategory(ctx context.Context, categoryName string, description string) (*domain.EquipmentCategory, error)
	UpdateCategory(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error)
	DeleteCategory(ctx context.Context, categoryID int) error
	GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error)
	GetCategoryByID(ctx context.Context, categoryID int) (*domain.EquipmentCategory, error)
}

type Repository struct {
	pg Postgres
}

func NewRepository(pg Postgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error) {
	return r.pg.GetEquipmentByCategory(ctx, categoryID)
}

func (r *Repository) AddCategory(ctx context.Context, categoryName string, description string) (*domain.EquipmentCategory, error) {
	return r.pg.AddCategory(ctx, categoryName, description)
}

func (r *Repository) UpdateCategory(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error) {
	return r.pg.UpdateCategory(ctx, category)
}

func (r *Repository) DeleteCategory(ctx context.Context, categoryID int) error {
	return r.pg.DeleteCategory(ctx, categoryID)
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error) {
	return r.pg.GetAllCategories(ctx)
}

func (r *Repository) GetCategoryByID(ctx context.Context, categoryID int) (*domain.EquipmentCategory, error) {
	return r.pg.GetCategoryByID(ctx, categoryID)
}
