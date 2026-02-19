package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type ResponsibleStorage interface {
	List(ctx context.Context) ([]domain.Responsible, error)
	Create(ctx context.Context, id uuid.UUID, username string) (domain.Responsible, error)
}

type ResponsibleService struct {
	storage ResponsibleStorage
}

func NewResponsibleService(storage ResponsibleStorage) *ResponsibleService {
	return &ResponsibleService{storage: storage}
}

func (s *ResponsibleService) List(ctx context.Context) ([]domain.Responsible, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.storage.List(ctx)
}

func (s *ResponsibleService) Create(ctx context.Context, id uuid.UUID, username string) (domain.Responsible, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.Responsible{}, err
	}
	return s.storage.Create(ctx, id, username)
}
