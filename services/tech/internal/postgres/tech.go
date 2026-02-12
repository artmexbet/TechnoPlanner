package postgres

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"tech/internal/domain"
	"tech/internal/postgres/queries"
)

func (p *Postgres) AddTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error) {
	req, err := queries.ConvertTechToQueryAdd(technic)
	if err != nil {
		return nil, err
	}
	id, err := p.q.AddTechnic(ctx, *req)
	if err != nil {
		slog.Error("Tech: DB: AddTechnic", "error", err)
		return nil, err
	}
	technic.ID = id
	return &technic, nil
}
func (p *Postgres) DeleteTechnic(ctx context.Context, techID uuid.UUID) error {
	err := p.q.DeleteTechnic(ctx, techID)
	if err != nil {
		slog.Error("Tech: DB: DeleteTechnic", "error", err)
		return err
	}
	return nil
}
func (p *Postgres) UpdateTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error) {
	req, err := queries.ConvertTechToQueryUpdate(technic)
	if err != nil {
		return nil, err
	}
	err = p.q.UpdateTechnic(ctx, *req)
	if err != nil {
		slog.Error("Tech: DB: UpdateTechnic", "error", err)
		return nil, err
	}
	return &technic, nil
}
func (p *Postgres) GetTechnicByID(ctx context.Context, techID uuid.UUID) (*domain.Technic, error) {
	dbTech, err := p.q.GetTechnicByID(ctx, techID)
	if err != nil {
		slog.Error("Tech: DB: GetTechnicByID", "error", err)
		return nil, err
	}
	tech, err := dbTech.ConvertToDomain()
	if err != nil {
		return nil, err
	}
	return tech, nil
}
