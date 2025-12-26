package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type CategoryServiceSuite struct {
	suite.Suite
	storage *MockCategoryStorage
	svc     *CategoryService
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}

func (s *CategoryServiceSuite) SetupTest() {
	s.storage = NewMockCategoryStorage(s.T())
	s.svc = NewCategoryService(s.storage)
}

func (s *CategoryServiceSuite) TestCreate_Admin() {
	adminID := uuid.New()
	ctx := withAdmin(adminID)
	input := domain.EquipmentCategory{Name: "Heavy"}
	s.storage.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cat domain.EquipmentCategory) bool {
		return cat.Name == input.Name && cat.Audit.CreatedBy != nil && cat.Audit.CreatedBy.String() == adminID.String()
	})).Return(input, nil)

	cat, err := s.svc.Create(ctx, input)

	s.Require().NoError(err)
	s.Equal(input, cat)
}

func (s *CategoryServiceSuite) TestCreate_Forbidden() {
	ctx := withPorter(uuid.New())

	_, err := s.svc.Create(ctx, domain.EquipmentCategory{})

	s.ErrorIs(err, ErrForbidden)
}

func (s *CategoryServiceSuite) TestDelete_Admin() {
	ctx := withAdmin()
	s.storage.EXPECT().SoftDelete(mock.Anything, int32(1), mock.Anything).Return(nil)

	err := s.svc.Delete(ctx, 1)

	s.Require().NoError(err)
}

func (s *CategoryServiceSuite) TestList() {
	ctx := context.Background()
	expected := []domain.EquipmentCategory{{ID: 1}}
	s.storage.EXPECT().List(mock.Anything).Return(expected, nil)

	cats, err := s.svc.List(ctx)

	s.Require().NoError(err)
	s.Equal(expected, cats)
}

func (s *CategoryServiceSuite) TestUpdate_Admin() {
	adminID := uuid.New()
	ctx := withAdmin(adminID)
	input := domain.EquipmentCategory{ID: 1, Name: "Heavy"}
	s.storage.EXPECT().Update(mock.Anything, mock.MatchedBy(func(cat domain.EquipmentCategory) bool {
		return cat.ID == input.ID && cat.Audit.UpdatedBy != nil && cat.Audit.UpdatedBy.String() == adminID.String()
	})).Return(input, nil)

	cat, err := s.svc.Update(ctx, input)

	s.Require().NoError(err)
	s.Equal(input, cat)
}
