package service

type userStorage interface {
	// Define methods for user storage operations
}

type UserService struct {
	storage userStorage
}

func NewUserService(storage userStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}
