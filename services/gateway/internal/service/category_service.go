package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type CategoryStorage interface {
	Create(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	List(ctx context.Context) ([]domain.EquipmentCategory, error)
	SoftDelete(ctx context.Context, id int32, userID *uuid.UUID) error
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
	cat.Audit.CreatedBy = userIDFromCtx(ctx)
	return s.storage.Create(ctx, cat)
}

func (s *CategoryService) Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.EquipmentCategory{}, err
	}
	cat.Audit.UpdatedBy = userIDFromCtx(ctx)
	return s.storage.Update(ctx, cat)
}

func (s *CategoryService) List(ctx context.Context) ([]domain.EquipmentCategory, error) {
	return s.storage.List(ctx)
}

func (s *CategoryService) Delete(ctx context.Context, id int32) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	return s.storage.SoftDelete(ctx, id, userIDFromCtx(ctx))
}
