package postgres

import (
	"context"
	"log/slog"
	"tech/internal/domain"
	"tech/internal/postgres/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *Postgres) GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error) {
	dbTech, err := p.q.GetTechnicByCategory(ctx, pgtype.UUID{Bytes: categoryID})
	if err != nil {
		slog.Error("Tech: DB: GetTechnicByCategory", "error", err)
		return nil, err
	}
	return queries.DBTechSliceToDomain(dbTech), nil
}
func (p *Postgres) AddCategory(ctx context.Context, categoryName string) (*domain.TechnicCategory, error) {
	categoryID, err := p.q.AddCategory(ctx, categoryName)
	if err != nil {
		slog.Error("Tech: DB: GetTechnicByCategory", "error", err)
		return nil, err
	}
	category := domain.TechnicCategory{
		Name: categoryName,
		ID:   categoryID,
	}
	return &category, nil
}
func (p *Postgres) UpdateCategoryName(ctx context.Context, category domain.TechnicCategory) (*domain.TechnicCategory, error) {
	err := p.q.UpdateCategoryName(ctx, queries.UpdateCategoryNameParams{Name: category.Name, ID: category.ID})
	if err != nil {
		slog.Error("Tech: DB: GetTechnicByCategory", "error", err)
		return nil, err
	}
	return &category, nil
}
func (p *Postgres) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	err := p.q.DeleteCategory(ctx, categoryID)
	if err != nil {
		slog.Error("Tech: DB: GetTechnicByCategory", "error", err)
		return err
	}
	return nil
}
func (p *Postgres) GetAllCategories(ctx context.Context) ([]domain.TechnicCategory, error) {
	dbCategories, err := p.q.GetAllCategories(ctx)
	if err != nil {
		slog.Error("Tech: DB: GetTechnicByCategory", "error", err)
		return nil, err
	}
	domainCategories := queries.DBCategorySliceToDomain(dbCategories)
	return domainCategories, nil
}
