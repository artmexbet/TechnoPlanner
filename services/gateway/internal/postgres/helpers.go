package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (d *DB) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d.cfg.Timeout)
}

func (d *DB) ensureTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := d.p.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
