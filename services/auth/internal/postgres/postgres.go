package postgres

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth/internal/postgres/queries"
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
	q    *queries.Queries
}

func New(ctx context.Context, cfg Config) (*Postgres, error) {
	pgCfg, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Pass,
		cfg.Host, cfg.Port,
		cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("could not parse postgres config: %w", err)
	}

	pgCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	q := queries.New(pool)

	if err = otelpgx.RecordStats(pool); err != nil {
		return nil, fmt.Errorf("could not record pgx stats: %w", err)
	}
	return &Postgres{pool: pool, q: q}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}
