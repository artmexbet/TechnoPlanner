package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
)

func (r *Router) InitResponsibleRoutes() *Router {
	group := r.r.Group("/api/v1/responsibles")
	group.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	group.Get("/", r.listResponsibles())
	group.Post("/", r.createResponsible())
	return r
}

func (r *Router) listResponsibles() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := r.userContext(c)
		responsibles, err := r.responsibleSvc.List(ctx)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := models.ResponsibleListResponse{Items: make([]models.ResponsibleResponse, 0, len(responsibles))}
		for _, responsible := range responsibles {
			resp.Items = append(resp.Items, models.ResponsibleResponse{
				ID:       responsible.ID.String(),
				Username: responsible.Username,
			})
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

func (r *Router) createResponsible() fiber.Handler {
	return func(c fiber.Ctx) error {
		var body models.ResponsibleCreateRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		id, err := uuid.Parse(body.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid id"})
		}
		ctx := r.userContext(c)
		responsible, err := r.responsibleSvc.Create(ctx, id, body.Username)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := models.ResponsibleResponse{
			ID:       responsible.ID.String(),
			Username: responsible.Username,
		}
		return c.Status(fiber.StatusCreated).JSON(resp)
	}
}
