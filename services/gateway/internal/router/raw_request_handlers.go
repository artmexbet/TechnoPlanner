package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
)

// InitRawRequestRoutes регистрирует маршруты для сырых запросов от бота
func (r *Router) InitRawRequestRoutes() *Router {
	group := r.r.Group("/api/v1/raw-requests")
	group.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	group.Get("/", r.listRawRequests())
	group.Get("/:id", r.getRawRequest())
	group.Post("/:id/process", r.processRawRequest())
	return r
}

// listRawRequests возвращает список сырых запросов от бота
func (r *Router) listRawRequests() fiber.Handler {
	return func(c fiber.Ctx) error {
		status := c.Query("status", "") // "new", "processed" или "" (все)
		ctx := r.userContext(c)

		requests, err := r.rawRequestSvc.List(ctx, status, 50, 0)
		if err != nil {
			return handleServiceError(c, err)
		}

		resp := models.RawRequestListResponse{Items: make([]models.RawRequestResponse, 0, len(requests))}
		for _, req := range requests {
			resp.Items = append(resp.Items, toRawRequestResponse(req))
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// getRawRequest возвращает сырой запрос по ID
func (r *Router) getRawRequest() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid raw request id"})
		}
		ctx := r.userContext(c)

		req, err := r.rawRequestSvc.Get(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toRawRequestResponse(req))
	}
}

// processRawRequest создаёт нормальную заявку из сырого запроса
func (r *Router) processRawRequest() fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid raw request id"})
		}

		var body models.RawRequestProcessRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}

		ctx := r.userContext(c)
		createdReq, updatedRaw, err := r.rawRequestSvc.Process(ctx, id, body)
		if err != nil {
			return handleServiceError(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(models.ProcessRawRequestResponse{
			Request:    toRequestResponse(createdReq),
			RawRequest: toRawRequestResponse(updatedRaw),
		})
	}
}

// toRawRequestResponse конвертирует domain.RawRequest в models.RawRequestResponse
func toRawRequestResponse(r domain.RawRequest) models.RawRequestResponse {
	resp := models.RawRequestResponse{
		ID:         r.ID.String(),
		TelegramID: r.TelegramID,
		Username:   r.Username,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		RawText:    r.RawText,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt.Format(time.RFC3339),
	}
	if r.ProcessedRequestID != nil {
		id := r.ProcessedRequestID.String()
		resp.ProcessedRequestID = &id
	}
	return resp
}
