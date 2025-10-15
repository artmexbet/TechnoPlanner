package service

type authStorage interface {
}

type Auth struct {
	storage authStorage
}

func NewAuthService(storage authStorage) *Auth {
	return &Auth{
		storage: storage,
	}
}

func (a *Auth) Login(username, password string) (string, error) {
	return "token", nil
}

func (a *Auth) Register(username, password, email string) error {
	return nil
}
