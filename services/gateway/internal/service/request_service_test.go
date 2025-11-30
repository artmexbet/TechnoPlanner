package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gateway/internal/domain"
)

type RequestServiceSuite struct {
	suite.Suite
	storage *MockRequestStorage
	svc     *RequestService
}

func TestRequestServiceSuite(t *testing.T) {
	suite.Run(t, new(RequestServiceSuite))
}

func (s *RequestServiceSuite) SetupTest() {
	s.storage = NewMockRequestStorage(s.T())
	s.svc = NewRequestService(s.storage)
}

func (s *RequestServiceSuite) TestList_AdminFiltered() {
	ctx := withAdmin()
	responsible := uuid.New()
	expected := []domain.Request{{ID: uuid.New()}}
	s.storage.EXPECT().List(mock.Anything, &responsible).Return(expected, nil)

	reqs, err := s.svc.List(ctx, &responsible)

	s.Require().NoError(err)
	s.Equal(expected, reqs)
}

func (s *RequestServiceSuite) TestList_PorterForbidden() {
	porterID := uuid.New()
	ctx := withPorter(porterID)
	other := uuid.New()

	_, err := s.svc.List(ctx, &other)

	s.ErrorIs(err, ErrForbidden)
}

func (s *RequestServiceSuite) TestGet_Admin() {
	ctx := withAdmin()
	id := uuid.New()
	expected := domain.Request{ID: id}
	s.storage.EXPECT().Get(mock.Anything, id).Return(expected, nil)

	req, err := s.svc.Get(ctx, id)

	s.Require().NoError(err)
	s.Equal(expected, req)
}

func (s *RequestServiceSuite) TestGet_PorterOwnRequest() {
	porterID := uuid.New()
	ctx := withPorter(porterID)
	reqID := uuid.New()
	expected := domain.Request{ID: reqID, ResponsibleUserID: &porterID}
	s.storage.EXPECT().Get(mock.Anything, reqID).Return(expected, nil)

	req, err := s.svc.Get(ctx, reqID)

	s.Require().NoError(err)
	s.Equal(expected, req)
}

func (s *RequestServiceSuite) TestGet_PorterForeign() {
	porterID := uuid.New()
	ctx := withPorter(porterID)
	reqID := uuid.New()
	other := uuid.New()
	s.storage.EXPECT().Get(mock.Anything, reqID).Return(domain.Request{ID: reqID, ResponsibleUserID: &other}, nil)

	_, err := s.svc.Get(ctx, reqID)

	s.ErrorIs(err, ErrForbidden)
}

func (s *RequestServiceSuite) TestAssignResponsible_Admin() {
	ctx := withAdmin()
	reqID := uuid.New()
	respID := uuid.New()
	expected := domain.Request{ID: reqID, ResponsibleUserID: &respID}
	s.storage.EXPECT().AssignResponsible(mock.Anything, reqID, &respID, mock.Anything).Return(expected, nil)

	req, err := s.svc.AssignResponsible(ctx, reqID, respID)

	s.Require().NoError(err)
	s.Equal(expected, req)
}

func (s *RequestServiceSuite) TestAssignResponsible_Forbidden() {
	ctx := withPorter(uuid.New())

	_, err := s.svc.AssignResponsible(ctx, uuid.New(), uuid.New())

	s.ErrorIs(err, ErrForbidden)
}
