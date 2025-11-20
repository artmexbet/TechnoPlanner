package integration

import (
	"context"
	"log/slog"
	"requests/internal/domain"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	nConntainer "github.com/testcontainers/testcontainers-go/modules/nats"
	"github.com/testcontainers/testcontainers-go/wait"

	"requests/internal/wrapnats"
)

func TestWrapNatsPublishSubscribe(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping NATS integration test in short mode")
	}

	ctx := context.Background()
	container, err := nConntainer.Run(ctx,
		"nats:latest",
		nConntainer.WithUsername("test"),
		nConntainer.WithPassword("test"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("4222/tcp")),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			slog.Error("failed to terminate container", "err", err)
		}
	}()
	require.NoError(t, err)
	defer require.NoError(t, container.Terminate(ctx))

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	natsConn, err := nats.Connect(connStr)
	require.NoError(t, err)
	defer natsConn.Close()

	publisher := wrapnats.NewNatsPublisher(natsConn)
	vc := make(chan struct{})
	setSubCalled := false

	_, err = natsConn.Subscribe("events.request.created", func(msg *nats.Msg) {
		setSubCalled = true
		vc <- struct{}{}
	})
	require.NoError(t, err)

	err = publisher.PublishRequestCreated(domain.Request{})
	require.NoError(t, err)
	<-vc
	require.True(t, setSubCalled)
}
