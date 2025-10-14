package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type techManagerSvc interface {
	// TODO: определить методы сервиса техники (List, Get, Create, ...)
}

type taskManagerSvc interface {
	// TODO: определить методы сервиса задач (List, Get, Create, ...)
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

// InitRoutes регистрирует HTTP-маршруты
func (r *Router) InitRoutes() {
	// Healthcheck
	r.r.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := r.r.Group("/api")
	v1 := api.Group("/v1")

	// Простая проверка доступности API
	v1.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	// Техника
	tech := v1.Group("/technic")
	tech.Get("/", func(c *fiber.Ctx) error {
		// TODO: использовать r.techMngr для получения списка техники
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "not implemented"})
	})
	tech.Get("/:id", func(c *fiber.Ctx) error {
		// TODO: использовать r.techMngr для получения техники по id
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "not implemented"})
	})

	// Задачи
	tasks := v1.Group("/tasks")
	tasks.Get("/", func(c *fiber.Ctx) error {
		// TODO: использовать r.taskMngr для получения списка задач
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "not implemented"})
	})
	tasks.Get("/:id", func(c *fiber.Ctx) error {
		// TODO: использовать r.taskMngr для получения задачи по id
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "not implemented"})
	})
}

func (r *Router) Run() {
	if err := r.r.Listen(fmt.Sprintf("%s:%s", r.cfg.Address, r.cfg.Port)); err != nil {
		panic(err)
	}
}
