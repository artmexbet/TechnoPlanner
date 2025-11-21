package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"requests/internal/domain"
	"testing"
	"time"

	"requests/internal/wrapnats"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	nConntainer "github.com/testcontainers/testcontainers-go/modules/nats"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Mock для сервисов
type MockRequestService struct {
	mock.Mock
}

func (m *MockRequestService) Add(ctx context.Context, newRequest domain.Request) (*domain.Request, error) {
	args := m.Called(ctx, newRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Request), args.Error(1)
}

func (m *MockRequestService) UpdateStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error {
	args := m.Called(ctx, requestID, status)
	return args.Error(0)
}

func (m *MockRequestService) Get(ctx context.Context, requestID uuid.UUID) (*domain.Request, error) {
	args := m.Called(ctx, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Request), args.Error(1)
}

func (m *MockRequestService) List(ctx context.Context, user domain.User, limit, offset int32) ([]domain.Request, error) {
	args := m.Called(ctx, user, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Request), args.Error(1)
}

func (m *MockRequestService) Cancel(ctx context.Context, requestID uuid.UUID) error {
	args := m.Called(ctx, requestID)
	return args.Error(0)
}

type MockEquipmentService struct {
	mock.Mock
}

func (m *MockEquipmentService) Add(ctx context.Context, technics []domain.Equipment) error {
	args := m.Called(ctx, technics)
	return args.Error(0)
}

// Вспомогательные функции
func setupNATSContainer(t *testing.T) (*nats.Conn, func()) {
	ctx := context.Background()
	container, err := nConntainer.Run(
		ctx,
		"nats:latest",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Server is ready").
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "4222/tcp")
	require.NoError(t, err)

	connStr := fmt.Sprintf("nats://127.0.0.1:%s", port.Port())
	natsConn, err := nats.Connect(connStr)
	require.NoError(t, err)

	cleanup := func() {
		natsConn.Close()
		if err := testcontainers.TerminateContainer(container); err != nil {
			slog.Error("failed to terminate container", "err", err)
		}
	}

	return natsConn, cleanup
}

func setupNATSWrapper(t *testing.T, natsConn *nats.Conn, mockReqSvc *MockRequestService, mockEqSvc *MockEquipmentService) *wrapnats.NatsWrapper {
	cfg := &wrapnats.Config{
		RequestTimeout: 5 * time.Second,
	}

	wrapper, err := wrapnats.New(cfg, natsConn, mockReqSvc, mockEqSvc)
	require.NoError(t, err)

	wrapper.HandleMsgs()
	time.Sleep(200 * time.Millisecond) // Даем время подписчикам зарегистрироваться

	return wrapper
}

// Тесты Publisher
func TestPublishRequestCreated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	publisher := wrapnats.NewNatsPublisher(natsConn)
	received := make(chan domain.Request, 1)

	_, err := natsConn.Subscribe("events.request.created", func(msg *nats.Msg) {
		var req domain.Request
		_ = req.UnmarshalJSON(msg.Data)
		received <- req
	})
	require.NoError(t, err)

	testReq := domain.Request{
		ID:          uuid.New(),
		RequestText: stringPtr("Test request"),
		Status:      domain.StatusPending,
	}

	err = publisher.PublishRequestCreated(testReq)
	require.NoError(t, err)

	select {
	case req := <-received:
		assert.Equal(t, testReq.ID, req.ID)
		assert.Equal(t, testReq.RequestText, req.RequestText)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestPublishRequestCanceled(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	publisher := wrapnats.NewNatsPublisher(natsConn)
	received := make(chan domain.Request, 1)

	_, err := natsConn.Subscribe("events.request.canceled", func(msg *nats.Msg) {
		var req domain.Request
		_ = req.UnmarshalJSON(msg.Data)
		received <- req
	})
	require.NoError(t, err)

	testReq := domain.Request{
		ID:     uuid.New(),
		Status: domain.StatusCanceled,
	}

	err = publisher.PublishRequestCanceled(testReq)
	require.NoError(t, err)

	select {
	case req := <-received:
		assert.Equal(t, testReq.ID, req.ID)
		assert.Equal(t, domain.StatusCanceled, req.Status)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestPublishUserAdded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	publisher := wrapnats.NewNatsPublisher(natsConn)
	received := make(chan domain.User, 1)

	_, err := natsConn.Subscribe("events.user.added", func(msg *nats.Msg) {
		var user domain.User
		_ = user.UnmarshalJSON(msg.Data)
		received <- user
	})
	require.NoError(t, err)

	testUser := domain.User{
		ID:         uuid.New(),
		TelegramID: 123456,
		Username:   "testuser",
		FirstName:  "Test",
	}

	err = publisher.PublishUserAdded(testUser)
	require.NoError(t, err)

	select {
	case user := <-received:
		assert.Equal(t, testUser.TelegramID, user.TelegramID)
		assert.Equal(t, testUser.Username, user.Username)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

// Тесты NatsWrapper handlers
func TestHandleCreateRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	mockReqSvc := new(MockRequestService)
	mockEqSvc := new(MockEquipmentService)

	_ = setupNATSWrapper(t, natsConn, mockReqSvc, mockEqSvc)

	testID := uuid.New()
	expectedReq := &domain.Request{
		ID:           testID,
		RequestText:  stringPtr("Test request"),
		Status:       domain.StatusPending,
		ScheduleTime: "2025-12-01T10:00:00Z",
		Address:      "Test Address",
	}

	mockReqSvc.On("Add", mock.Anything, mock.Anything).Return(expectedReq, nil)

	reqData := map[string]interface{}{
		"text":          "Test request",
		"schedule_time": "2025-12-01T10:00:00Z",
		"telegram_id":   int64(123456),
		"equipments": []map[string]int{
			{"id": 1, "quantity": 2},
		},
		"address": "Test Address",
	}

	data, err := json.Marshal(reqData)
	require.NoError(t, err)

	msg, err := natsConn.Request("requests.create", data, 2*time.Second)
	require.NoError(t, err)

	var resp map[string]interface{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp["message"])
	assert.False(t, resp["isError"].(bool))
	mockReqSvc.AssertExpectations(t)
}

func TestHandleUpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	mockReqSvc := new(MockRequestService)
	mockEqSvc := new(MockEquipmentService)

	_ = setupNATSWrapper(t, natsConn, mockReqSvc, mockEqSvc)

	testID := uuid.New()
	mockReqSvc.On("UpdateStatus", mock.Anything, testID, domain.StatusCompleted).Return(nil)

	reqData := map[string]interface{}{
		"request_id": testID.String(),
		"status":     "completed",
	}

	data, err := json.Marshal(reqData)
	require.NoError(t, err)

	msg, err := natsConn.Request("requests.status.update", data, 2*time.Second)
	require.NoError(t, err)

	var resp map[string]interface{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp["message"])
	assert.False(t, resp["isError"].(bool))
	mockReqSvc.AssertExpectations(t)
}

func TestHandleGetRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	mockReqSvc := new(MockRequestService)
	mockEqSvc := new(MockEquipmentService)

	_ = setupNATSWrapper(t, natsConn, mockReqSvc, mockEqSvc)

	testID := uuid.New()
	expectedReq := &domain.Request{
		ID:           testID,
		RequestText:  stringPtr("Test request"),
		Status:       domain.StatusPending,
		ScheduleTime: "2025-12-01T10:00:00Z",
	}

	mockReqSvc.On("Get", mock.Anything, testID).Return(expectedReq, nil)

	reqData := map[string]interface{}{
		"request_id": testID.String(),
	}

	data, err := json.Marshal(reqData)
	require.NoError(t, err)

	msg, err := natsConn.Request("requests.get", data, 5*time.Second)
	require.NoError(t, err)

	var resp map[string]interface{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp["message"])
	assert.False(t, resp["isError"].(bool))
	mockReqSvc.AssertExpectations(t)
}

func TestHandleListRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	mockReqSvc := new(MockRequestService)
	mockEqSvc := new(MockEquipmentService)

	_ = setupNATSWrapper(t, natsConn, mockReqSvc, mockEqSvc)

	expectedReqs := []domain.Request{
		{
			ID:           uuid.New(),
			RequestText:  stringPtr("Request 1"),
			Status:       domain.StatusPending,
			ScheduleTime: "2025-12-01T10:00:00Z",
		},
		{
			ID:           uuid.New(),
			RequestText:  stringPtr("Request 2"),
			Status:       domain.StatusCompleted,
			ScheduleTime: "2025-12-02T10:00:00Z",
		},
	}

	telegramID := int64(123456)
	user := domain.User{TelegramID: telegramID}

	mockReqSvc.On("List", mock.Anything, user, int32(10), int32(0)).Return(expectedReqs, nil)

	reqData := map[string]interface{}{
		"telegram_id": telegramID,
		"limit":       10,
		"offset":      0,
	}

	data, err := json.Marshal(reqData)
	require.NoError(t, err)

	msg, err := natsConn.Request("requests.list", data, 2*time.Second)
	require.NoError(t, err)

	var resp map[string]interface{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp["message"])
	assert.False(t, resp["isError"].(bool))
	mockReqSvc.AssertExpectations(t)
}

func TestHandleCancelRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	mockReqSvc := new(MockRequestService)
	mockEqSvc := new(MockEquipmentService)

	_ = setupNATSWrapper(t, natsConn, mockReqSvc, mockEqSvc)

	testID := uuid.New()
	mockReqSvc.On("Cancel", mock.Anything, testID).Return(nil)

	reqData := map[string]interface{}{
		"request_id": testID.String(),
	}

	data, err := json.Marshal(reqData)
	require.NoError(t, err)

	msg, err := natsConn.Request("requests.cancel", data, 2*time.Second)
	require.NoError(t, err)

	var resp map[string]interface{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp["message"])
	assert.False(t, resp["isError"].(bool))
	mockReqSvc.AssertExpectations(t)
}

func TestHandleAddEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	mockReqSvc := new(MockRequestService)
	mockEqSvc := new(MockEquipmentService)

	_ = setupNATSWrapper(t, natsConn, mockReqSvc, mockEqSvc)

	mockEqSvc.On("Add", mock.Anything, mock.MatchedBy(func(eqs []domain.Equipment) bool {
		return len(eqs) == 2
	})).Return(nil)

	reqData := map[string]interface{}{
		"equipments": []map[string]interface{}{
			{"id": 1, "name": "Equipment 1", "quantity": 5},
			{"id": 2, "name": "Equipment 2", "quantity": 3},
		},
	}

	data, err := json.Marshal(reqData)
	require.NoError(t, err)

	msg, err := natsConn.Request("equipment.add", data, 2*time.Second)
	require.NoError(t, err)

	var resp map[string]interface{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp["message"])
	assert.False(t, resp["isError"].(bool))
	mockEqSvc.AssertExpectations(t)
}

// Тест на обработку ошибок валидации
func TestHandleCreateRequestValidationError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	natsConn, cleanup := setupNATSContainer(t)
	defer cleanup()

	mockReqSvc := new(MockRequestService)
	mockEqSvc := new(MockEquipmentService)

	_ = setupNATSWrapper(t, natsConn, mockReqSvc, mockEqSvc)

	// Невалидные данные - отсутствует обязательное поле schedule_time
	reqData := map[string]interface{}{
		"text":        "Test request",
		"telegram_id": int64(123456),
		"address":     "Test Address",
	}

	data, err := json.Marshal(reqData)
	require.NoError(t, err)

	msg, err := natsConn.Request("requests.create", data, 2*time.Second)
	require.NoError(t, err)

	var resp map[string]interface{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	assert.Equal(t, "validation error", resp["message"])
	assert.True(t, resp["isError"].(bool))
	assert.Equal(t, float64(400), resp["statusCode"])
}

// Вспомогательная функция
func stringPtr(s string) *string {
	return &s
}
