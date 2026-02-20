package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
)

func (r *Router) InitRequestRoutes() *Router {
	group := r.r.Group("/api/v1/requests")
	group.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	group.Get("/", r.listRequests())
	group.Get(":id", r.getRequest())
	group.Patch(":id", r.updateRequest())
	group.Post(":id/responsible", r.assignResponsible())
	group.Get(":id/history", r.listRequestHistory())
	group.Post(":id/history", r.addRequestHistory())
	return r
}

func (r *Router) listRequests() fiber.Handler {
	return func(c fiber.Ctx) error {
		var filters models.RequestFilter
		if responsible := c.Query("responsible_id"); responsible != "" {
			filters.ResponsibleID = &responsible
		}
		ctx := r.userContext(c)
		var responsibleUUID *uuid.UUID
		if filters.ResponsibleID != nil {
			id, err := uuid.Parse(*filters.ResponsibleID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid responsible_id"})
			}
			responsibleUUID = &id
		}
		requests, err := r.requestSvc.List(ctx, responsibleUUID)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := models.RequestListResponse{Items: make([]models.RequestResponse, 0, len(requests))}
		for _, req := range requests {
			resp.Items = append(resp.Items, toRequestResponse(req))
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

func (r *Router) getRequest() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request id"})
		}
		ctx := r.userContext(c)
		req, err := r.requestSvc.Get(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toRequestResponse(req))
	}
}

func (r *Router) updateRequest() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request id"})
		}
		var body models.RequestUpdateRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}

		// Преобразуем models.RequestUpdateRequest в domain.RequestUpdate
		updates := domain.RequestUpdate{
			RequestText:  body.RequestText,
			ScheduleTime: body.ScheduleTime,
			Address:      body.Address,
		}

		if body.Status != nil {
			status := domain.RequestStatus(*body.Status)
			updates.Status = &status
		}

		if body.ResponsibleID != nil {
			responsibleID, err := uuid.Parse(*body.ResponsibleID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid responsible_id"})
			}
			updates.ResponsibleID = &responsibleID
		}

		ctx := r.userContext(c)
		req, err := r.requestSvc.UpdateRequest(ctx, id, updates)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toRequestResponse(req))
	}
}

func (r *Router) assignResponsible() fiber.Handler {
	type payload struct {
		ResponsibleID string `json:"responsible_id" validate:"required,uuid4"`
	}
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request id"})
		}
		var body payload
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		responsibleID, err := uuid.Parse(body.ResponsibleID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid responsible_id"})
		}
		ctx := r.userContext(c)
		req, err := r.requestSvc.AssignResponsible(ctx, id, responsibleID)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toRequestResponse(req))
	}
}

func (r *Router) listRequestHistory() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request id"})
		}
		ctx := r.userContext(c)
		history, err := r.historySvc.List(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := models.RequestStatusHistoryListResponse{Items: make([]models.RequestStatusHistoryResponse, 0, len(history))}
		for _, entry := range history {
			resp.Items = append(resp.Items, toHistoryResponse(entry))
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

func (r *Router) addRequestHistory() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request id"})
		}
		var body models.RequestStatusUpdateRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		ctx := r.userContext(c)
		entry := domain.RequestStatusHistory{
			RequestID: id,
			Status:    domain.RequestStatus(body.Status),
			Comment:   body.Comment,
		}
		created, err := r.historySvc.Add(ctx, entry)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(toHistoryResponse(created))
	}
}

// Response mappers

func toRequestResponse(req domain.Request) models.RequestResponse {
	resp := models.RequestResponse{
		ID:           req.ID.String(),
		RequestText:  req.RequestText,
		Status:       string(req.Status),
		ScheduleTime: req.ScheduleTime,
		Address:      req.Address,
		CreatedAt:    req.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    req.Audit.UpdatedAt.Format(time.RFC3339),
	}
	if !req.EndTime.IsZero() {
		resp.EndTime = req.EndTime.Format(time.RFC3339)
	}
	if req.ResponsibleUserID != nil {
		id := req.ResponsibleUserID.String()
		resp.ResponsibleID = &id
	}
	if req.Audit.CreatedBy != nil {
		id := req.Audit.CreatedBy.String()
		resp.CreatedBy = &id
	}
	if req.Audit.UpdatedBy != nil {
		id := req.Audit.UpdatedBy.String()
		resp.UpdatedBy = &id
	}
	if req.Audit.DeletedAt != nil {
		dt := req.Audit.DeletedAt.Format(time.RFC3339)
		resp.DeletedAt = &dt
	}
	resp.Equipment = make([]models.Equipment, 0, len(req.Equipment))
	for _, eq := range req.Equipment {
		resp.Equipment = append(resp.Equipment, models.Equipment{
			ID:        int(eq.EquipmentID),
			Quantity:  eq.Quantity,
			CreatedAt: eq.CreatedAt.Format(time.RFC3339),
			UpdatedAt: eq.UpdatedAt.Format(time.RFC3339),
		})
	}
	return resp
}

func toHistoryResponse(entry domain.RequestStatusHistory) models.RequestStatusHistoryResponse {
	resp := models.RequestStatusHistoryResponse{
		ID:        entry.ID,
		RequestID: entry.RequestID.String(),
		Status:    string(entry.Status),
		Comment:   entry.Comment,
		ChangedAt: entry.ChangedAt.Format(time.RFC3339),
	}
	if entry.ChangedBy != nil {
		id := entry.ChangedBy.String()
		resp.ChangedBy = &id
	}
	return resp
}
