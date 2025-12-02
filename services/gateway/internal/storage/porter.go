package storage

import (
	"context"

	"github.com/google/uuid"

	"gateway/internal/domain"
)

type porterRepository interface {
	ListPorters(ctx context.Context, roleID int32) ([]domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type PorterStorage struct {
	repo porterRepository
}

func NewPorterStorage(repo porterRepository) *PorterStorage {
	return &PorterStorage{repo: repo}
}

func (s *PorterStorage) List(ctx context.Context, roleID int32) ([]domain.User, error) {
	return s.repo.ListPorters(ctx, roleID)
}

func (s *PorterStorage) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
