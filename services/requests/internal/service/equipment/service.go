package equipment

import (
	"context"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type Repository interface {
	CreateEquipment(ctx context.Context, technics []domain.Equipment) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Add(ctx context.Context, technics []domain.Equipment) error {
	return s.repository.CreateEquipment(ctx, technics)
}
