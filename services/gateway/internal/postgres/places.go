package postgres

import (
	"context"
	"fmt"

	"gateway/internal/domain"
	"gateway/internal/postgres/queries"

	"github.com/jackc/pgx/v5/pgtype"
)

func (d *DB) AddPlace(ctx context.Context, place domain.Place) (domain.Place, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cfg.Timeout)
	defer cancel()

	tx, err := d.p.Begin(ctx)
	if err != nil {
		return domain.Place{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	res, err := d.q.WithTx(tx).AddPlace(ctx, queries.AddPlaceParams{
		Name:        place.Name,
		Description: pgtype.Text{String: place.Description, Valid: place.Description != ""},
		Latitude:    place.Latitude,
		Longitude:   place.Longitude,
	})
	if err != nil {
		return domain.Place{}, fmt.Errorf("AddPlace: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.Place{}, fmt.Errorf("AddPlace: %w", err)
	}

	return res.ToDomain(), nil
}
