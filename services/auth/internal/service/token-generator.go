package service

import (
	"auth/internal/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenGenerator struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	jwtSecret  []byte
}

func NewTokenGenerator(accessTTL, refreshTTL time.Duration, jwtSecret string) *TokenGenerator {
	return &TokenGenerator{
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		jwtSecret:  []byte(jwtSecret),
	}
}

func (t *TokenGenerator) GenerateTokenPair(userid string) (models.TokenPair, error) {
	sessionId, _ := uuid.NewUUID()
	claims := models.Claims{
		UserID:    userid,
		SessionID: sessionId.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth_service",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.jwtSecret)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := models.Claims{
		UserID:    userid,
		SessionID: sessionId.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth_service",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(t.jwtSecret)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
