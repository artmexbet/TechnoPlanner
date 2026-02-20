package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
)

func (r *Router) InitProtectedUserRoutes() *Router {
	users := r.r.Group("/api/v1/users")
	users.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	users.Post("/logout", r.LogoutUser())

	return r
}

func (r *Router) LogoutUser() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		token := keyauth.TokenFromContext(ctx)
		err := r.authSvc.Logout(ctx.RequestCtx(), token)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "could not logout user",
				Details: err.Error(),
			})
		}
		return nil
	}
}

// fiber:context-methods migrated
