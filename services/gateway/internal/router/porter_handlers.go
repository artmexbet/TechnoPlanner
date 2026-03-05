package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
)

func (r *Router) InitPorterRoutes() *Router {
	group := r.r.Group("/api/v1/porters")
	group.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	group.Get("/", r.listPorters())
	group.Post("/", r.createPorter())
	group.Get(":id", r.getPorter())
	group.Put(":id", r.updatePorter())
	group.Delete(":id", r.deletePorter())
	return r
}

func (r *Router) listPorters() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := r.userContext(c)
		porters, err := r.porterSvc.List(ctx)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := models.PorterListResponse{Items: make([]models.PorterResponse, 0, len(porters))}
		for _, p := range porters {
			resp.Items = append(resp.Items, toPorterResponseFromDomainPorter(p))
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

func (r *Router) getPorter() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid porter id"})
		}
		ctx := r.userContext(c)
		porter, err := r.porterSvc.Get(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toPorterResponseFromDomainPorter(porter))
	}
}

func (r *Router) createPorter() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.PorterCreateRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		ctx := r.userContext(c)
		userID, err := r.porterSvc.Create(ctx, req.Username, req.Email, req.Password)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": userID})
	}
}

func (r *Router) updatePorter() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid porter id"})
		}
		var req models.PorterUpdateRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		ctx := r.userContext(c)
		user, err := r.porterSvc.Update(ctx, id, req.Username, req.Email)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toPorterResponse(user))
	}
}

func (r *Router) deletePorter() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid porter id"})
		}
		ctx := r.userContext(c)
		if err := r.porterSvc.Delete(ctx, id); err != nil {
			return handleServiceError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func toPorterResponse(u domain.User) models.PorterResponse {
	return models.PorterResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

func toPorterResponseFromDomainPorter(p domain.Porter) models.PorterResponse {
	return models.PorterResponse{
		ID:       p.ID.String(),
		Username: p.Username,
	}
}
