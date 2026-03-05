package postgres

import (
	"context"
	"fmt"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/postgres/queries"
)

// SaveResponsible сохраняет или обновляет портера
func (p *Postgres) SaveResponsible(ctx context.Context, id uuid.UUID, username string) error {
	_, err := p.q.SaveResponsible(ctx, queries.SaveResponsibleParams{
		ID:       id,
		Username: username,
	})
	return err
}

// GetResponsibleByID возвращает портера по ID
func (p *Postgres) GetResponsibleByID(ctx context.Context, id uuid.UUID) (queries.Porter, error) {
	return p.q.GetResponsibleByID(ctx, id)
}

// GetResponsibleByUsername возвращает портера по username
func (p *Postgres) GetResponsibleByUsername(ctx context.Context, username string) (queries.Porter, error) {
	return p.q.GetResponsibleByUsername(ctx, username)
}

// ListResponsibles возвращает список всех портеров
func (p *Postgres) ListResponsibles(ctx context.Context) ([]domain.Porter, error) {
	qPorters, err := p.q.ListResponsibles(ctx)
	if err != nil {
		return nil, fmt.Errorf("list porters: %w", err)
	}

	res := make([]domain.Porter, len(qPorters))
	for i, r := range qPorters {
		res[i] = domain.Porter{
			ID:       r.ID,
			Username: r.Username,
		}
	}
	return res, nil
}

// GetResponsible возвращает портера по ID
func (p *Postgres) GetResponsible(ctx context.Context, id uuid.UUID) (domain.Porter, error) {
	r, err := p.q.GetResponsibleByID(ctx, id)
	if err != nil {
		return domain.Porter{}, fmt.Errorf("get porter: %w", err)
	}
	return domain.Porter{
		ID:       r.ID,
		Username: r.Username,
	}, nil
}

// DeleteResponsible удаляет портера по ID
func (p *Postgres) DeleteResponsible(ctx context.Context, id uuid.UUID) error {
	return p.q.DeleteResponsible(ctx, id)
}
