package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gateway/internal/domain"
)

var testUserID = uuid.New()

type PorterServiceSuite struct {
	suite.Suite
	storage *MockPorterStorage
	svc     *PorterService
}

func TestPorterServiceSuite(t *testing.T) {
	suite.Run(t, new(PorterServiceSuite))
}

func (s *PorterServiceSuite) SetupTest() {
	s.storage = NewMockPorterStorage(s.T())
	s.svc = NewPorterService(s.storage)
}

func (s *PorterServiceSuite) TestList() {
	ctx := withAdmin()
	expected := []domain.User{{ID: testUserID, RoleID: porterRoleID}}
	s.storage.EXPECT().List(mock.Anything, porterRoleID).Return(expected, nil)

	users, err := s.svc.List(ctx)

	s.Require().NoError(err)
	s.Equal(expected, users)
	s.storage.AssertExpectations(s.T())
}

func (s *PorterServiceSuite) TestList_Forbidden() {
	ctx := withPorter(uuid.New())

	_, err := s.svc.List(ctx)

	s.Require().Error(err)
	s.ErrorIs(err, ErrForbidden)
}

func (s *PorterServiceSuite) TestGet() {
	ctx := withAdmin()
	targetID := uuid.New()
	expected := domain.User{ID: targetID, RoleID: porterRoleID}
	s.storage.EXPECT().Get(mock.Anything, targetID).Return(expected, nil)

	user, err := s.svc.Get(ctx, targetID)

	s.Require().NoError(err)
	s.Equal(expected, user)
	s.storage.AssertExpectations(s.T())
}

func (s *PorterServiceSuite) TestGet_NotPorter() {
	ctx := withAdmin()
	targetID := uuid.New()
	notPorter := domain.User{ID: targetID, RoleID: 1}
	s.storage.EXPECT().Get(mock.Anything, targetID).Return(notPorter, nil)

	_, err := s.svc.Get(ctx, targetID)

	s.Require().Error(err)
	s.ErrorIs(err, ErrNotFound)
}

func withAdmin(id ...uuid.UUID) context.Context {
	userID := uuid.New()
	if len(id) > 0 {
		userID = id[0]
	}
	return WithUserContext(context.Background(), userID.String(), string(RoleAdmin))
}

func withPorter(id uuid.UUID) context.Context {
	return WithUserContext(context.Background(), id.String(), string(RolePorter))
}
