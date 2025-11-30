package storage

import (
	"context"

	"github.com/google/uuid"

	"gateway/internal/domain"
)

const (
	categoryCreatedSubject = "equipment.category.created"
	categoryUpdatedSubject = "equipment.category.updated"
	categoryDeletedSubject = "equipment.category.deleted"
)

type categoryRepository interface {
	CreateEquipmentCategory(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	UpdateEquipmentCategory(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	ListEquipmentCategories(ctx context.Context) ([]domain.EquipmentCategory, error)
	SoftDeleteEquipmentCategory(ctx context.Context, id int32, userID *uuid.UUID) error
}

type CategoryStorage struct {
	repo      categoryRepository
	publisher EventPublisher
}

func NewCategoryStorage(repo categoryRepository, publisher EventPublisher) *CategoryStorage {
	return &CategoryStorage{repo: repo, publisher: publisher}
}

func (s *CategoryStorage) Create(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	created, err := s.repo.CreateEquipmentCategory(ctx, cat)
	if err != nil {
		return domain.EquipmentCategory{}, err
	}
	if err := s.publish(ctx, categoryCreatedSubject, created); err != nil {
		return domain.EquipmentCategory{}, err
	}
	return created, nil
}

func (s *CategoryStorage) Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	updated, err := s.repo.UpdateEquipmentCategory(ctx, cat)
	if err != nil {
		return domain.EquipmentCategory{}, err
	}
	if err := s.publish(ctx, categoryUpdatedSubject, updated); err != nil {
		return domain.EquipmentCategory{}, err
	}
	return updated, nil
}

func (s *CategoryStorage) List(ctx context.Context) ([]domain.EquipmentCategory, error) {
	return s.repo.ListEquipmentCategories(ctx)
}

func (s *CategoryStorage) SoftDelete(ctx context.Context, id int32, userID *uuid.UUID) error {
	if err := s.repo.SoftDeleteEquipmentCategory(ctx, id, userID); err != nil {
		return err
	}
	payload := struct {
		ID        int32      `json:"id"`
		DeletedBy *uuid.UUID `json:"deleted_by,omitempty"`
	}{ID: id, DeletedBy: userID}
	return s.publish(ctx, categoryDeletedSubject, payload)
}

func (s *CategoryStorage) publish(ctx context.Context, subject string, payload interface{}) error {
	if s.publisher == nil || subject == "" {
		return nil
	}
	return s.publisher.Publish(ctx, subject, payload)
}
