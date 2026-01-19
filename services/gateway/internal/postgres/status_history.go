package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/postgres/queries"
)

func (d *DB) AddRequestStatusHistory(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	row, err := d.q.InsertRequestStatusHistory(ctx, queries.InsertRequestStatusHistoryParams{
		RequestID: entry.RequestID,
		Status:    queries.RequestStatus(entry.Status),
		Comment:   entry.Comment,
		ChangedBy: entry.ChangedBy,
	})
	if err != nil {
		return domain.RequestStatusHistory{}, fmt.Errorf("InsertRequestStatusHistory: %w", err)
	}

	return mapStatusHistory(row), nil
}

func (d *DB) ListRequestStatusHistory(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	rows, err := d.q.ListRequestStatusHistory(ctx, requestID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("ListRequestStatusHistory: %w", err)
	}

	res := make([]domain.RequestStatusHistory, 0, len(rows))
	for _, row := range rows {
		res = append(res, mapStatusHistory(row))
	}
	return res, nil
}

func mapStatusHistory(row queries.RequestStatusHistory) domain.RequestStatusHistory {
	return domain.RequestStatusHistory{
		ID:        row.ID,
		RequestID: row.RequestID,
		Status:    domain.RequestStatus(row.Status),
		Comment:   row.Comment,
		ChangedBy: row.ChangedBy,
		ChangedAt: row.ChangedAt,
	}
}
