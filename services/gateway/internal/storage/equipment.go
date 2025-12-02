package storage

import (
	"context"

	"github.com/google/uuid"

	"gateway/internal/domain"
)

const (
	equipmentCreatedSubject = "equipment.item.created"
	equipmentUpdatedSubject = "equipment.item.updated"
	equipmentDeletedSubject = "equipment.item.deleted"
)

type equipmentRepository interface {
	CreateEquipment(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	UpdateEquipment(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	GetEquipment(ctx context.Context, id int32) (domain.Equipment, error)
	ListEquipment(ctx context.Context) ([]domain.Equipment, error)
	SoftDeleteEquipment(ctx context.Context, id int32, userID *uuid.UUID) error
}

type EquipmentStorage struct {
	repo      equipmentRepository
	publisher EventPublisher
}

func NewEquipmentStorage(repo equipmentRepository, publisher EventPublisher) *EquipmentStorage {
	return &EquipmentStorage{repo: repo, publisher: publisher}
}

func (s *EquipmentStorage) Create(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	created, err := s.repo.CreateEquipment(ctx, eq)
	if err != nil {
		return domain.Equipment{}, err
	}
	if err := s.publish(ctx, equipmentCreatedSubject, created); err != nil {
		return domain.Equipment{}, err
	}
	return created, nil
}

func (s *EquipmentStorage) Update(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	updated, err := s.repo.UpdateEquipment(ctx, eq)
	if err != nil {
		return domain.Equipment{}, err
	}
	if err := s.publish(ctx, equipmentUpdatedSubject, updated); err != nil {
		return domain.Equipment{}, err
	}
	return updated, nil
}

func (s *EquipmentStorage) Get(ctx context.Context, id int32) (domain.Equipment, error) {
	return s.repo.GetEquipment(ctx, id)
}

func (s *EquipmentStorage) List(ctx context.Context) ([]domain.Equipment, error) {
	return s.repo.ListEquipment(ctx)
}

func (s *EquipmentStorage) SoftDelete(ctx context.Context, id int32, userID *uuid.UUID) error {
	if err := s.repo.SoftDeleteEquipment(ctx, id, userID); err != nil {
		return err
	}
	payload := struct {
		ID        int32      `json:"id"`
		DeletedBy *uuid.UUID `json:"deleted_by,omitempty"`
	}{ID: id, DeletedBy: userID}
	return s.publish(ctx, equipmentDeletedSubject, payload)
}

func (s *EquipmentStorage) publish(ctx context.Context, subject string, payload interface{}) error {
	if s.publisher == nil || subject == "" {
		return nil
	}
	return s.publisher.Publish(ctx, subject, payload)
}
