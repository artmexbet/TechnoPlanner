package middlwares

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
)

const (
	ContextUserIDKey   = "user_id"
	ContextUserRoleKey = "user_role"
)

type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (models.TokenValidationResponse, error)
}

func CheckJWTMiddleware(authClient AuthClient) fiber.Handler {
	return func(c fiber.Ctx) error {
		h := c.GetReqHeaders()["Authorization"]
		if len(h) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "missing authorization header",
			})
		}
		authHeader := h[0]
		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid authorization header format",
			})
		}
		token := authHeader[len(bearerPrefix):]
		validateRes, err := authClient.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "invalid token",
				Details: err.Error(),
			})
		}

		if validateRes.State != models.TokenStateValid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid token state",
			})
		}

		c.Locals(ContextUserIDKey, validateRes.UserID)
		c.Locals(ContextUserRoleKey, validateRes.Role)
		c.Locals("token", token)
		return c.Next()
	}
}

// fiber:context-methods migrated
