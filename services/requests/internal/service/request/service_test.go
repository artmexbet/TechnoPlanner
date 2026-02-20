package request

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

var errMock = errors.New("repo error")

type requestServiceTestSuite struct {
	suite.Suite
	repo         *MockRepository
	userProvider *mockUserProvider
	service      *Service
}

func (s *requestServiceTestSuite) SetupTest() {
	s.repo = NewMockRepository(s.T())
	s.userProvider = newMockUserProvider(s.T())
	s.service = New(s.repo, s.userProvider)
}

// mockUserProvider - простой мок для UserProvider
type mockUserProvider struct {
	mock.Mock
}

func newMockUserProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockUserProvider {
	m := &mockUserProvider{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *mockUserProvider) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (s *requestServiceTestSuite) TestAddSuccess() {
	issuerID := uuid.New()
	created := &domain.Request{ID: uuid.New(), Issuer: domain.User{ID: issuerID}}

	s.repo.EXPECT().SaveTelegramUser(mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.TelegramID == 123
	})).Return(domain.User{ID: issuerID, TelegramID: 123}, nil)

	s.repo.EXPECT().CreateRequest(mock.Anything, mock.MatchedBy(func(r domain.Request) bool {
		return r.Issuer.ID == issuerID
	})).Return(created, nil)

	got, err := s.service.Add(context.Background(), domain.Request{Issuer: domain.User{TelegramID: 123}})
	s.NoError(err)
	s.Equal(created, got)
}

func (s *requestServiceTestSuite) TestAdd_saveUserError() {
	s.repo.EXPECT().SaveTelegramUser(mock.Anything, mock.Anything).
		Return(domain.User{}, errMock)

	_, err := s.service.Add(context.Background(), domain.Request{})
	s.Error(err)
	s.True(errors.Is(err, errMock))
}

func (s *requestServiceTestSuite) TestListSuccess() {
	userID := uuid.New()
	expected := []domain.Request{{ID: uuid.New()}, {ID: uuid.New()}}

	s.repo.EXPECT().SaveTelegramUser(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, u domain.User) {
			u.ID = userID
		}).Return(domain.User{ID: userID}, nil)

	s.repo.EXPECT().GetRequestsByUserID(mock.Anything, userID, int32(10), int32(5)).
		Return(expected, nil)

	got, err := s.service.List(context.Background(), domain.User{TelegramID: 1}, 10, 5)
	s.NoError(err)
	s.Equal(expected, got)
}

func (s *requestServiceTestSuite) TestListRepoError() {
	userID := uuid.New()
	s.repo.EXPECT().SaveTelegramUser(mock.Anything, mock.Anything).
		Return(domain.User{ID: userID}, nil)

	s.repo.EXPECT().GetRequestsByUserID(mock.Anything, userID, int32(1), int32(0)).
		Return(nil, errMock)

	_, err := s.service.List(context.Background(), domain.User{}, 1, 0)
	s.Error(err)
	s.True(errors.Is(err, errMock))
}

func (s *requestServiceTestSuite) TestGetSuccess() {
	expected := &domain.Request{ID: uuid.New()}
	s.repo.EXPECT().GetRequestByID(mock.Anything, expected.ID).
		Return(expected, nil)

	got, err := s.service.Get(context.Background(), expected.ID)
	s.NoError(err)
	s.Equal(expected, got)
}

func (s *requestServiceTestSuite) TestCancelSuccess() {
	requestID := uuid.New()
	s.repo.EXPECT().UpdateRequestStatus(mock.Anything, requestID, domain.StatusCanceled).
		Return(nil)

	s.NoError(s.service.Cancel(context.Background(), requestID))
}

func (s *requestServiceTestSuite) TestUpdateStatusSuccess() {
	requestID := uuid.New()
	s.repo.EXPECT().UpdateRequestStatus(mock.Anything, requestID, domain.StatusInProgress).
		Return(nil)

	s.NoError(s.service.UpdateStatus(context.Background(), requestID, domain.StatusInProgress))
}

func TestRequestServiceSuite(t *testing.T) {
	suite.Run(t, new(requestServiceTestSuite))
}
