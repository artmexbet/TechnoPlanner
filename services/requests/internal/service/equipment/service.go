package equipment

import (
	"context"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type Repository interface {
	CreateEquipment(ctx context.Context, technics []domain.Equipment) error
	UpsertEquipment(ctx context.Context, eq domain.Equipment) error
	DeleteEquipment(ctx context.Context, id int) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Add(ctx context.Context, technics []domain.Equipment) error {
	return s.repository.CreateEquipment(ctx, technics)
}

// SyncCreate сохраняет новое оборудование из equipment сервиса в локальную копию.
func (s *Service) SyncCreate(ctx context.Context, eq domain.Equipment) error {
	return s.repository.UpsertEquipment(ctx, eq)
}

// SyncUpdate обновляет оборудование в локальной копии.
func (s *Service) SyncUpdate(ctx context.Context, eq domain.Equipment) error {
	return s.repository.UpsertEquipment(ctx, eq)
}

// SyncDelete удаляет оборудование из локальной копии.
func (s *Service) SyncDelete(ctx context.Context, id int) error {
	return s.repository.DeleteEquipment(ctx, id)
}
