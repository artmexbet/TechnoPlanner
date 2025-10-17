package app

import (
	"gateway/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/mailru/easyjson"
)

func (r *Router) InitUserRoutes() *Router {
	users := r.r.Group("/api/v1/users")

	users.Post("/register", r.RegisterUser()) // регистрация пользователя
	users.Post("/login", r.LoginUser())       // логин пользователя
	users.Post("/refresh", r.RefreshToken())  // обновление токенов

	return r
}

func (r *Router) RegisterUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req models.RegisterRequest
		err := easyjson.Unmarshal(ctx.Body(), &req)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "invalid request body",
				Details: err.Error(),
			})
		}

		if err = r.validator.Struct(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "validation failed",
				Details: err.Error(),
			})
		}

		userID, err := r.authSvc.Register(ctx.Context(), req.Username, req.Password, req.Email)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "could not register user",
				Details: err.Error(),
			})
		}

		return ctx.Status(fiber.StatusCreated).JSON(models.RegisterResponse{
			UserID: userID,
		})
	}
}

func (r *Router) LoginUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req models.LoginRequest
		err := easyjson.Unmarshal(ctx.Body(), &req)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "invalid request body",
				Details: err.Error(),
			})
		}

		if err = r.validator.Struct(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "validation failed",
				Details: err.Error(),
			})
		}

		// Capture additional info
		req.IP = ctx.IP()
		req.UserAgent = string(ctx.Request().Header.UserAgent())
		req.DeviceID = string(ctx.Request().Header.Host())

		tokens, err := r.authSvc.Login(ctx.Context(), req)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "invalid credentials",
				Details: err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(tokens)
	}
}

func (r *Router) RefreshToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req models.TokenRefreshRequest
		err := easyjson.Unmarshal(ctx.Body(), &req)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "invalid request body",
				Details: err.Error(),
			})
		}

		if err = r.validator.Struct(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "validation failed",
				Details: err.Error(),
			})
		}

		// Capture additional info
		req.IP = ctx.IP()
		req.UserAgent = string(ctx.Request().Header.UserAgent())
		req.DeviceID = string(ctx.Request().Header.Host())

		tokens, err := r.authSvc.Refresh(ctx.Context(), req)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "could not refresh token",
				Details: err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(tokens)
	}
}
