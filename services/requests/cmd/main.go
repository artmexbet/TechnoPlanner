package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	"requests/internal/postgres"
	"requests/internal/repository"
	"requests/internal/service/equipment"
	"requests/internal/service/request"
	"requests/internal/wrapnats"

	"config"
)

const (
	defaultConfigPath = "configs/config.yaml"
	configPathKey     = "CONFIG_PATH"
)

type Config struct {
	Nats     *wrapnats.Config `yaml:"nats" env-prefix:"NATS_"`
	Postgres *config.Postgres `yaml:"postgres" env-prefix:"POSTGRES_"`
}

func main() {
	cfgPath := os.Getenv(configPathKey)
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}
	cfg := config.MustParseConfig[Config](cfgPath)
	ctx := context.Background()

	conn, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		panic(err)
	}

	publisher := wrapnats.NewNatsPublisher(conn)

	pg := postgres.New(pool)
	repo := repository.NewRepository(pg, publisher)

	requestService := request.New(repo)
	equipmentService := equipment.New(repo)

	wrapper, err := wrapnats.New(cfg.Nats, conn, requestService, equipmentService)
	if err != nil {
		panic(err)
	}
	wrapper.HandleMsgs()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	conn.Close()
	pool.Close()
}
