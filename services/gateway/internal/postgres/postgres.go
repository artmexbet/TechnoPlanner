package postgres

import (
	"context"
	"fmt"

	"github.com/artmexbet/TechnoPlanner/libs/config"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/postgres/queries"
)

//go:generate sqlc generate -f ./queries/sqlc.yaml

type DB struct {
	cfg config.Postgres
	p   *pgxpool.Pool
	q   *queries.Queries
}

func New(ctx context.Context, cfg config.Postgres) (*DB, error) {
	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName),
	)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
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
