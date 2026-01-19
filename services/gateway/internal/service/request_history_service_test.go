package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

type RequestHistoryServiceSuite struct {
	suite.Suite
	storage *MockRequestHistoryStorage
	svc     *RequestHistoryService
}

func TestRequestHistoryServiceSuite(t *testing.T) {
	suite.Run(t, new(RequestHistoryServiceSuite))
}

func (s *RequestHistoryServiceSuite) SetupTest() {
	s.storage = NewMockRequestHistoryStorage(s.T())
	s.svc = NewRequestHistoryService(s.storage)
}

func (s *RequestHistoryServiceSuite) TestAdd_Admin() {
	adminID := uuid.New()
	ctx := withAdmin(adminID)
	entry := domain.RequestStatusHistory{RequestID: uuid.New()}
	s.storage.EXPECT().Add(mock.Anything, mock.MatchedBy(func(e domain.RequestStatusHistory) bool {
		return e.RequestID == entry.RequestID && e.ChangedBy != nil && e.ChangedBy.String() == adminID.String()
	})).Return(entry, nil)

	res, err := s.svc.Add(ctx, entry)

	s.Require().NoError(err)
	s.Equal(entry, res)
}

func (s *RequestHistoryServiceSuite) TestAdd_Forbidden() {
	ctx := withPorter(uuid.New())

	_, err := s.svc.Add(ctx, domain.RequestStatusHistory{})

	s.ErrorIs(err, ErrForbidden)
}

func (s *RequestHistoryServiceSuite) TestList() {
	ctx := context.Background()
	reqID := uuid.New()
	expected := []domain.RequestStatusHistory{{RequestID: reqID}}
	s.storage.EXPECT().List(mock.Anything, reqID).Return(expected, nil)

	items, err := s.svc.List(ctx, reqID)

	s.Require().NoError(err)
	s.Equal(expected, items)
}
