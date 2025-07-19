package main

import (
	"technoBro/internal/router"
	"technoBro/pkg/config"
)

type Config struct {
	Router router.Config `yaml:"router" env:"ROUTER"`
}

func main() {
	cfg := config.MustParseConfig[Config]("config.yaml")

	r := router.NewRouter(cfg.Router)
	r.InitMiddlewares()
	r.Run()
}
