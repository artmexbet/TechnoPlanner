package repository

import (
	"auth/internal/models"
	"context"
	"encoding/json"
	"time"
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
}

type iPostgres interface {
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
	data, _ := json.Marshal(session)
	return r.r.StoreSession(ctx, session.SessionID, session.UserID, refreshToken, data, r.refreshTokenTTL)
}
