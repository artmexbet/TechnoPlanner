package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gateway/internal/domain"
)

type EquipmentServiceSuite struct {
	suite.Suite
	storage *MockEquipmentStorage
	svc     *EquipmentService
}

func TestEquipmentServiceSuite(t *testing.T) {
	suite.Run(t, new(EquipmentServiceSuite))
}

func (s *EquipmentServiceSuite) SetupTest() {
	s.storage = NewMockEquipmentStorage(s.T())
	s.svc = NewEquipmentService(s.storage)
}

func (s *EquipmentServiceSuite) TestCreate_Admin() {
	adminID := uuid.New()
	ctx := withAdmin(adminID)
	input := domain.Equipment{Name: "Crane"}
	expected := input
	s.storage.EXPECT().Create(mock.Anything, mock.MatchedBy(func(eq domain.Equipment) bool {
		return eq.Name == input.Name && eq.Audit.CreatedBy != nil && eq.Audit.CreatedBy.String() == adminID.String()
	})).Return(expected, nil)

	eq, err := s.svc.Create(ctx, input)

	s.Require().NoError(err)
	s.Equal(expected, eq)
}

func (s *EquipmentServiceSuite) TestCreate_Forbidden() {
	ctx := withPorter(uuid.New())

	_, err := s.svc.Create(ctx, domain.Equipment{})

	s.ErrorIs(err, ErrForbidden)
}

func (s *EquipmentServiceSuite) TestUpdate_Admin() {
	adminID := uuid.New()
	ctx := withAdmin(adminID)
	input := domain.Equipment{ID: 1, Name: "Lift"}
	expected := input
	s.storage.EXPECT().Update(mock.Anything, mock.MatchedBy(func(eq domain.Equipment) bool {
		return eq.ID == input.ID && eq.Audit.UpdatedBy != nil && eq.Audit.UpdatedBy.String() == adminID.String()
	})).Return(expected, nil)

	eq, err := s.svc.Update(ctx, input)

	s.Require().NoError(err)
	s.Equal(expected, eq)
}

func (s *EquipmentServiceSuite) TestDelete_Admin() {
	adminID := uuid.New()
	ctx := withAdmin(adminID)
	s.storage.EXPECT().SoftDelete(mock.Anything, int32(1), mock.Anything).Return(nil)

	err := s.svc.Delete(ctx, 1)

	s.Require().NoError(err)
}

func (s *EquipmentServiceSuite) TestDelete_Forbidden() {
	ctx := withPorter(uuid.New())

	err := s.svc.Delete(ctx, 1)

	s.ErrorIs(err, ErrForbidden)
}

func (s *EquipmentServiceSuite) TestList() {
	ctx := context.Background()
	expected := []domain.Equipment{{ID: 1}}
	s.storage.EXPECT().List(mock.Anything).Return(expected, nil)

	items, err := s.svc.List(ctx)

	s.Require().NoError(err)
	s.Equal(expected, items)
}

func (s *EquipmentServiceSuite) TestGet() {
	ctx := context.Background()
	s.storage.EXPECT().Get(mock.Anything, int32(1)).Return(domain.Equipment{ID: 1}, nil)

	eq, err := s.svc.Get(ctx, 1)

	s.Require().NoError(err)
	s.Equal(int32(1), eq.ID)
}
