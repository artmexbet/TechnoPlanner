package postgres

import (
	"context"
	"fmt"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/postgres/queries"
)

// SaveResponsible сохраняет или обновляет ответственного пользователя
func (p *Postgres) SaveResponsible(ctx context.Context, id uuid.UUID, username string) error {
	_, err := p.q.SaveResponsible(ctx, queries.SaveResponsibleParams{
		ID:       id,
		Username: username,
	})
	return err
}

// GetResponsibleByID возвращает ответственного по ID
func (p *Postgres) GetResponsibleByID(ctx context.Context, id uuid.UUID) (queries.Responsible, error) {
	return p.q.GetResponsibleByID(ctx, id)
}

// GetResponsibleByUsername возвращает ответственного по username
func (p *Postgres) GetResponsibleByUsername(ctx context.Context, username string) (queries.Responsible, error) {
	return p.q.GetResponsibleByUsername(ctx, username)
}

// ListResponsibles возвращает список всех ответственных
func (p *Postgres) ListResponsibles(ctx context.Context) ([]domain.Responsible, error) {
	qResponsibles, err := p.q.ListResponsibles(ctx)
	if err != nil {
		return nil, fmt.Errorf("list responsibles: %w", err)
	}

	res := make([]domain.Responsible, len(qResponsibles))
	for i, r := range qResponsibles {
		res[i] = domain.Responsible{
			ID:       r.ID,
			Username: r.Username,
		}
	}
	return res, nil
}

// GetResponsible возвращает ответственного по ID
func (p *Postgres) GetResponsible(ctx context.Context, id uuid.UUID) (domain.Responsible, error) {
	r, err := p.q.GetResponsibleByID(ctx, id)
	if err != nil {
		return domain.Responsible{}, fmt.Errorf("get responsible: %w", err)
	}
	return domain.Responsible{
		ID:       r.ID,
		Username: r.Username,
	}, nil
}

// DeleteResponsible удаляет ответственного по ID
func (p *Postgres) DeleteResponsible(ctx context.Context, id uuid.UUID) error {
	return p.q.DeleteResponsible(ctx, id)
}
