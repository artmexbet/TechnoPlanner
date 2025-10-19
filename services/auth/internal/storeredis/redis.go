package storeredis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string `yaml:"addr" env:"ADDR"`
	User     string `yaml:"user" env:"USER"`
	Password string `yaml:"password" env:"PASSWORD"`
	DB       int    `yaml:"db" env:"DB"`
	PoolSize int    `yaml:"poolSize" env:"POOL_SIZE"`
}

type Redis struct {
	client *redis.Client
}

func New(cfg Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to redis: %w", err)
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, fmt.Errorf("could not instrument redis client: %w", err)
	}
	if err := redisotel.InstrumentMetrics(client); err != nil {
		return nil, fmt.Errorf("could not instrument redis client: %w", err)
	}

	return &Redis{client: client}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) StoreSession(
	ctx context.Context,
	sessionID, userID, refreshToken string,
	sessionData []byte,
	tokenTTL time.Duration) error {
	pipe := r.client.Pipeline()
	pipe.Set(ctx, sessionKey(sessionID), sessionData, tokenTTL)
	pipe.Set(ctx, fmt.Sprintf("refresh:%s", refreshToken), sessionID, tokenTTL)

	pipe.SAdd(ctx, fmt.Sprintf("user_sessions:%s", userID), sessionID)
	pipe.Expire(ctx, fmt.Sprintf("user_sessions:%s", userID), tokenTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func (r *Redis) GetSession(ctx context.Context, refreshToken string) ([]byte, error) {
	sessionID, err := r.client.Get(ctx, fmt.Sprintf("refresh:%s", refreshToken)).Result()
	if err != nil {
		return nil, fmt.Errorf("could not get session: %w", err)
	}

	data, err := r.client.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("could not get session data: %w", err)
	}
	return data, nil
}

func (r *Redis) DeleteSession(ctx context.Context, sessionID, userID string) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("session:%s", sessionID))
	pipe.SRem(ctx, fmt.Sprintf("user_sessions:%s", userID), sessionID)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *Redis) GetUserSessions(ctx context.Context, userID string) ([][]byte, error) {
	sessionIDs, err := r.client.SMembers(ctx, fmt.Sprintf("user_sessions:%s", userID)).Result()
	if err != nil {
		return nil, fmt.Errorf("could not get user sessions: %w", err)
	}

	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(sessionIDs))
	for i, sessionID := range sessionIDs {
		cmds[i] = pipe.Get(ctx, fmt.Sprintf("session:%s", sessionID))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("could not execute pipeline: %w", err)
	}
	sessions := make([][]byte, 0, len(sessionIDs))
	for _, cmd := range cmds {
		data, err := cmd.Bytes()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("could not get session data: %w", err)
		}
		if err == nil {
			sessions = append(sessions, data)
		}
	}
	return sessions, nil
}

func (r *Redis) DeleteAllUserSessions(ctx context.Context, userID string) error {
	sessionIDs, err := r.client.SMembers(ctx, fmt.Sprintf("user_sessions:%s", userID)).Result()
	if err != nil {
		return fmt.Errorf("could not get user sessions: %w", err)
	}

	pipe := r.client.Pipeline()
	for _, sessionID := range sessionIDs {
		pipe.Del(ctx, fmt.Sprintf("session:%s", sessionID))
	}
	pipe.Del(ctx, fmt.Sprintf("user_sessions:%s", userID))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("could not delete user sessions: %w", err)
	}
	return nil
}
