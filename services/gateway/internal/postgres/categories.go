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

func (d *DB) CreateEquipmentCategory(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	res, err := d.q.CreateEquipmentCategory(ctx, queries.CreateEquipmentCategoryParams{
		Name:        cat.Name,
		Description: cat.Description,
		CreatedBy:   cat.Audit.CreatedBy,
	})
	if err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("CreateEquipmentCategory: %w", err)
	}

	return mapCategory(res), nil
}

func (d *DB) UpdateEquipmentCategory(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	res, err := d.q.UpdateEquipmentCategory(ctx, queries.UpdateEquipmentCategoryParams{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		UpdatedBy:   cat.Audit.UpdatedBy,
	})
	if err != nil {
		return domain.EquipmentCategory{}, fmt.Errorf("UpdateEquipmentCategory: %w", err)
	}

	return mapCategory(res), nil
}

func (d *DB) SoftDeleteEquipmentCategory(ctx context.Context, id int32, userID *uuid.UUID) error {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	return d.q.SoftDeleteEquipmentCategory(ctx, queries.SoftDeleteEquipmentCategoryParams{
		ID:        id,
		UpdatedBy: userID,
	})
}

func (d *DB) ListEquipmentCategories(ctx context.Context) ([]domain.EquipmentCategory, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	rows, err := d.q.ListEquipmentCategories(ctx)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("ListEquipmentCategories: %w", err)
	}

	res := make([]domain.EquipmentCategory, 0, len(rows))
	for _, row := range rows {
		res = append(res, mapCategory(row))
	}
	return res, nil
}

func mapCategory(cat queries.EquipmentCategory) domain.EquipmentCategory {
	return domain.EquipmentCategory{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		Audit: domain.AuditFields{
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
			DeletedAt: cat.DeletedAt,
			CreatedBy: cat.CreatedBy,
			UpdatedBy: cat.UpdatedBy,
		},
	}
}
