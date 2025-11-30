package service

import (
	"context"

	"github.com/google/uuid"

	"gateway/internal/domain"
)

type PorterStorage interface {
	List(ctx context.Context, roleID int32) ([]domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
}

const porterRoleID int32 = 2

var PorterRoleID int32 = porterRoleID

type PorterService struct {
	storage PorterStorage
}

func NewPorterService(storage PorterStorage) *PorterService {
	return &PorterService{storage: storage}
}

func (s *PorterService) List(ctx context.Context) ([]domain.User, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.storage.List(ctx, porterRoleID)
}

func (s *PorterService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.User{}, err
	}
	user, err := s.storage.Get(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	if user.RoleID != porterRoleID {
		return domain.User{}, ErrNotFound
	}
	return user, nil
}
