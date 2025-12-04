package models

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/artmexbet/TechnoPlanner/libs/proto"
)

type TokenState int32

func (t TokenState) ToProto() int32 {
	return int32(t)
}

const (
	TokenStateValid   TokenState = 0
	TokenStateExpired TokenState = 1
	TokenStateInvalid TokenState = 2
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (p *TokenPair) ToProto() *proto.TokenPair {
	return &proto.TokenPair{
		Token:        p.AccessToken,
		RefreshToken: p.RefreshToken,
	}
}

func TokenPairFromProto(in *proto.TokenPair) TokenPair {
	return TokenPair{
		AccessToken:  in.GetToken(),
		RefreshToken: in.GetRefreshToken(),
	}
}

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

type TokenValidateResult struct {
	State  TokenState
	UserID string
	Role   string
}

type TokenRefreshRequest struct {
	Pair      TokenPair
	DeviceID  string
	UserAgent string
	IP        string
}
