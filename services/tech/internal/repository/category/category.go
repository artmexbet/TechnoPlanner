package categoryrepo

import (
	"context"
	"tech/internal/domain"

	"github.com/google/uuid"
)

type IPostgres interface {
	GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error)
	AddCategory(ctx context.Context, categoryName string) (*domain.TechnicCategory, error)
	UpdateCategoryName(ctx context.Context, category domain.TechnicCategory) (*domain.TechnicCategory, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetAllCategories(ctx context.Context) ([]domain.TechnicCategory, error)
}

type Repository struct {
	pg IPostgres
}

func NewRepository(pg IPostgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error) {
	return r.pg.GetTechnicByCategory(ctx, categoryID)
}

func (r *Repository) AddCategory(ctx context.Context, categoryName string) (*domain.TechnicCategory, error) {
	return r.pg.AddCategory(ctx, categoryName)
}

func (r *Repository) UpdateCategoryName(ctx context.Context, category domain.TechnicCategory) (*domain.TechnicCategory, error) {
	return r.pg.UpdateCategoryName(ctx, category)
}

func (r *Repository) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return r.pg.DeleteCategory(ctx, categoryID)
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]domain.TechnicCategory, error) {
	return r.pg.GetAllCategories(ctx)
}
