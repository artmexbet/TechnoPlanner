package service

import (
	"context"
	"errors"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type CategoryStorage interface {
	Create(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	List(ctx context.Context) ([]domain.EquipmentCategory, error)
	Delete(ctx context.Context, id int) error
}

type CategoryService struct {
	storage CategoryStorage
}

func NewCategoryService(storage CategoryStorage) *CategoryService {
	return &CategoryService{storage: storage}
}

func (s *CategoryService) Create(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.EquipmentCategory{}, err
	}
	if s.storage == nil {
		return domain.EquipmentCategory{}, errors.New("category storage not implemented")
	}
	return s.storage.Create(ctx, cat)
}

func (s *CategoryService) Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.EquipmentCategory{}, err
	}
	if s.storage == nil {
		return domain.EquipmentCategory{}, errors.New("category storage not implemented")
	}
	return s.storage.Update(ctx, cat)
}

func (s *CategoryService) List(ctx context.Context) ([]domain.EquipmentCategory, error) {
	if s.storage == nil {
		return nil, errors.New("category storage not implemented")
	}
	return s.storage.List(ctx)
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	if s.storage == nil {
		return errors.New("category storage not implemented")
	}
	return s.storage.Delete(ctx, id)
}
