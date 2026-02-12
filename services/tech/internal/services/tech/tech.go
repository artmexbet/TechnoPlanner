package techservice

import (
	"context"

	"tech/internal/domain"
)

type Repository interface {
	AddEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error)
	DeleteEquipment(ctx context.Context, equipmentID int) error
	UpdateEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error)
	GetEquipmentByID(ctx context.Context, equipmentID int) (*domain.Equipment, error)
	GetAllEquipment(ctx context.Context) ([]domain.Equipment, error)
	GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error) {
	return s.repository.AddEquipment(ctx, equipment)
}

func (s *Service) DeleteEquipment(ctx context.Context, equipmentID int) error {
	return s.repository.DeleteEquipment(ctx, equipmentID)
}

func (s *Service) UpdateEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error) {
	return s.repository.UpdateEquipment(ctx, equipment)
}

func (s *Service) GetEquipmentByID(ctx context.Context, equipmentID int) (*domain.Equipment, error) {
	return s.repository.GetEquipmentByID(ctx, equipmentID)
}

func (s *Service) GetAllEquipment(ctx context.Context) ([]domain.Equipment, error) {
	return s.repository.GetAllEquipment(ctx)
}

func (s *Service) GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error) {
	return s.repository.GetEquipmentByCategory(ctx, categoryID)
}
