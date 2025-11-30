package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"auth/internal/models"
)

type Config struct {
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"ACCESS_TOKEN_TTL"`
}

type iRedis interface {
	StoreSession(
		ctx context.Context,
		sessionID, userID, refreshToken string,
		sessionData []byte,
		tokenTTL time.Duration) error
	GetSession(ctx context.Context, refreshToken string) ([]byte, error)
	DeleteSession(ctx context.Context, sessionID, userID string) error
	DeleteAllUserSessions(ctx context.Context, userID string) error
	GetUserSessions(ctx context.Context, userID string) ([][]byte, error)
}

type iPostgres interface {
	FindUserByUsername(ctx context.Context, username string) (models.User, error)
	CreateUser(ctx context.Context, username, email, passwordHash string, roleID int32) (models.User, error)
}

type Repository struct {
	r iRedis
	p iPostgres

	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
}

func New(cfg Config, r iRedis, p iPostgres) (*Repository, error) {
	return &Repository{
		r:               r,
		p:               p,
		refreshTokenTTL: cfg.RefreshTokenTTL,
		accessTokenTTL:  cfg.AccessTokenTTL,
	}, nil
}

func (r *Repository) StoreToken(ctx context.Context, session *models.Session, refreshToken string) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}
	return r.r.StoreSession(ctx, session.SessionID, session.UserID, refreshToken, data, r.refreshTokenTTL)
}

func (r *Repository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	data, err := r.r.GetSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	var session models.Session
	if err = json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}
	return &session, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return r.p.FindUserByUsername(ctx, username)
}

func (r *Repository) CreateUser(ctx context.Context, username, email, passwordHash string) (models.User, error) {
	const defaultRoleID int32 = 1
	return r.p.CreateUser(ctx, username, email, passwordHash, defaultRoleID)
}

func (r *Repository) DeleteSession(ctx context.Context, sessionID, userID string) error {
	return r.r.DeleteSession(ctx, sessionID, userID)
}

func (r *Repository) DeleteAllUserSessions(ctx context.Context, userID string) error {
	return r.r.DeleteAllUserSessions(ctx, userID)
}

func (r *Repository) GetUserSessions(ctx context.Context, userID string) ([]*models.Session, error) {
	dataList, err := r.r.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user sessions: %w", err)
	}

	sessions := make([]*models.Session, 0, len(dataList))
	for _, data := range dataList {
		var session models.Session
		if err = json.Unmarshal(data, &session); err != nil {
			return nil, fmt.Errorf("unmarshal session: %w", err)
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}
