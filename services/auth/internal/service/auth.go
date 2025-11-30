package service

import (
	"context"
	"errors"
	"fmt"

	"auth/internal/models"
)

type iTokenizer interface {
	GenerateTokenPair(userid string, role string) (models.TokenPair, error)
	GenerateSession(u models.User, deviceID, userAgent, ip string) *models.Session
	DecodeToken(tokenStr string) (*models.Claims, error)
}

type iRepository interface {
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	StoreToken(ctx context.Context, session *models.Session, refreshToken string) error
	CreateUser(ctx context.Context, username, email, passwordHash string) (models.User, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID, userID string) error
	DeleteAllUserSessions(ctx context.Context, userID string) error
}

type Auth struct {
	tokenizer  iTokenizer
	repository iRepository
}

func NewAuth(tokenizer iTokenizer, repository iRepository) *Auth {
	return &Auth{
		tokenizer:  tokenizer,
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
	tokenPair, err := a.tokenizer.GenerateTokenPair(u.ID.String(), roleNameFromID(u.RoleID))
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to generate token pair: %w", err)
	}

	session := a.tokenizer.GenerateSession(u, loginRequest.DeviceID, loginRequest.UserAgent, loginRequest.IP)

	// Store session and tokens
	if err = a.repository.StoreToken(ctx, session, tokenPair.RefreshToken); err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to store token: %w", err)
	}
	return tokenPair, nil
}

func (a *Auth) Register(ctx context.Context, req models.RegisterRequest) (models.User, error) {
	h, err := req.HashPassword()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	u, err := a.repository.CreateUser(ctx, req.Username, req.Email, string(h))
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return u, nil
}

func (a *Auth) ValidateToken(_ context.Context, token string) (models.TokenValidateResult, error) {
	claims, err := a.tokenizer.DecodeToken(token)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return models.TokenValidateResult{State: models.TokenStateInvalid}, fmt.Errorf("failed to decode token: %w", err)
	} else if err != nil && errors.Is(err, ErrTokenExpired) {
		return models.TokenValidateResult{
			State:  models.TokenStateExpired,
			UserID: claims.UserID,
			Role:   claims.Role,
		}, nil
	}
	return models.TokenValidateResult{State: models.TokenStateValid, UserID: claims.UserID, Role: claims.Role}, nil
}

func (a *Auth) Refresh(ctx context.Context, req models.TokenRefreshRequest) (models.TokenPair, error) {
	session, err := a.repository.GetSessionByRefreshToken(ctx, req.Pair.RefreshToken)
	if err != nil {
		_ = a.LogoutAll(ctx, req.Pair.RefreshToken) // invalidate all sessions if refresh token is invalid
		return models.TokenPair{}, fmt.Errorf("failed to get session by refresh token: %w", err)
	}
	if err = session.Validate(req.DeviceID, req.UserAgent, req.IP); err != nil {
		_ = a.LogoutAll(ctx, req.Pair.RefreshToken) // invalidate all sessions if refresh token is invalid
		return models.TokenPair{}, fmt.Errorf("session validation failed: %w", err)
	}

	return a.tokenizer.GenerateTokenPair(session.UserID, session.Role)
}

func (a *Auth) Logout(ctx context.Context, token string) error {
	claims, err := a.tokenizer.DecodeToken(token)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return fmt.Errorf("failed to decode token: %w", err)
	}
	return a.repository.DeleteSession(ctx, claims.SessionID, claims.UserID)
}

func (a *Auth) LogoutAll(ctx context.Context, token string) error {
	claims, err := a.tokenizer.DecodeToken(token)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return fmt.Errorf("failed to decode token: %w", err)
	}
	userID := claims.UserID
	return a.repository.DeleteAllUserSessions(ctx, userID)
}

const (
	RoleAdmin  = "admin"
	RolePorter = "porter"
)

func roleNameFromID(id int32) string {
	switch id {
	case 1:
		return RoleAdmin
	case 2:
		return RolePorter
	default:
		return RolePorter
	}
}
