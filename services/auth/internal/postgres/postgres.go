package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate sqlc generate -f ./queries/sqlc.yaml

type Config struct {
	Host   string `yaml:"host" env:"HOST"`
	Port   int    `yaml:"port" env:"PORT"`
	User   string `yaml:"user" env:"USER"`
	Pass   string `yaml:"pass" env:"PASS"`
	DBName string `yaml:"db_name" env:"DB_NAME"`

	SSLMode string `yaml:"sslmode" env:"SSLMODE"`
}

type Postgres struct {
	pool *pgxpool.Pool
}

func New(cfg Config) (*Postgres, error) {
	pool, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.User, cfg.Pass,
			cfg.Host, cfg.Port,
			cfg.DBName, cfg.SSLMode),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	return &Postgres{pool: pool}, nil
}
