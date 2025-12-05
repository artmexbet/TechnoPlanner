package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"gateway/internal/domain"
)

func (d *DB) ListEquipmentForRequest(ctx context.Context, requestID uuid.UUID) ([]domain.RequestEquipment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	rows, err := d.q.ListEquipmentForRequest(ctx, requestID)
	if err != nil && !errors.Is(err, driver.ErrBadConn) {
		return nil, fmt.Errorf("ListEquipmentForRequest: %w", err)
	}

	res := make([]domain.RequestEquipment, 0, len(rows))
	for _, row := range rows {
		res = append(res, domain.RequestEquipment{
			RequestID:   row.RequestID,
			EquipmentID: row.EquipmentID,
			Quantity:    row.Quantity,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return res, nil
}
