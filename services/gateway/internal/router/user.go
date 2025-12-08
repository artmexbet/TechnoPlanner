package router

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
)

func (r *Router) InitUserRoutes() *Router {
	users := r.r.Group("/api/v1/users")

	users.Post("/register", r.RegisterUser()) // user registration
	users.Post("/login", r.LoginUser())       // user login
	users.Post("/refresh", r.RefreshToken())  // token refresh
	users.Get("/me", middlwares.CheckJWTMiddleware(r.authSvc), r.Me())

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

		userID, err := r.authSvc.Register(ctx.UserContext(), req.Username, req.Password, req.Email)
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
		slog.InfoContext(ctx.UserContext(), "LoginUser handler called", "context", ctx.UserContext())

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

		tokens, err := r.authSvc.Login(ctx.UserContext(), req)
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
		err := req.UnmarshalJSON(ctx.Body())
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

		tokens, err := r.authSvc.Refresh(ctx.UserContext(), req)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "could not refresh token",
				Details: err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(tokens)
	}
}

func (r *Router) Me() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := c.Locals(middlwares.ContextUserIDKey).(string)
		role, _ := c.Locals(middlwares.ContextUserRoleKey).(string)

		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "user not authenticated",
			})
		}

		ctx := r.userContext(c)
		id, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "invalid user id",
			})
		}

		user, err := r.porterSvc.GetCurrentUser(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(models.MeResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Role:     role,
		})
	}
}
