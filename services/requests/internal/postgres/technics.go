package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"requests/internal/postgres/queries"

	"requests/internal/domain"

	"github.com/google/uuid"
)

func (p *Postgres) GetTechnicsByRequestID(ctx context.Context, requestID uuid.UUID) ([]domain.Technic, error) {
	technics, err := p.q.GetTechnicsByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("GetTechnicsByRequestID: %w", err)
	}
	var result []domain.Technic
	for _, t := range technics {
		result = append(result, *t.ToDomain())
	}
	return result, nil
}

func (p *Postgres) GetTechnicByRequestIDs(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Technic, error) {
	technicsRes := p.q.BatchGetTechnicsByRequestID(ctx, requestIDs)
	defer technicsRes.Close()
	result := make(map[uuid.UUID][]domain.Technic, len(requestIDs))
	technicsRes.Query(func(i int, technics []queries.Technic, err error) {
		result[requestIDs[i]] = make([]domain.Technic, 0, len(technics))
		for _, t := range technics {
			result[requestIDs[i]] = append(result[requestIDs[i]], *t.ToDomain())
		}
	})
	return result, nil
}

func (p *Postgres) CreateTechnics(ctx context.Context, technics []domain.Technic) error {
	dbTechnics := make([]queries.AddTechnicParams, 0, len(technics))
	for _, t := range technics {
		dbTechnics = append(dbTechnics, queries.AddTechnicParams{
			Name:        t.Name,
			Description: t.Description,
			Quantity:    int32(t.Quantity),
		})
	}

	br := p.q.AddTechnic(ctx, dbTechnics)
	defer br.Close()
	br.QueryRow(func(i int, _ queries.Technic, err error) {
		if err != nil {
			// Log the error but continue processing other technics
			slog.Error("error adding technic", "err", err)
		}
	})
	return nil
}
