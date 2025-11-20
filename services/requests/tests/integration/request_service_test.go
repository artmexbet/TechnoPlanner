package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"

	"requests/internal/domain"
	"requests/internal/postgres"
	"requests/internal/repository"
	requestservice "requests/internal/service/request"
	"requests/internal/wrapnats"
)

func TestRequestServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	pgContainer, pgPool := setupTestPostgres(ctx, t)
	defer func() {
		pgPool.Close()
		require.NoError(t, pgContainer.Terminate(ctx))
	}()

	applyMigrations(t, ctx, pgPool)

	equipID := insertEquipment(t, ctx, pgPool, "Crane")

	natsContainer, natsURL := setupTestNats(ctx, t)
	defer func() {
		require.NoError(t, natsContainer.Terminate(ctx))
	}()

	natsConn := newNatsConn(ctx, t, natsURL)
	defer natsConn.Close()

	publisher := wrapnats.NewNatsPublisher(natsConn)

	repo := repository.NewRepository(postgres.New(pgPool), publisher)
	svc := requestservice.New(repo)

	payloadReceived := make(chan struct{}, 1)
	sub, err := natsConn.Subscribe("events.request.created", func(_ *nats.Msg) {
		select {
		case payloadReceived <- struct{}{}:
		default:
		}
	})
	require.NoError(t, err)
	defer require.NoError(t, sub.Unsubscribe())

	reqText := "integration request"
	req, err := svc.Add(ctx, domain.Request{
		RequestText:  &reqText,
		ScheduleTime: "2025-11-25T10:00:00Z",
		Address:      "Test facility",
		Issuer: domain.User{
			TelegramID: 123456789,
			Username:   "int-bot",
			FirstName:  "Ivan",
			LastName:   strPtr("Petrov"),
		},
		Equipments: []domain.Equipment{{ID: equipID, Quantity: 1}},
	})
	require.NoError(t, err)
	require.NotZero(t, req.ID)
	require.Len(t, req.Equipments, 1)
	require.Equal(t, equipID, req.Equipments[0].ID)

	require.Eventually(t, func() bool {
		select {
		case <-payloadReceived:
			return true
		default:
			return false
		}
	}, time.Second, 50*time.Millisecond)

	listed, err := svc.List(ctx, req.Issuer, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, listed)

	fetched, err := svc.Get(ctx, req.ID)
	require.NoError(t, err)
	require.Equal(t, req.ID, fetched.ID)
	require.Equal(t, req.Address, fetched.Address)
}
