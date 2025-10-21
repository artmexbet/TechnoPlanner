package router

import (
	"gateway/internal/models"
	"gateway/internal/router/middlwares"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) InitProtectedUserRoutes() *Router {
	users := r.r.Group("/api/v1/users")
	users.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	users.Post("/logout", r.LogoutUser())

	return r
}

func (r *Router) LogoutUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Locals("token").(string)
		err := r.authSvc.Logout(ctx.Context(), token)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "could not logout user",
				Details: err.Error(),
			})
		}
		return nil
	}
}
