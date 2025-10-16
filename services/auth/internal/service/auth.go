package service

import (
	"auth/internal/models"

	"context"
	"fmt"
)

type iTokenGenerator interface {
	GenerateTokenPair(userid string) (models.TokenPair, error)
	GenerateSession(u models.User, deviceID, userAgent, ip string) *models.Session
}

type iRepository interface {
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	StoreToken(ctx context.Context, session *models.Session, refreshToken string) error
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

func (a *Auth) Login(ctx context.Context, loginRequest models.LoginRequest) (models.TokenPair, error) {
	u, err := a.repository.GetUserByUsername(ctx, loginRequest.Username)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to get user by username: %w", err)
	}
	ok, err := u.CheckPassword(loginRequest.Password)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to check password: %w", err)
	}
	if !ok {
		return models.TokenPair{}, fmt.Errorf("invalid password")
	}

	// Generate token pair
	tokenPair, err := a.generator.GenerateTokenPair(u.ID.String())
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to generate token pair: %w", err)
	}

	session := a.generator.GenerateSession(u, loginRequest.DeviceID, loginRequest.UserAgent, loginRequest.IP)

	// Store session and tokens
	if err = a.repository.StoreToken(ctx, session, tokenPair.RefreshToken); err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to store token: %w", err)
	}
	return tokenPair, nil
}

func (a *Auth) Register(ctx context.Context, username, password string) error {
	return nil
}

func (a *Auth) ValidateToken(ctx context.Context, token string) (string, error) {
	return "username", nil
}
