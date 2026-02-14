package postgres

import (
	"context"

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
func (p *Postgres) ListResponsibles(ctx context.Context) ([]queries.Responsible, error) {
	return p.q.ListResponsibles(ctx)
}
