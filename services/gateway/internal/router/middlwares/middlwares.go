package middlwares

import (
	"context"

	"gateway/internal/models"

	"github.com/gofiber/fiber/v2"
)

type iAuthClient interface {
	ValidateToken(ctx context.Context, token string) (models.TokenValidationResponse, error)
}

func CheckJWTMiddleware(authClient iAuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.GetReqHeaders()["Authorization"]
		if len(h) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "missing authorization header",
			})
		}
		token := h[0][len("Bearer "):]
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

		c.Locals("user_id", validateRes.UserID)
		c.Locals("token", token)
		return c.Next()
	}
}
