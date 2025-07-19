package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Config struct {
	Address string `yaml:"address" env:"ADDRESS"`
	Port    string `yaml:"port" env:"PORT"`
}

type Router struct {
	r   *fiber.App
	cfg Config
}

func NewRouter(cfg Config) *Router {
	return &Router{
		r:   fiber.New(),
		cfg: cfg,
	}
}

func (r *Router) InitMiddlewares() {
	r.r.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
		},
	))
	r.r.Use(recover.New())
	r.r.Use(requestid.New())
}

func (r *Router) Run() {
	// extract the config from the environment or a config file
	if err := r.r.Listen(fmt.Sprintf("%s:%s", r.cfg.Address, r.cfg.Port)); err != nil {
		panic(err)
	}
}
