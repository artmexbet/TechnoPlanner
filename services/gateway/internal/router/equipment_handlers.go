package router

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
)

func (r *Router) InitEquipmentRoutes() *Router {
	group := r.r.Group("/api/v1/equipment")
	group.Use(middlwares.CheckJWTMiddleware(r.authSvc))

	cat := group.Group("/categories")

	cat.Get("/", r.listCategories())
	cat.Post("/", r.createCategory())
	cat.Put(":id", r.updateCategory())
	cat.Delete(":id", r.deleteCategory())

	group.Get("/", r.listEquipment())
	group.Post("/", r.createEquipment())
	group.Get(":id", r.getEquipment())
	group.Put(":id", r.updateEquipment())
	group.Delete(":id", r.deleteEquipment())
	return r
}

// Equipment handlers

func (r *Router) listEquipment() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := r.userContext(c)
		items, err := r.equipmentSvc.List(ctx)
		if err != nil {
			return handleServiceError(c, err)
		}
		categoryMap := r.loadCategoryMap(ctx)
		resp := make([]models.Equipment, 0, len(items))
		for _, eq := range items {
			resp = append(resp, toEquipmentResponse(eq, categoryMap))
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": resp})
	}
}

func (r *Router) getEquipment() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid equipment id"})
		}
		ctx := r.userContext(c)
		eq, err := r.equipmentSvc.Get(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toEquipmentResponse(eq, r.loadCategoryMap(ctx)))
	}
}

func (r *Router) createEquipment() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.EquipmentCreateRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		eq := domain.Equipment{
			Name:        req.Name,
			Description: derefString(req.Description),
			Quantity:    req.Quantity,
		}
		// CategoryID из первого элемента CategoryIDs
		if len(req.CategoryIDs) > 0 {
			eq.CategoryID = req.CategoryIDs[0]
		}
		ctx := r.userContext(c)
		created, err := r.equipmentSvc.Create(ctx, eq)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(toEquipmentResponse(created, r.loadCategoryMap(ctx)))
	}
}

func (r *Router) updateEquipment() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid equipment id"})
		}
		var req models.EquipmentUpdateRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		eq := domain.Equipment{
			ID:          id,
			Name:        req.Name,
			Description: derefString(req.Description),
			Quantity:    req.Quantity,
		}
		// CategoryID из первого элемента CategoryIDs
		if len(req.CategoryIDs) > 0 {
			eq.CategoryID = req.CategoryIDs[0]
		}
		ctx := r.userContext(c)
		updated, err := r.equipmentSvc.Update(ctx, eq)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toEquipmentResponse(updated, r.loadCategoryMap(ctx)))
	}
}

func (r *Router) deleteEquipment() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid equipment id"})
		}
		ctx := r.userContext(c)
		if err := r.equipmentSvc.Delete(ctx, id); err != nil {
			return handleServiceError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Category handlers

func (r *Router) listCategories() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := r.userContext(c)
		cats, err := r.categorySvc.List(ctx)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := models.EquipmentCategoryListResponse{Items: make([]models.EquipmentCategory, 0, len(cats))}
		for _, cat := range cats {
			resp.Items = append(resp.Items, toCategoryResponse(cat))
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

func (r *Router) createCategory() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.EquipmentCategoryCreateRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		cat := domain.EquipmentCategory{Name: req.Name, Description: derefString(req.Description)}
		ctx := r.userContext(c)
		created, err := r.categorySvc.Create(ctx, cat)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(toCategoryResponse(created))
	}
}

func (r *Router) updateCategory() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid category id"})
		}
		var req models.EquipmentCategoryUpdateRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		cat := domain.EquipmentCategory{ID: id, Name: req.Name, Description: derefString(req.Description)}
		ctx := r.userContext(c)
		updated, err := r.categorySvc.Update(ctx, cat)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toCategoryResponse(updated))
	}
}

func (r *Router) deleteCategory() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid category id"})
		}
		ctx := r.userContext(c)
		if err := r.categorySvc.Delete(ctx, id); err != nil {
			return handleServiceError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Response mappers

func (r *Router) loadCategoryMap(ctx context.Context) map[int]domain.EquipmentCategory {
	cats, err := r.categorySvc.List(ctx)
	if err != nil {
		return nil
	}
	categoryMap := make(map[int]domain.EquipmentCategory, len(cats))
	for _, cat := range cats {
		categoryMap[cat.ID] = cat
	}
	return categoryMap
}

func toEquipmentResponse(eq domain.Equipment, categoryMap map[int]domain.EquipmentCategory) models.Equipment {
	resp := models.Equipment{
		ID:          eq.ID,
		Name:        eq.Name,
		Description: &eq.Description,
		Quantity:    eq.Quantity,
		CreatedAt:   eq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   eq.UpdatedAt.Format(time.RFC3339),
	}
	if categoryMap != nil {
		if cat, ok := categoryMap[eq.CategoryID]; ok {
			resp.Categories = []models.EquipmentCategory{toCategoryResponse(cat)}
		}
	}
	return resp
}

func toCategoryResponse(cat domain.EquipmentCategory) models.EquipmentCategory {
	resp := models.EquipmentCategory{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: &cat.Description,
		CreatedAt:   cat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   cat.UpdatedAt.Format(time.RFC3339),
	}
	return resp
}
