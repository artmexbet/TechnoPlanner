package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"tech/internal/domain"
	"tech/internal/postgres/queries"
)

// AddEquipment добавляет новое оборудование в БД.
func (p *Postgres) AddEquipment(ctx context.Context, eq domain.Equipment) (*domain.Equipment, error) {
	chars, err := json.Marshal(eq.AdditionalCharacteristics)
	if err != nil {
		return nil, fmt.Errorf("marshal additional_characteristics: %w", err)
	}
	catID := int32(eq.CategoryID)

	id, err := p.q.AddEquipment(ctx, queries.AddEquipmentParams{
		CategoryID:                &catID,
		Name:                      eq.Name,
		Description:               strPtr(eq.Description),
		AdditionalCharacteristics: chars,
		Quantity:                  int32(eq.Quantity),
	})
	if err != nil {
		return nil, fmt.Errorf("AddEquipment query: %w", err)
	}
	eq.ID = int(id)
	return &eq, nil
}

// DeleteEquipment удаляет оборудование по ID.
func (p *Postgres) DeleteEquipment(ctx context.Context, equipmentID int) error {
	if err := p.q.DeleteEquipment(ctx, int32(equipmentID)); err != nil {
		return fmt.Errorf("DeleteEquipment query: %w", err)
	}
	return nil
}

// UpdateEquipment обновляет поля оборудования.
func (p *Postgres) UpdateEquipment(ctx context.Context, eq domain.Equipment) (*domain.Equipment, error) {
	chars, err := json.Marshal(eq.AdditionalCharacteristics)
	if err != nil {
		return nil, fmt.Errorf("marshal additional_characteristics: %w", err)
	}
	catID := int32(eq.CategoryID)

	if err := p.q.UpdateEquipment(ctx, queries.UpdateEquipmentParams{
		CategoryID:                &catID,
		Name:                      eq.Name,
		Description:               strPtr(eq.Description),
		AdditionalCharacteristics: chars,
		Quantity:                  int32(eq.Quantity),
		ID:                        int32(eq.ID),
	}); err != nil {
		return nil, fmt.Errorf("UpdateEquipment query: %w", err)
	}
	return &eq, nil
}

// GetEquipmentByID возвращает оборудование по ID.
func (p *Postgres) GetEquipmentByID(ctx context.Context, equipmentID int) (*domain.Equipment, error) {
	row, err := p.q.GetEquipmentByID(ctx, int32(equipmentID))
	if err != nil {
		return nil, fmt.Errorf("GetEquipmentByID query: %w", err)
	}
	eq := mapEquipment(row)
	return &eq, nil
}

// GetAllEquipment возвращает весь каталог оборудования.
func (p *Postgres) GetAllEquipment(ctx context.Context) ([]domain.Equipment, error) {
	rows, err := p.q.GetAllEquipment(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllEquipment query: %w", err)
	}
	result := make([]domain.Equipment, len(rows))
	for i, r := range rows {
		result[i] = mapEquipment(r)
	}
	return result, nil
}

// GetEquipmentByCategory возвращает оборудование по категории.
func (p *Postgres) GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error) {
	catID := int32(categoryID)
	rows, err := p.q.GetTechnicByCategory(ctx, &catID)
	if err != nil {
		return nil, fmt.Errorf("GetEquipmentByCategory query: %w", err)
	}
	result := make([]domain.Equipment, len(rows))
	for i, r := range rows {
		result[i] = mapEquipment(r)
	}
	return result, nil
}

// ReserveEquipment резервирует указанные единицы оборудования.
// Возвращает ошибку если хотя бы для одной позиции недостаточно свободных единиц.
func (p *Postgres) ReserveEquipment(ctx context.Context, items []domain.ReserveItem) error {
	for _, item := range items {

		if err := p.q.ReserveEquipment(
			ctx,
			queries.ReserveEquipmentParams{
				ID:               int32(item.EquipmentID),
				ReservedQuantity: int32(item.Quantity),
			},
		); err != nil {
			return fmt.Errorf("ReserveEquipment id=%d qty=%d: %w", item.EquipmentID, item.Quantity, err)
		}
	}
	return nil
}

// ReleaseEquipment освобождает зарезервированные единицы оборудования.
func (p *Postgres) ReleaseEquipment(ctx context.Context, items []domain.ReserveItem) error {
	for _, item := range items {
		if err := p.q.ReleaseEquipment(ctx,
			queries.ReleaseEquipmentParams{
				ID:               int32(item.EquipmentID),
				ReservedQuantity: int32(item.Quantity),
			},
		); err != nil {
			return fmt.Errorf("ReleaseEquipment id=%d qty=%d: %w", item.EquipmentID, item.Quantity, err)
		}
	}
	return nil
}

// CheckAvailability проверяет, доступно ли нужное количество для каждой позиции.
func (p *Postgres) CheckAvailability(ctx context.Context, items []domain.ReserveItem) (bool, []int, error) {
	var unavailable []int
	for _, item := range items {
		available, err := p.q.GetAvailableQuantity(ctx, int32(item.EquipmentID))
		if err != nil {
			return false, nil, fmt.Errorf("GetAvailableQuantity id=%d: %w", item.EquipmentID, err)
		}
		if int(available) < item.Quantity {
			unavailable = append(unavailable, item.EquipmentID)
		}
	}
	return len(unavailable) == 0, unavailable, nil
}

func mapEquipment(r queries.Equipment) domain.Equipment {
	var chars map[string]string
	if len(r.AdditionalCharacteristics) > 0 {
		_ = json.Unmarshal(r.AdditionalCharacteristics, &chars)
	}
	catID := 0
	if r.CategoryID != nil {
		catID = int(*r.CategoryID)
	}
	desc := ""
	if r.Description != nil {
		desc = *r.Description
	}
	return domain.Equipment{
		ID:                        int(r.ID),
		CategoryID:                catID,
		Name:                      r.Name,
		Description:               desc,
		AdditionalCharacteristics: chars,
		Quantity:                  int(r.Quantity),
		ReservedQuantity:          int(r.ReservedQuantity),
		CreatedAt:                 r.CreatedAt,
		UpdatedAt:                 r.UpdatedAt,
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
