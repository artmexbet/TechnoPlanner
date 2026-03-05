package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

// PorterStorage интерфейс для работы с хранилищем портеров (Requests сервис)
type PorterStorage interface {
	List(ctx context.Context) ([]domain.Porter, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Porter, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, id uuid.UUID, username string) error
}

// AuthServiceConnector используется для создания портера через auth сервис
type AuthServiceConnector interface {
	RegisterPorter(ctx context.Context, username, password, email string) (string, error)
}

// UserStorage используется для обновления/удаления пользователя через auth
type UserStorage interface {
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, username, email string) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

const porterRoleID int32 = 2

var PorterRoleID int32 = porterRoleID

type PorterService struct {
	storage     PorterStorage
	userStorage UserStorage
	authSvc     AuthServiceConnector
}

func NewPorterService(storage PorterStorage, userStorage UserStorage, authSvc AuthServiceConnector) *PorterService {
	return &PorterService{
		storage:     storage,
		userStorage: userStorage,
		authSvc:     authSvc,
	}
}

func (s *PorterService) List(ctx context.Context) ([]domain.Porter, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.storage.List(ctx)
}

func (s *PorterService) Get(ctx context.Context, id uuid.UUID) (domain.Porter, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.Porter{}, err
	}
	return s.storage.Get(ctx, id)
}

// GetCurrentUser возвращает текущего пользователя по ID (для /me endpoint)
func (s *PorterService) GetCurrentUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.userStorage.Get(ctx, id)
}

func (s *PorterService) Create(ctx context.Context, username, email, password string) (string, error) {
	if err := requireAdmin(ctx); err != nil {
		return "", err
	}
	// Вызываем auth service для регистрации нового porter'а
	// При создании auth публикует UserCreated → Requests сохраняет как Porter автоматически
	userID, err := s.authSvc.RegisterPorter(ctx, username, password, email)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *PorterService) Update(ctx context.Context, id uuid.UUID, username, email string) (domain.User, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.User{}, err
	}
	return s.userStorage.Update(ctx, id, username, email)
}

func (s *PorterService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	// Удаляем из хранилища портеров (Requests)
	if err := s.storage.Delete(ctx, id); err != nil {
		return err
	}
	// Удаляем пользователя из auth сервиса
	return s.userStorage.Delete(ctx, id)
}
