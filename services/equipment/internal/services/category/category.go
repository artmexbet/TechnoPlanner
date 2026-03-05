package categoryservice

import (
	"context"

	"tech/internal/domain"
)

type IRepository interface {
	AddCategory(ctx context.Context, categoryName string, description string) (*domain.EquipmentCategory, error)
	UpdateCategory(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error)
	DeleteCategory(ctx context.Context, categoryID int) error
	GetCategoryByID(ctx context.Context, categoryID int) (*domain.EquipmentCategory, error)
	GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error)
	GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error)
}

type Service struct {
	repository IRepository
}

func New(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddCategory(ctx context.Context, categoryName string, description string) (*domain.EquipmentCategory, error) {
	return s.repository.AddCategory(ctx, categoryName, description)
}

func (s *Service) UpdateCategory(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error) {
	return s.repository.UpdateCategory(ctx, category)
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID int) error {
	return s.repository.DeleteCategory(ctx, categoryID)
}

func (s *Service) GetCategoryByID(ctx context.Context, categoryID int) (*domain.EquipmentCategory, error) {
	return s.repository.GetCategoryByID(ctx, categoryID)
}

func (s *Service) GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error) {
	return s.repository.GetAllCategories(ctx)
}

func (s *Service) GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error) {
	return s.repository.GetEquipmentByCategory(ctx, categoryID)
}
