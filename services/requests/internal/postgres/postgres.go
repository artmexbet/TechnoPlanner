package postgres

import (
	"requests/internal/postgres/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate sqlc generate -f ./queries/sqlc.yaml

// Postgres is a struct that holds the Postgres connection pool and queries.
type Postgres struct {
	pool *pgxpool.Pool
	q    *queries.Queries
}

// New creates a new Postgres instance.
func New(pool *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: pool,
		q:    queries.New(pool),
	}
}
