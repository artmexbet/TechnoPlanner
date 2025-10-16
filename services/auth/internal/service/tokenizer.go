package service

import (
	"auth/internal/models"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired = fmt.Errorf("token expired")
)

type Tokenizer struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	jwtSecret  []byte
}

func NewTokenizer(accessTTL, refreshTTL time.Duration, jwtSecret string) *Tokenizer {
	return &Tokenizer{
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		jwtSecret:  []byte(jwtSecret),
	}
}

func (t *Tokenizer) GenerateTokenPair(userid string) (models.TokenPair, error) {
	sessionID := uuid.New()
	claims := models.Claims{
		UserID:    userid,
		SessionID: sessionID.String(),
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
		SessionID: sessionID.String(),
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

func (t *Tokenizer) GenerateSession(u models.User, deviceID, userAgent, ip string) *models.Session {
	sessionID := uuid.NewString()
	session := &models.Session{
		UserID:    u.ID.String(),
		SessionID: sessionID,
		DeviceID:  deviceID,
		UserAgent: userAgent,
		IP:        ip,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(t.refreshTTL),
	}
	return session
}

func (t *Tokenizer) DecodeToken(tokenStr string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return claims, ErrTokenExpired
	}
	return claims, nil
}
