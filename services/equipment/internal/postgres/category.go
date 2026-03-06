package postgres

import (
	"context"
	"fmt"
	"time"

	"tech/internal/domain"
	"tech/internal/postgres/queries"
)

// AddCategory создаёт новую категорию оборудования.
func (p *Postgres) AddCategory(ctx context.Context, categoryName string, description string) (*domain.EquipmentCategory, error) {
	var descPtr *string
	if description != "" {
		descPtr = &description
	}
	id, err := p.q.AddCategory(ctx, queries.AddCategoryParams{
		Name:        categoryName,
		Description: descPtr,
	})
	if err != nil {
		return nil, fmt.Errorf("AddCategory query: %w", err)
	}
	return &domain.EquipmentCategory{
		ID:          int(id),
		Name:        categoryName,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// UpdateCategory обновляет имя и описание категории.
func (p *Postgres) UpdateCategory(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error) {
	var descPtr *string
	if category.Description != "" {
		descPtr = &category.Description
	}
	if err := p.q.UpdateCategory(ctx, queries.UpdateCategoryParams{
		Name:        category.Name,
		Description: descPtr,
		ID:          int32(category.ID),
	}); err != nil {
		return nil, fmt.Errorf("UpdateCategory query: %w", err)
	}
	return &category, nil
}

// DeleteCategory удаляет категорию по ID.
func (p *Postgres) DeleteCategory(ctx context.Context, categoryID int) error {
	if err := p.q.DeleteCategory(ctx, int32(categoryID)); err != nil {
		return fmt.Errorf("DeleteCategory query: %w", err)
	}
	return nil
}

// GetCategoryByID возвращает категорию по ID.
func (p *Postgres) GetCategoryByID(ctx context.Context, categoryID int) (*domain.EquipmentCategory, error) {
	rows, err := p.q.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllCategories query: %w", err)
	}
	for _, r := range rows {
		if int(r.ID) == categoryID {
			return mapCategory(r), nil
		}
	}
	return nil, fmt.Errorf("category %d not found", categoryID)
}

// GetAllCategories возвращает все категории.
func (p *Postgres) GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error) {
	rows, err := p.q.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllCategories query: %w", err)
	}
	result := make([]domain.EquipmentCategory, len(rows))
	for i, r := range rows {
		result[i] = *mapCategory(r)
	}
	return result, nil
}

func mapCategory(r queries.Category) *domain.EquipmentCategory {
	desc := ""
	if r.Description != nil {
		desc = *r.Description
	}
	return &domain.EquipmentCategory{
		ID:          int(r.ID),
		Name:        r.Name,
		Description: desc,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
