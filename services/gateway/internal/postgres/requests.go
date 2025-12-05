package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gateway/internal/domain"
	"gateway/internal/postgres/queries"
)

func (d *DB) AssignRequestResponsible(ctx context.Context, requestID uuid.UUID, responsibleID *uuid.UUID, userID *uuid.UUID) (domain.Request, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	row, err := d.q.AssignResponsible(ctx, queries.AssignResponsibleParams{
		ID:                requestID,
		ResponsibleUserID: responsibleID,
		UpdatedBy:         userID,
	})
	if err != nil {
		return domain.Request{}, fmt.Errorf("AssignResponsible: %w", err)
	}

	return d.enrichRequest(ctx, row)
}

func (d *DB) GetRequestByID(ctx context.Context, id uuid.UUID) (domain.Request, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	row, err := d.q.GetRequestByID(ctx, id)
	if err != nil {
		return domain.Request{}, fmt.Errorf("GetRequestByID: %w", err)
	}

	return d.enrichRequest(ctx, row)
}

func (d *DB) ListRequests(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	if responsibleID == nil {
		return nil, fmt.Errorf("ListRequests: responsibleID is required")
	}

	rows, err := d.q.ListRequests(ctx, responsibleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("ListRequests: %w", err)
	}

	res := make([]domain.Request, 0, len(rows))
	for _, row := range rows {
		req, err := d.enrichRequest(ctx, row)
		if err != nil {
			return nil, err
		}
		res = append(res, req)
	}

	return res, nil
}

func (d *DB) enrichRequest(ctx context.Context, row queries.Request) (domain.Request, error) {
	req := mapRequest(row)
	equip, err := d.ListEquipmentForRequest(ctx, row.ID)
	if err != nil {
		return domain.Request{}, err
	}
	req.Equipment = equip
	return req, nil
}

func mapRequest(row queries.Request) domain.Request {
	return domain.Request{
		ID:                row.ID,
		TelegramUserInfo:  row.TelegramUserInfo,
		RequestText:       row.RequestText,
		Status:            domain.RequestStatus(row.Status),
		ScheduleTime:      row.ScheduleTime,
		EndTime:           row.EndTime,
		Address:           row.Address,
		ResponsibleUserID: row.ResponsibleUserID,
		Audit: domain.AuditFields{
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			DeletedAt: row.DeletedAt,
			CreatedBy: row.CreatedBy,
			UpdatedBy: row.UpdatedBy,
		},
	}
}
