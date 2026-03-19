package router

import (
	"testing"
	"time"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

func TestToEquipmentResponse_IncludesQuantityAndCategory(t *testing.T) {
	now := time.Date(2026, time.March, 19, 12, 0, 0, 0, time.UTC)
	eq := domain.Equipment{
		ID:          11,
		CategoryID:  3,
		Name:        "Projector",
		Description: "4K",
		Quantity:    9,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	categoryMap := map[int]domain.EquipmentCategory{
		3: {
			ID:          3,
			Name:        "Screens",
			Description: "Display equipment",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	resp := toEquipmentResponse(eq, categoryMap)
	if resp.Quantity != 9 {
		t.Fatalf("expected quantity 9, got %d", resp.Quantity)
	}
	if len(resp.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(resp.Categories))
	}
	if resp.Categories[0].ID != 3 || resp.Categories[0].Name != "Screens" {
		t.Fatalf("unexpected categories payload: %#v", resp.Categories)
	}
}
