package service

import (
	"context"
	"fmt"
	"proto"

	"gateway/internal/models"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceConfig struct {
	Address string `yaml:"address" env:"AUTH_SERVICE_ADDRESS"`
}

type GRPCWrapper struct {
	client   proto.AuthClient
	grpcConn *grpc.ClientConn
}

func NewGRPCWrapper(cfg AuthServiceConfig) (*GRPCWrapper, error) {
	conn, err := grpc.NewClient(cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(otelgrpc.WithPublicEndpoint())),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc client: %w", err)
	}
	conn.Connect()
	return &GRPCWrapper{
		client:   proto.NewAuthClient(conn),
		grpcConn: conn,
	}, nil
}

func (g *GRPCWrapper) Close() error {
	return g.grpcConn.Close()
}

func (g *GRPCWrapper) Login(ctx context.Context, req models.LoginRequest) (models.TokenPair, error) {
	resp, err := g.client.Login(ctx, &proto.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		DeviceId:  req.DeviceID,
		IpAddress: req.IP,
		UserAgent: req.UserAgent,
	})
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("grpc login error: %w", err)
	}
	return models.TokenPair{
		AccessToken:  resp.Token,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (g *GRPCWrapper) Register(ctx context.Context, username, password, email string) (string, error) {
	u, err := g.client.Register(ctx, &proto.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
	if err != nil {
		return "", fmt.Errorf("grpc register error: %w", err)
	}
	return u.UserId, nil
}

func (g *GRPCWrapper) ValidateToken(ctx context.Context, token string) (models.TokenValidationResponse, error) {
	resp, err := g.client.Validate(ctx, &proto.TokenRequest{
		Token: token,
	})
	if err != nil {
		return models.TokenValidationResponse{}, fmt.Errorf("grpc validate token error: %w", err)
	}
	return models.TokenValidationResponse{
		UserID: resp.UserId,
		State:  models.TokenState(resp.State.String()),
	}, nil
}

func (g *GRPCWrapper) Refresh(ctx context.Context, req models.TokenRefreshRequest) (models.TokenPair, error) {
	resp, err := g.client.Refresh(ctx, &proto.RefreshRequest{
		Pair: &proto.TokenPair{
			Token:        req.Pair.AccessToken,
			RefreshToken: req.Pair.RefreshToken,
		},
		DeviceId:  req.DeviceID,
		UserAgent: req.UserAgent,
		IpAddress: req.IP,
	})
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("grpc refresh token error: %w", err)
	}
	return models.TokenPair{
		AccessToken:  resp.Token,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (g *GRPCWrapper) Logout(ctx context.Context, token string) error {
	_, err := g.client.Logout(ctx, &proto.TokenRequest{
		Token: token,
	})
	if err != nil {
		return fmt.Errorf("grpc logout error: %w", err)
	}
	return nil
}

func (g *GRPCWrapper) LogoutAll(ctx context.Context, token string) error {
	_, err := g.client.LogoutAll(ctx, &proto.TokenRequest{
		Token: token,
	})
	if err != nil {
		return fmt.Errorf("grpc logout all error: %w", err)
	}
	return nil
}
