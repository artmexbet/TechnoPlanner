package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gateway/internal/domain"
	"gateway/internal/postgres/queries"
)

func (d *DB) CreateEquipment(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	tx, err := d.ensureTx(ctx)
	if err != nil {
		return domain.Equipment{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	res, err := d.q.WithTx(tx).CreateEquipment(ctx, queries.CreateEquipmentParams{
		Name:        eq.Name,
		Description: eq.Description,
		Quantity:    eq.Quantity,
		CreatedBy:   eq.Audit.CreatedBy,
	})
	if err != nil {
		return domain.Equipment{}, fmt.Errorf("CreateEquipment: %w", err)
	}

	if err := d.syncEquipmentCategories(ctx, tx, res.ID, eq.Categories); err != nil {
		return domain.Equipment{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.Equipment{}, fmt.Errorf("CreateEquipment: %w", err)
	}

	return d.getEquipmentWithCategories(ctx, res.ID)
}

func (d *DB) UpdateEquipment(ctx context.Context, eq domain.Equipment) (domain.Equipment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	tx, err := d.ensureTx(ctx)
	if err != nil {
		return domain.Equipment{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	res, err := d.q.WithTx(tx).UpdateEquipment(ctx, queries.UpdateEquipmentParams{
		ID:          eq.ID,
		Name:        eq.Name,
		Description: eq.Description,
		Quantity:    eq.Quantity,
		UpdatedBy:   eq.Audit.UpdatedBy,
	})
	if err != nil {
		return domain.Equipment{}, fmt.Errorf("UpdateEquipment: %w", err)
	}

	if err := d.syncEquipmentCategories(ctx, tx, res.ID, eq.Categories); err != nil {
		return domain.Equipment{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.Equipment{}, fmt.Errorf("UpdateEquipment: %w", err)
	}

	return d.getEquipmentWithCategories(ctx, res.ID)
}

func (d *DB) syncEquipmentCategories(ctx context.Context, tx pgx.Tx, equipmentID int32, cats []domain.EquipmentCategory) error {
	if err := d.q.WithTx(tx).ClearEquipmentCategories(ctx, equipmentID); err != nil {
		return fmt.Errorf("ClearEquipmentCategories: %w", err)
	}
	if len(cats) == 0 {
		return nil
	}
	ids := make([]int32, 0, len(cats))
	for _, c := range cats {
		ids = append(ids, c.ID)
	}
	if err := d.q.WithTx(tx).UpsertEquipmentCategories(ctx, queries.UpsertEquipmentCategoriesParams{
		EquipmentID: equipmentID,
		Column2:     ids,
	}); err != nil {
		return fmt.Errorf("UpsertEquipmentCategories: %w", err)
	}
	return nil
}

func (d *DB) GetEquipment(ctx context.Context, id int32) (domain.Equipment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	return d.getEquipmentWithCategories(ctx, id)
}

func (d *DB) ListEquipment(ctx context.Context) ([]domain.Equipment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	rows, err := d.q.ListEquipment(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListEquipment: %w", err)
	}

	result := make([]domain.Equipment, 0, len(rows))
	for _, row := range rows {
		item, err := d.getEquipmentWithCategories(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (d *DB) SoftDeleteEquipment(ctx context.Context, id int32, userID *uuid.UUID) error {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	return d.q.SoftDeleteEquipment(ctx, queries.SoftDeleteEquipmentParams{
		ID:        id,
		UpdatedBy: userID,
	})
}

func (d *DB) getEquipmentWithCategories(ctx context.Context, id int32) (domain.Equipment, error) {
	item, err := d.q.GetEquipmentByID(ctx, id)
	if err != nil {
		return domain.Equipment{}, err
	}
	cats, err := d.q.ListCategoriesForEquipment(ctx, id)
	if err != nil {
		return domain.Equipment{}, err
	}
	return mapEquipment(item, cats), nil
}

func mapEquipment(eq queries.Equipment, cats []queries.EquipmentCategory) domain.Equipment {
	return domain.Equipment{
		ID:          eq.ID,
		Name:        eq.Name,
		Description: eq.Description,
		Quantity:    eq.Quantity,
		Categories:  mapEquipmentCategories(cats),
		Audit: domain.AuditFields{
			CreatedAt: eq.CreatedAt,
			UpdatedAt: eq.UpdatedAt,
			DeletedAt: eq.DeletedAt,
			CreatedBy: eq.CreatedBy,
			UpdatedBy: eq.UpdatedBy,
		},
	}
}

func mapEquipmentCategories(rows []queries.EquipmentCategory) []domain.EquipmentCategory {
	if len(rows) == 0 {
		return nil
	}
	res := make([]domain.EquipmentCategory, 0, len(rows))
	for _, row := range rows {
		res = append(res, domain.EquipmentCategory{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Audit: domain.AuditFields{
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
				DeletedAt: row.DeletedAt,
				CreatedBy: row.CreatedBy,
				UpdatedBy: row.UpdatedBy,
			},
		})
	}
	return res
}
