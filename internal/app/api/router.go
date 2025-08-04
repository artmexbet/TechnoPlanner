package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type techManagerSvc interface {
}

type taskManagerSvc interface {
}

type Config struct {
	Address string `yaml:"address" env:"ADDRESS"`
	Port    string `yaml:"port" env:"PORT"`
}

type Router struct {
	r        *fiber.App
	cfg      Config
	techMngr techManagerSvc
	taskMngr taskManagerSvc
}

func NewRouter(cfg Config, techMngr techManagerSvc, taskMngr taskManagerSvc) *Router {
	return &Router{
		r:        fiber.New(),
		cfg:      cfg,
		techMngr: techMngr,
		taskMngr: taskMngr,
	}
}

func (r *Router) InitMiddlewares() {
	r.r.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
		},
	))
	r.r.Use(recover.New())
	r.r.Use(requestid.New()) // Trace
}

func (r *Router) Run() {
	if err := r.r.Listen(fmt.Sprintf("%s:%s", r.cfg.Address, r.cfg.Port)); err != nil {
		panic(err)
	}
}
