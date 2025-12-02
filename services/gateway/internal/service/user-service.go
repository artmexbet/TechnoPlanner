package service

import "context"

// userStorage provides the operations needed by UserService.
type userStorage interface {
	// TODO: define methods when user flows are implemented.
}

type UserService struct {
	storage userStorage
}

func NewUserService(storage userStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) Placeholder(ctx context.Context) error {
	_ = ctx
	return nil
}
