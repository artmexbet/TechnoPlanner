package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mailru/easyjson"
)

func (r *Router) InitUserRoutes() *Router {
	users := r.r.Group("/api/v1/users")

	users.Post("/register", r.RegisterUser()) // регистрация пользователя
	users.Post("/login", r.LoginUser())       // логин пользователя

	return r
}

func (r *Router) RegisterUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req UserRequest
		err := easyjson.Unmarshal(ctx.Body(), &req)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if err = r.validator.Struct(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		if err = r.authSvc.Register(req.Username, req.Password, req.Email); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "could not register user",
				"details": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "user registered successfully",
		})
	}
}

func (r *Router) LoginUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
