package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/artmexbet/TechnoPlanner/services/auth/internal/models"
)

const (
	serviceName = "auth_service"
	tracerName  = "auth_service_tracer"
)

var (
	ErrTokenExpired = fmt.Errorf("token expired")
)

type Tokenizer struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	jwtSecret  []byte
	tracer     trace.Tracer
}

func NewTokenizer(accessTTL, refreshTTL time.Duration, jwtSecret string) *Tokenizer {
	return &Tokenizer{
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		jwtSecret:  []byte(jwtSecret),
		tracer:     otel.GetTracerProvider().Tracer(tracerName),
	}
}

func (t *Tokenizer) GenerateTokenPair(ctx context.Context, userid string, role string) (models.TokenPair, error) {
	_, span := t.tracer.Start(ctx, "GenerateTokenPair") //nolint:ineffassign
	defer span.End()

	sessionID := uuid.New()
	claims := models.Claims{
		UserID:    userid,
		SessionID: sessionID.String(),
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.accessTTL).UTC()),
			// TODO: подумать о том, чтобы вынести time.Now()
			// 	в репозиторий, чтобы можно было мокать время в тестах
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:   serviceName,
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.jwtSecret)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := models.Claims{
		UserID:    userid,
		SessionID: sessionID.String(),
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.refreshTTL).UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
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
		Role:      roleNameFromID(u.RoleID),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().Add(t.refreshTTL).UTC(),
	}
	return session
}

func (t *Tokenizer) DecodeToken(ctx context.Context, tokenStr string) (*models.Claims, error) {
	_, span := t.tracer.Start(ctx, "DecodeToken") //nolint:ineffassign
	defer span.End()

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
	if claims.ExpiresAt.Before(time.Now().UTC()) {
		return claims, ErrTokenExpired
	}
	return claims, nil
}
