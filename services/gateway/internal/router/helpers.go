package router

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
)

// handleServiceError обрабатывает ошибки сервисов и возвращает соответствующий HTTP-ответ
func handleServiceError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "forbidden"})
	case errors.Is(err, domain.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "not found"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal error", Details: err.Error()})
	}
}

// derefString безопасно разыменовывает указатель на строку
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
