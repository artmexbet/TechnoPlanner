package equipment

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

var errRepo = errors.New("repo error")

type equipmentServiceTestSuite struct {
	suite.Suite
	repo    *MockRepository
	service *Service
}

func (s *equipmentServiceTestSuite) SetupTest() {
	s.repo = NewMockRepository(s.T())
	s.service = New(s.repo)
}

func (s *equipmentServiceTestSuite) TestAddSuccess() {
	technics := []domain.Equipment{{ID: 1}, {ID: 2}}

	s.repo.EXPECT().CreateEquipment(mock.Anything, technics).
		Return(nil)

	s.NoError(s.service.Add(context.Background(), technics))
}

func (s *equipmentServiceTestSuite) TestAdd_Error() {
	technics := []domain.Equipment{{ID: 1}}

	s.repo.EXPECT().CreateEquipment(mock.Anything, technics).
		Return(errRepo)

	s.Error(s.service.Add(context.Background(), technics))
}

func TestEquipmentServiceSuite(t *testing.T) {
	suite.Run(t, new(equipmentServiceTestSuite))
}
