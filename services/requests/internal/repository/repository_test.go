package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type postgresStub struct {
	listRequestsFn      func(ctx context.Context, offset, limit int32) ([]domain.Request, error)
	getEquipmentByIDsFn func(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Equipment, error)
}

func (p *postgresStub) CreateRequest(_ context.Context, _ domain.Request) (*domain.Request, error) {
	panic("unexpected call")
}
func (p *postgresStub) CreateEquipment(_ context.Context, _ []domain.Equipment) error {
	panic("unexpected call")
}
func (p *postgresStub) UpsertEquipment(_ context.Context, _ domain.Equipment) error {
	panic("unexpected call")
}
func (p *postgresStub) DeleteEquipment(_ context.Context, _ int) error { panic("unexpected call") }
func (p *postgresStub) GetRequestsByUserID(_ context.Context, _ uuid.UUID, _, _ int32) ([]domain.Request, error) {
	panic("unexpected call")
}
func (p *postgresStub) UpdateRequestStatus(_ context.Context, _ uuid.UUID, _ domain.StatusType) error {
	panic("unexpected call")
}
func (p *postgresStub) GetRequestByID(_ context.Context, _ uuid.UUID) (*domain.Request, error) {
	panic("unexpected call")
}
func (p *postgresStub) AssignEquipmentToRequest(_ context.Context, _ uuid.UUID, _ []domain.Equipment) []error {
	panic("unexpected call")
}
func (p *postgresStub) GetEquipmentByRequestID(_ context.Context, _ uuid.UUID) ([]domain.Equipment, error) {
	panic("unexpected call")
}
func (p *postgresStub) GetEquipmentByRequestIDs(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Equipment, error) {
	return p.getEquipmentByIDsFn(ctx, requestIDs)
}
func (p *postgresStub) GetUserByTelegramID(_ context.Context, _ int64) (domain.User, error) {
	panic("unexpected call")
}
func (p *postgresStub) SaveTelegramUser(_ context.Context, _ domain.User) (domain.User, error) {
	panic("unexpected call")
}
func (p *postgresStub) GetUserByID(_ context.Context, _ uuid.UUID) (domain.User, error) {
	panic("unexpected call")
}
func (p *postgresStub) GetRequestsByResponsibleID(_ context.Context, _ *uuid.UUID) ([]domain.Request, error) {
	panic("unexpected call")
}
func (p *postgresStub) ListRequests(ctx context.Context, offset, limit int32) ([]domain.Request, error) {
	return p.listRequestsFn(ctx, offset, limit)
}
func (p *postgresStub) ListResponsibles(_ context.Context) ([]domain.Porter, error) {
	panic("unexpected call")
}
func (p *postgresStub) GetResponsible(_ context.Context, _ uuid.UUID) (domain.Porter, error) {
	panic("unexpected call")
}
func (p *postgresStub) DeleteResponsible(_ context.Context, _ uuid.UUID) error {
	panic("unexpected call")
}
func (p *postgresStub) AssignResponsible(_ context.Context, _ uuid.UUID, _ *uuid.UUID) error {
	panic("unexpected call")
}
func (p *postgresStub) UpdateRequest(_ context.Context, _ uuid.UUID, _ domain.RequestUpdate) (*domain.Request, error) {
	panic("unexpected call")
}
func (p *postgresStub) SaveResponsible(_ context.Context, _ uuid.UUID, _ string) error {
	panic("unexpected call")
}
func (p *postgresStub) CreateRawRequest(_ context.Context, _ domain.RawRequest) (*domain.RawRequest, error) {
	panic("unexpected call")
}
func (p *postgresStub) GetRawRequests(_ context.Context, _ string, _, _ int32) ([]domain.RawRequest, error) {
	panic("unexpected call")
}
func (p *postgresStub) GetRawRequestByID(_ context.Context, _ uuid.UUID) (*domain.RawRequest, error) {
	panic("unexpected call")
}
func (p *postgresStub) MarkRawRequestProcessed(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.RawRequest, error) {
	panic("unexpected call")
}

func TestRepositoryListRequests_ForwardsLimitAndOffsetAndEnrichesEquipment(t *testing.T) {
	requestID := uuid.New()
	var gotOffset, gotLimit int32
	pg := &postgresStub{
		listRequestsFn: func(ctx context.Context, offset, limit int32) ([]domain.Request, error) {
			gotOffset = offset
			gotLimit = limit
			return []domain.Request{{ID: requestID}}, nil
		},
		getEquipmentByIDsFn: func(ctx context.Context, requestIDs []uuid.UUID) (map[uuid.UUID][]domain.Equipment, error) {
			if len(requestIDs) != 1 || requestIDs[0] != requestID {
				t.Fatalf("unexpected request ids: %#v", requestIDs)
			}
			return map[uuid.UUID][]domain.Equipment{
				requestID: {{ID: 7, Quantity: 3}},
			}, nil
		},
	}

	repo := &Repository{pg: pg}
	items, err := repo.ListRequests(context.Background(), 25, 10)
	if err != nil {
		t.Fatalf("ListRequests returned error: %v", err)
	}
	if gotLimit != 25 || gotOffset != 10 {
		t.Fatalf("expected limit=25 offset=10, got limit=%d offset=%d", gotLimit, gotOffset)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 request, got %d", len(items))
	}
	if len(items[0].Equipments) != 1 || items[0].Equipments[0].ID != 7 || items[0].Equipments[0].Quantity != 3 {
		t.Fatalf("expected enriched equipment on request, got %#v", items[0].Equipments)
	}
}
