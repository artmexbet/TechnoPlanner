package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/postgres/queries"
)

func (p *Postgres) GetEquipmentByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.Equipment, error) {
	technics, err := p.q.GetEquipmentByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("GetEquipmentByRequestID: %w", err)
	}
	var result []domain.Equipment
	for _, t := range technics {
		result = append(result, *t.ToDomain())
	}
	return result, nil
}

func (p *Postgres) GetEquipmentByRequestIDs(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Equipment, error) {
	technicsRes := p.q.BatchGetEquipmentByRequestID(ctx, requestIDs)
	defer technicsRes.Close() //nolint:errcheck
	result := make(map[uuid.UUID][]domain.Equipment, len(requestIDs))
	errs := make([]error, len(requestIDs))
	technicsRes.Query(func(i int, technics []queries.BatchGetEquipmentByRequestIDRow, err error) {
		result[requestIDs[i]] = make([]domain.Equipment, 0, len(technics))
		for _, t := range technics {
			result[requestIDs[i]] = append(result[requestIDs[i]], *t.ToDomain())
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("GetEquipmentByRequestID for requestID %v: %w", requestIDs[i], err))
		}
	})
	return result, errors.Join(errs...)
}

func (p *Postgres) CreateEquipment(ctx context.Context, technics []domain.Equipment) error {
	dbTechnics := make([]queries.AddEquipmentParams, 0, len(technics))
	for _, t := range technics {
		dbTechnics = append(dbTechnics, queries.AddEquipmentParams{
			Name:        t.Name,
			Description: t.Description,
			Quantity:    int32(t.Quantity),
		})
	}

	br := p.q.AddEquipment(ctx, dbTechnics)
	defer br.Close() //nolint:errcheck
	errs := make([]error, len(technics))
	br.QueryRow(func(_ int, _ queries.Equipment, err error) {
		if err != nil {
			// Log the error but continue processing other technics
			errs = append(errs, err)
		}
	})
	return errors.Join(errs...)
}

// UpsertEquipment выполняет INSERT ON CONFLICT UPDATE для синхронизации с equipment сервисом.
func (p *Postgres) UpsertEquipment(ctx context.Context, eq domain.Equipment) error {
	return p.q.UpsertEquipment(ctx, queries.UpsertEquipmentParams{
		ID:               int32(eq.ID),
		Name:             eq.Name,
		Description:      eq.Description,
		Quantity:         int32(eq.Quantity),
		ReservedQuantity: int32(eq.ReservedQuantity),
	})
}

// DeleteEquipment удаляет оборудование из локальной копии.
func (p *Postgres) DeleteEquipment(ctx context.Context, id int) error {
	return p.q.DeleteEquipmentByID(ctx, int32(id))
}
