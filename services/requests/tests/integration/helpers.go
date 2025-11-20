package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pgContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestPostgres(ctx context.Context, t *testing.T) (*pgContainer.PostgresContainer, *pgxpool.Pool) {
	t.Helper()

	container, err := pgContainer.Run(ctx,
		"postgres:16",
		pgContainer.WithDatabase("requests_test"),
		pgContainer.WithUsername("requests"),
		pgContainer.WithPassword("requests"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	require.NoError(t, err)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	return container, pool
}

func setupTestNats(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "nats:2",
		ExposedPorts: []string{"4222/tcp"},
		WaitingFor:   wait.ForListeningPort("4222/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "4222/tcp")
	require.NoError(t, err)

	url := fmt.Sprintf("nats://%s:%s", host, port.Port())
	return container, url
}

func newNatsConn(ctx context.Context, t *testing.T, url string) *nats.Conn {
	t.Helper()

	conn, err := nats.Connect(url)
	require.NoError(t, err)
	return conn
}

func applyMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	entries, err := os.ReadDir("migrations")
	require.NoError(t, err)

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join("migrations", entry.Name()))
		require.NoError(t, err)
		_, err = pool.Exec(ctx, string(data))
		require.NoError(t, err)
	}
}

func insertEquipment(t *testing.T, ctx context.Context, pool *pgxpool.Pool, name string) int {
	t.Helper()

	var id int
	err := pool.QueryRow(ctx,
		"INSERT INTO equipment (name, description, quantity) VALUES ($1, $2, $3) RETURNING id",
		name,
		"available",
		10,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func strPtr(value string) *string {
	return &value
}
