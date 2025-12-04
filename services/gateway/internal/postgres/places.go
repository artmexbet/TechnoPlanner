package postgres

import (
	"context"
	"fmt"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/postgres/queries"
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
		Description: place.Description,
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
