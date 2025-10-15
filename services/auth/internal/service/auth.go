package service

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Login(username, password string) (string, error) {
	return "token", nil
}

func (a *Auth) Register(username, password string) error {
	return nil
}

func (a *Auth) ValidateToken(token string) (string, error) {
	return "username", nil
}
