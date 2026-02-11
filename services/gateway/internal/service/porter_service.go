package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type PorterStorage interface {
	List(ctx context.Context, roleID int32) ([]domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type AuthServiceConnector interface {
	RegisterPorter(ctx context.Context, username, password, email string) (string, error)
}

const porterRoleID int32 = 2

var PorterRoleID int32 = porterRoleID

type PorterService struct {
	storage PorterStorage
	authSvc AuthServiceConnector
}

func NewPorterService(storage PorterStorage, authSvc AuthServiceConnector) *PorterService {
	return &PorterService{
		storage: storage,
		authSvc: authSvc,
	}
}

func (s *PorterService) List(ctx context.Context) ([]domain.User, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if s.storage == nil {
		return nil, errors.New("porter storage not implemented")
	}
	return s.storage.List(ctx, porterRoleID)
}

func (s *PorterService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.User{}, err
	}
	if s.storage == nil {
		return domain.User{}, errors.New("porter storage not implemented")
	}
	user, err := s.storage.Get(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	if user.RoleID != porterRoleID {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

// GetCurrentUser returns user by ID without role check (for /me endpoint)
func (s *PorterService) GetCurrentUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if s.storage == nil {
		return domain.User{}, errors.New("porter storage not implemented")
	}
	return s.storage.Get(ctx, id)
}

func (s *PorterService) Create(ctx context.Context, username, email, password string) (string, error) {
	if err := requireAdmin(ctx); err != nil {
		return "", err
	}
	// Вызываем auth service для регистрации нового porter'а
	userID, err := s.authSvc.RegisterPorter(ctx, username, password, email)
	if err != nil {
		return "", err
	}
	return userID, nil
}
