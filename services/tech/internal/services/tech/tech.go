package techservice

import (
	"context"
	"tech/internal/domain"

	"github.com/google/uuid"
)

type IRepository interface {
	AddTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error)
	DeleteTechnic(ctx context.Context, techID uuid.UUID) error
	UpdateTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error)
	GetTechnicByID(ctx context.Context, techID uuid.UUID) (*domain.Technic, error)
}

type Service struct {
	repository IRepository
}

func New(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error) {
	return s.repository.AddTechnic(ctx, technic)
}
func (s *Service) DeleteTechnic(ctx context.Context, techID uuid.UUID) error {
	return s.repository.DeleteTechnic(ctx, techID)
}
func (s *Service) UpdateTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error) {
	return s.repository.UpdateTechnic(ctx, technic)
}
func (s *Service) GetTechnicByID(ctx context.Context, techID uuid.UUID) (*domain.Technic, error) {
	return s.repository.GetTechnicByID(ctx, techID)
}
