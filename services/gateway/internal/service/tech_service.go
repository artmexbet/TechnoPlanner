package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

// TechCall defines storage contract for equipment CRUD operations.
type TechCall interface {
	AddTechnic(ctx context.Context, technic domain.Technic) (domain.Technic, error)
	DeleteTechnic(ctx context.Context, techID uuid.UUID) error
	UpdateTechnic(ctx context.Context, technic domain.Technic) (domain.Technic, error)
	GetTechnicByID(ctx context.Context, techID uuid.UUID) (domain.Technic, error)
	GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error)
	AddCategory(ctx context.Context, categoryName string) (domain.TechnicCategory, error)
	UpdateCategoryName(ctx context.Context, category domain.TechnicCategory) (domain.TechnicCategory, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetAllCategories(ctx context.Context) ([]domain.TechnicCategory, error)
}

type TechService struct {
	caller TechCall
}

func NewTechService(caller TechCall) *TechService {
	return &TechService{caller: caller}
}

func (s *TechService) AddTechnic(ctx context.Context, eq domain.Technic) (domain.Technic, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.Technic{}, err
	}

	return s.caller.AddTechnic(ctx, eq)
}
