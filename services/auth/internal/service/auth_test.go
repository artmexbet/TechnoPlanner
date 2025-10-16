package service

import (
	"context"
	"errors"
	"testing"

	"auth/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth_Login(t *testing.T) {
	mockPass := "pass"
	mockPassHash, _ := bcrypt.GenerateFromPassword([]byte(mockPass), bcrypt.DefaultCost)
	testUser := models.User{
		ID:           uuid.New(),
		Username:     "user1",
		PasswordHash: string(mockPassHash), // hash for "pass"
	}

	tests := []struct {
		name     string
		loginReq models.LoginRequest
		setup    func(repo *mockiRepository, tokenizer *mockiTokenizer)
		wantErr  bool
	}{
		{
			name: "успешный логин",
			loginReq: models.LoginRequest{
				Username:  "user1",
				Password:  "pass",
				DeviceID:  "dev1",
				UserAgent: "agent",
				IP:        "127.0.0.1",
			},
			setup: func(repo *mockiRepository, tokenizer *mockiTokenizer) {
				repo.EXPECT().
					GetUserByUsername(mock.Anything, "user1").
					Return(testUser, nil)

				tokenizer.EXPECT().
					GenerateTokenPair(testUser.ID.String()).
					Return(models.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil)

				tokenizer.EXPECT().
					GenerateSession(mock.Anything, "dev1", "agent", "127.0.0.1").
					Return(&models.Session{})

				repo.EXPECT().
					StoreToken(mock.Anything, mock.Anything, "refresh").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "неверный пароль",
			loginReq: models.LoginRequest{
				Username:  "user1",
				Password:  "wrong",
				DeviceID:  "dev1",
				UserAgent: "agent",
				IP:        "127.0.0.1",
			},
			setup: func(repo *mockiRepository, tokenizer *mockiTokenizer) {
				repo.EXPECT().
					GetUserByUsername(mock.Anything, "user1").
					Return(models.User{Username: "user1", PasswordHash: "$2a$10$7EqJtq98hPqEX7fNZaFWoO5r5h0g6vF2pQ5QyQp1pQ5QyQp1pQ5Qy"}, nil)
			},
			wantErr: true,
		},
		{
			name: "пользователь не найден",
			loginReq: models.LoginRequest{
				Username:  "nouser",
				Password:  "pass",
				DeviceID:  "dev1",
				UserAgent: "agent",
				IP:        "127.0.0.1",
			},
			setup: func(repo *mockiRepository, tokenizer *mockiTokenizer) {
				repo.EXPECT().
					GetUserByUsername(mock.Anything, "nouser").
					Return(models.User{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockiRepository(t)
			tokenizer := newMockiTokenizer(t)
			if tt.setup != nil {
				tt.setup(repo, tokenizer)
			}
			auth := NewAuth(tokenizer, repo)
			_, err := auth.Login(context.Background(), tt.loginReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("ожидалась ошибка: %v, получено: %v", tt.wantErr, err)
			}
		})
	}
}

func TestAuth_Register(t *testing.T) {
	tests := []struct {
		name    string
		regReq  models.RegisterRequest
		setup   func(repo *mockiRepository)
		wantErr bool
	}{
		{
			name: "успешная регистрация",
			regReq: models.RegisterRequest{
				Username: "newuser",
				Email:    "new@user.com",
				Password: "pass",
			},
			setup: func(repo *mockiRepository) {
				repo.EXPECT().
					CreateUser(mock.Anything, "newuser", "new@user.com", mock.Anything).
					Return(models.User{Username: "newuser", Email: "new@user.com"}, nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка создания пользователя",
			regReq: models.RegisterRequest{
				Username: "failuser",
				Email:    "fail@user.com",
				Password: "pass",
			},
			setup: func(repo *mockiRepository) {
				repo.EXPECT().
					CreateUser(mock.Anything, "failuser", "fail@user.com", mock.Anything).
					Return(models.User{}, errors.New("fail"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockiRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			auth := NewAuth(newMockiTokenizer(t), repo)
			_, err := auth.Register(context.Background(), tt.regReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("ожидалась ошибка: %v, получено: %v", tt.wantErr, err)
			}
		})
	}
}

func TestAuth_ValidateToken(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(tokenizer *mockiTokenizer)
		token     string
		wantState models.TokenState
		wantErr   bool
	}{
		{
			name: "валидный токен",
			setup: func(tokenizer *mockiTokenizer) {
				tokenizer.EXPECT().
					DecodeToken("valid").
					Return(&models.Claims{UserID: "id"}, nil)
			},
			token:     "valid",
			wantState: models.TokenStateValid,
			wantErr:   false,
		},
		{
			name: "просроченный токен",
			setup: func(tokenizer *mockiTokenizer) {
				tokenizer.EXPECT().
					DecodeToken("expired").
					Return(&models.Claims{UserID: "id"}, ErrTokenExpired)
			},
			token:     "expired",
			wantState: models.TokenStateExpired,
			wantErr:   false,
		},
		{
			name: "невалидный токен",
			setup: func(tokenizer *mockiTokenizer) {
				tokenizer.EXPECT().
					DecodeToken("invalid").
					Return(nil, errors.New("invalid"))
			},
			token:     "invalid",
			wantState: models.TokenStateInvalid,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := newMockiTokenizer(t)
			if tt.setup != nil {
				tt.setup(tokenizer)
			}
			auth := NewAuth(tokenizer, newMockiRepository(t))
			res, err := auth.ValidateToken(context.Background(), tt.token)
			if res.State != tt.wantState {
				t.Errorf("ожидалось состояние %v, получено %v", tt.wantState, res.State)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ожидалась ошибка: %v, получено: %v", tt.wantErr, err)
			}
		})
	}
}

func TestAuth_Refresh(t *testing.T) {
	tests := []struct {
		name       string
		refreshReq models.TokenRefreshRequest
		setup      func(repo *mockiRepository, tokenizer *mockiTokenizer)
		wantErr    bool
	}{
		{
			name: "успешный рефреш",
			refreshReq: models.TokenRefreshRequest{
				Pair:      models.TokenPair{RefreshToken: "refresh"},
				DeviceID:  "dev1",
				UserAgent: "agent",
				IP:        "127.0.0.1",
			},
			setup: func(repo *mockiRepository, tokenizer *mockiTokenizer) {
				repo.EXPECT().
					GetSessionByRefreshToken(mock.Anything, "refresh").
					Return(&models.Session{DeviceID: "dev1", UserAgent: "agent", IP: "127.0.0.1"}, nil)
				tokenizer.EXPECT().
					GenerateTokenPair(mock.Anything).
					Return(models.TokenPair{AccessToken: "access", RefreshToken: "refresh2"}, nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка получения сессии",
			refreshReq: models.TokenRefreshRequest{
				Pair:      models.TokenPair{RefreshToken: "bad"},
				DeviceID:  "dev1",
				UserAgent: "agent",
				IP:        "127.0.0.1",
			},
			setup: func(repo *mockiRepository, tokenizer *mockiTokenizer) {
				tokenizer.EXPECT().
					DecodeToken("bad").
					Return(&models.Claims{
						UserID: "user1",
					}, nil)
				repo.EXPECT().
					GetSessionByRefreshToken(mock.Anything, "bad").
					Return(nil, errors.New("not found"))
				repo.EXPECT().
					DeleteAllUserSessions(mock.Anything, "user1").
					Return(nil)
			},
			wantErr: true,
		},
		{
			name: "ошибка валидации сессии",
			refreshReq: models.TokenRefreshRequest{
				Pair:      models.TokenPair{RefreshToken: "refresh"},
				DeviceID:  "wrong",
				UserAgent: "agent",
				IP:        "127.0.0.1",
			},
			setup: func(repo *mockiRepository, tokenizer *mockiTokenizer) {
				tokenizer.EXPECT().
					DecodeToken("refresh").
					Return(&models.Claims{
						UserID: "user1",
					}, nil)
				repo.EXPECT().
					GetSessionByRefreshToken(mock.Anything, "refresh").
					Return(&models.Session{DeviceID: "dev1", UserAgent: "agent", IP: "127.0.0.1"}, nil)
				repo.EXPECT().
					DeleteAllUserSessions(mock.Anything, "user1").
					Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockiRepository(t)
			tokenizer := newMockiTokenizer(t)
			if tt.setup != nil {
				tt.setup(repo, tokenizer)
			}
			auth := NewAuth(tokenizer, repo)
			_, err := auth.Refresh(context.Background(), tt.refreshReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("ожидалась ошибка: %v, получено: %v", tt.wantErr, err)
			}
		})
	}
}
