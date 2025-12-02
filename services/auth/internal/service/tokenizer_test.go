package service

import (
	"context"
	"time"
)

func (s *ServiceTestSuite) TestTokenizer_Tokenize() {
	s.Run("valid token generation", func() {
		t := NewTokenizer(time.Hour, 24*time.Hour, "test")
		ctx := context.Background()
		userID := "user123"
		tokenPair, err := t.GenerateTokenPair(ctx, userID, RoleAdmin)
		s.Require().NoError(err)
		s.Require().NotEmpty(tokenPair.AccessToken)
		s.Require().NotEmpty(tokenPair.RefreshToken)
	})
	s.Run("invalid token", func() {
		t := NewTokenizer(time.Hour, 24*time.Hour, "test")
		ctx := context.Background()
		claims, err := t.DecodeToken(ctx, "invalid.token.here")
		s.Require().Error(err)
		s.Require().Nil(claims)
	})
	s.Run("valid token decoding", func() {
		t := NewTokenizer(time.Hour, 24*time.Hour, "test")
		ctx := context.Background()
		userID := "user123"
		tokenPair, err := t.GenerateTokenPair(ctx, userID, RoleAdmin)
		s.Require().NoError(err)

		claims, err := t.DecodeToken(ctx, tokenPair.AccessToken)
		s.Require().NoError(err)
		s.Require().Equal(userID, claims.UserID)
	})
}
