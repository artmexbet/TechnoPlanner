package storeredis

import (
	"context"
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
	pipe.Set(ctx, fmt.Sprintf("session:%s", sessionID), sessionData, tokenTTL)
	pipe.Set(ctx, fmt.Sprintf("refresh:%s", refreshToken), sessionID, tokenTTL)

	pipe.SAdd(ctx, fmt.Sprintf("user_sessions:%s", userID), sessionID)
	pipe.Expire(ctx, fmt.Sprintf("user_sessions:%s", userID), tokenTTL)
	_, err := pipe.Exec(ctx)
	return err
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
