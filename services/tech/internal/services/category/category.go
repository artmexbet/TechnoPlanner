package categoryservice

import (
	"context"
	"tech/internal/domain"

	"github.com/google/uuid"
)

type IRepository interface {
	GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error)
	AddCategory(ctx context.Context, categoryName string) (*domain.TechnicCategory, error)
	UpdateCategoryName(ctx context.Context, category domain.TechnicCategory) (*domain.TechnicCategory, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetAllCategories(ctx context.Context) ([]domain.TechnicCategory, error)
}

type Service struct {
	repository IRepository
}

func New(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error) {
	return s.repository.GetTechnicByCategory(ctx, categoryID)
}

func (s *Service) AddCategory(ctx context.Context, categoryName string) (*domain.TechnicCategory, error) {
	return s.repository.AddCategory(ctx, categoryName)
}

func (s *Service) UpdateCategoryName(ctx context.Context, category domain.TechnicCategory) (*domain.TechnicCategory, error) {
	return s.repository.UpdateCategoryName(ctx, category)
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return s.repository.DeleteCategory(ctx, categoryID)
}

func (s *Service) GetAllCategories(ctx context.Context) ([]domain.TechnicCategory, error) {
	return s.repository.GetAllCategories(ctx)
}
