package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"gateway/internal/postgres/queries"
)

//go:generate sqlc generate -f ./queries/sqlc.yaml

type Config struct {
	User string `yaml:"user" env:"USER"`
	Pass string `yaml:"password" env:"PASSWORD"`
	Host string `yaml:"host" env:"HOST"`
	Port string `yaml:"port" env:"PORT"`
	Db   string `yaml:"db" env:"DB"`

	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT"`
}

type DB struct {
	cfg Config
	p   *pgxpool.Pool
	q   *queries.Queries
}

func New(ctx context.Context, cfg Config) (*DB, error) {
	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Db),
	)
	if err != nil {
		return nil, err
	}
	return &DB{
		cfg: cfg,
		p:   pool,
		q:   queries.New(pool),
	}, nil
}

func (d *DB) Close() {
	d.p.Close()
}
