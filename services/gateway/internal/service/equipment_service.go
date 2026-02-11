package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

// EquipmentStorage defines storage contract for equipment CRUD operations.
type EquipmentStorage interface {
	Create(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	Update(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	Get(ctx context.Context, id int32) (domain.Equipment, error)
	List(ctx context.Context) ([]domain.Equipment, error)
	SoftDelete(ctx context.Context, id int32, userID *uuid.UUID) error
}

type EquipmentService struct {
	storage EquipmentStorage
}

func NewEquipmentService(storage EquipmentStorage) *EquipmentService {
	return &EquipmentService{storage: storage}
}

func (s *EquipmentService) Create(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.Equipment{}, err
	}
	if s.storage == nil {
		return domain.Equipment{}, errors.New("equipment storage not implemented")
	}
	eq.Audit.CreatedBy = userIDFromCtx(ctx)
	return s.storage.Create(ctx, eq)
}

func (s *EquipmentService) Update(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.Equipment{}, err
	}
	if s.storage == nil {
		return domain.Equipment{}, errors.New("equipment storage not implemented")
	}
	eq.Audit.UpdatedBy = userIDFromCtx(ctx)
	return s.storage.Update(ctx, eq)
}

func (s *EquipmentService) List(ctx context.Context) ([]domain.Equipment, error) {
	if s.storage == nil {
		return nil, errors.New("equipment storage not implemented")
	}
	return s.storage.List(ctx)
}

func (s *EquipmentService) Get(ctx context.Context, id int32) (domain.Equipment, error) {
	if s.storage == nil {
		return domain.Equipment{}, errors.New("equipment storage not implemented")
	}
	return s.storage.Get(ctx, id)
}

func (s *EquipmentService) Delete(ctx context.Context, id int32) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	if s.storage == nil {
		return errors.New("equipment storage not implemented")
	}
	return s.storage.SoftDelete(ctx, id, userIDFromCtx(ctx))
}
