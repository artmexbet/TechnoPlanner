package main

import (
	"net"

	"auth/internal/postgres"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"auth/internal/storeredis"
	"proto"

	"config"

	"google.golang.org/grpc"
)

type Config struct {
	Repository repository.Config `yaml:"repository" env:"REPOSITORY"`
	Redis      storeredis.Config `yaml:"redis" env:"REDIS"`
	Postgres   postgres.Config   `yaml:"postgres" env:"POSTGRES"`

	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET"`
}

func main() {
	cfg := config.MustParseConfig[Config]("./cmd/config/cfg.yaml")

	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	redisClient, err := storeredis.New(cfg.Redis)
	if err != nil {
		panic(err)
	}

	repo, err := repository.New(cfg.Repository, redisClient, pg)
	if err != nil {
		panic(err)
	}

	gen := service.NewTokenizer(cfg.Repository.AccessTokenTTL, cfg.Repository.RefreshTokenTTL, cfg.JWTSecret)

	svc := service.NewAuth(gen, repo)
	handler := server.NewHandler(svc)

	grpcServer := grpc.NewServer() //todo: накинуть интерцептор логирования
	proto.RegisterAuthServer(grpcServer, handler)
	grpcListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	if err = grpcServer.Serve(grpcListener); err != nil {
		panic(err) // TODO: handle error properly
	}
}
