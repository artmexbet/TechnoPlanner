package service

import "auth/internal/models"

type iTokenGenerator interface {
	GenerateTokenPair(userid string) (models.TokenPair, error)
}

type iRepository interface {
	// Define repository methods here
}

type Auth struct {
	generator  iTokenGenerator
	repository iRepository
}

func NewAuth(generator iTokenGenerator, repository iRepository) *Auth {
	return &Auth{
		generator:  generator,
		repository: repository,
	}
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
