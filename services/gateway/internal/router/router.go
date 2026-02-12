package router

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	slogfiber "github.com/samber/slog-fiber"
	"go.opentelemetry.io/otel/trace"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/router/middlwares"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/service"
)

type UserService interface {
	// Define methods for user service operations
}

type AuthSvcConnector interface {
	Login(ctx context.Context, req models.LoginRequest) (models.TokenPair, error)
	Register(ctx context.Context, username, password, email string) (string, error)
	ValidateToken(ctx context.Context, token string) (models.TokenValidationResponse, error)
	Refresh(ctx context.Context, req models.TokenRefreshRequest) (models.TokenPair, error)
	Logout(ctx context.Context, token string) error
	LogoutAll(ctx context.Context, token string) error
}

type PorterService interface {
	List(ctx context.Context) ([]domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetCurrentUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	Create(ctx context.Context, username, email, password string) (string, error)
}

type EquipmentService interface {
	Create(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	Update(ctx context.Context, eq domain.Equipment) (domain.Equipment, error)
	Get(ctx context.Context, id int32) (domain.Equipment, error)
	List(ctx context.Context) ([]domain.Equipment, error)
	Delete(ctx context.Context, id int32) error
}

type CategoryService interface {
	Create(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	Update(ctx context.Context, cat domain.EquipmentCategory) (domain.EquipmentCategory, error)
	List(ctx context.Context) ([]domain.EquipmentCategory, error)
	Delete(ctx context.Context, id int32) error
}

type RequestService interface {
	List(ctx context.Context, responsibleID *uuid.UUID) ([]domain.Request, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Request, error)
	AssignResponsible(ctx context.Context, requestID uuid.UUID, responsibleID uuid.UUID) (domain.Request, error)
}

type HistoryService interface {
	List(ctx context.Context, requestID uuid.UUID) ([]domain.RequestStatusHistory, error)
	Add(ctx context.Context, entry domain.RequestStatusHistory) (domain.RequestStatusHistory, error)
}

type Config struct {
	Address string `yaml:"address" env:"ADDRESS"`
	Port    string `yaml:"port" env:"PORT"`
}

type Router struct {
	r         *fiber.App
	validator *validator.Validate

	cfg          Config
	userSvc      UserService
	authSvc      AuthSvcConnector
	porterSvc    PorterService
	equipmentSvc EquipmentService
	categorySvc  CategoryService
	requestSvc   RequestService
	historySvc   HistoryService
}

func NewRouter(
	cfg Config,
	userSvc UserService,
	authSvc AuthSvcConnector,
	porterSvc PorterService,
	equipmentSvc EquipmentService,
	categorySvc CategoryService,
	requestSvc RequestService,
	historySvc HistoryService,
) *Router {
	return &Router{
		r:            fiber.New(),
		validator:    validator.New(validator.WithRequiredStructEnabled()),
		cfg:          cfg,
		userSvc:      userSvc,
		authSvc:      authSvc,
		porterSvc:    porterSvc,
		equipmentSvc: equipmentSvc,
		categorySvc:  categorySvc,
		requestSvc:   requestSvc,
		historySvc:   historySvc,
	}
}

func (r *Router) InitMiddlewares(provider trace.TracerProvider) *Router {
	r.r.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
		},
	))
	r.r.Use(recover.New())
	r.r.Use(requestid.New()) // Trace
	r.r.Use(
		otelfiber.Middleware(
			otelfiber.WithTracerProvider(provider),
		),
	)
	r.r.Use(slogfiber.New(slog.Default()))
	return r
}

// InitBaseRoutes регистрирует HTTP-маршруты
func (r *Router) InitBaseRoutes() *Router {
	// Healthcheck
	r.r.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Простая проверка доступности API
	r.r.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	return r
}

func (r *Router) Run() {
	if err := r.r.Listen(fmt.Sprintf("%s:%s", r.cfg.Address, r.cfg.Port)); err != nil {
		panic(err)
	}
}

func (r *Router) userContext(c *fiber.Ctx) context.Context {
	userID, _ := c.Locals(middlwares.ContextUserIDKey).(string)
	role, _ := c.Locals(middlwares.ContextUserRoleKey).(string)
	return service.WithUserContext(c.UserContext(), userID, role)
}

func (r *Router) InitPorterRoutes() *Router {
	group := r.r.Group("/api/v1/porters")
	group.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	group.Get("/", r.listPorters())
	group.Post("/", r.createPorter())
	group.Get(":id", r.getPorter())
	return r
}

func (r *Router) listPorters() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := r.userContext(c)
		users, err := r.porterSvc.List(ctx)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := models.PorterListResponse{Items: make([]models.PorterResponse, 0, len(users))}
		for _, u := range users {
			resp.Items = append(resp.Items, toPorterResponse(u))
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

func (r *Router) getPorter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid porter id"})
		}
		ctx := r.userContext(c)
		user, err := r.porterSvc.Get(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toPorterResponse(user))
	}
}

func (r *Router) createPorter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.PorterCreateRequest
		if err := c.BodyParser(&req); err != nil {
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

func (r *Router) listEquipment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := r.userContext(c)
		items, err := r.equipmentSvc.List(ctx)
		if err != nil {
			return handleServiceError(c, err)
		}
		resp := make([]models.Equipment, 0, len(items))
		for _, eq := range items {
			resp = append(resp, toEquipmentResponse(eq))
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": resp})
	}
}

func (r *Router) getEquipment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseInt32Param(c, "id")
		if err != nil {
			return err
		}
		ctx := r.userContext(c)
		eq, err := r.equipmentSvc.Get(ctx, id)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toEquipmentResponse(eq))
	}
}

func (r *Router) createEquipment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.EquipmentCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		eq := domain.Equipment{
			Name:        req.Name,
			Description: req.Description,
			Quantity:    req.Quantity,
			Categories:  make([]domain.EquipmentCategory, 0, len(req.CategoryIDs)),
		}
		for _, catID := range req.CategoryIDs {
			eq.Categories = append(eq.Categories, domain.EquipmentCategory{ID: catID})
		}
		ctx := r.userContext(c)
		created, err := r.equipmentSvc.Create(ctx, eq)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(toEquipmentResponse(created))
	}
}

func (r *Router) updateEquipment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseInt32Param(c, "id")
		if err != nil {
			return err
		}
		var req models.EquipmentUpdateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		eq := domain.Equipment{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			Quantity:    req.Quantity,
			Categories:  make([]domain.EquipmentCategory, 0, len(req.CategoryIDs)),
		}
		for _, cID := range req.CategoryIDs {
			eq.Categories = append(eq.Categories, domain.EquipmentCategory{ID: cID})
		}
		ctx := r.userContext(c)
		updated, err := r.equipmentSvc.Update(ctx, eq)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toEquipmentResponse(updated))
	}
}

func (r *Router) deleteEquipment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseInt32Param(c, "id")
		if err != nil {
			return err
		}
		ctx := r.userContext(c)
		if err := r.equipmentSvc.Delete(ctx, id); err != nil {
			return handleServiceError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (r *Router) listCategories() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	return func(c *fiber.Ctx) error {
		var req models.EquipmentCategoryCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		cat := domain.EquipmentCategory{Name: req.Name, Description: req.Description}
		ctx := r.userContext(c)
		created, err := r.categorySvc.Create(ctx, cat)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(toCategoryResponse(created))
	}
}

func (r *Router) updateCategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseInt32Param(c, "id")
		if err != nil {
			return err
		}
		var req models.EquipmentCategoryUpdateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid body", Details: err.Error()})
		}
		if err := r.validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "validation failed", Details: err.Error()})
		}
		cat := domain.EquipmentCategory{ID: id, Name: req.Name, Description: req.Description}
		ctx := r.userContext(c)
		updated, err := r.categorySvc.Update(ctx, cat)
		if err != nil {
			return handleServiceError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(toCategoryResponse(updated))
	}
}

func (r *Router) deleteCategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseInt32Param(c, "id")
		if err != nil {
			return err
		}
		ctx := r.userContext(c)
		if err := r.categorySvc.Delete(ctx, id); err != nil {
			return handleServiceError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (r *Router) InitRequestRoutes() *Router {
	group := r.r.Group("/api/v1/requests")
	group.Use(middlwares.CheckJWTMiddleware(r.authSvc))
	group.Get("/", r.listRequests())
	group.Get(":id", r.getRequest())
	group.Post(":id/responsible", r.assignResponsible())
	group.Get(":id/history", r.listRequestHistory())
	group.Post(":id/history", r.addRequestHistory())
	return r
}

func (r *Router) listRequests() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	return func(c *fiber.Ctx) error {
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

func (r *Router) assignResponsible() fiber.Handler {
	type payload struct {
		ResponsibleID string `json:"responsible_id" validate:"required,uuid4"`
	}
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request id"})
		}
		var body payload
		if err := c.BodyParser(&body); err != nil {
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
	return func(c *fiber.Ctx) error {
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
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request id"})
		}
		var body models.RequestStatusUpdateRequest
		if err := c.BodyParser(&body); err != nil {
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

func toPorterResponse(u domain.User) models.PorterResponse {
	return models.PorterResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

func toEquipmentResponse(eq domain.Equipment) models.Equipment {
	resp := models.Equipment{
		ID:          eq.ID,
		Name:        eq.Name,
		Description: eq.Description,
		Quantity:    eq.Quantity,
		CreatedAt:   eq.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   eq.Audit.UpdatedAt.Format(time.RFC3339),
	}
	if eq.Audit.CreatedBy != nil {
		id := eq.Audit.CreatedBy.String()
		resp.CreatedBy = &id
	}
	if eq.Audit.UpdatedBy != nil {
		id := eq.Audit.UpdatedBy.String()
		resp.UpdatedBy = &id
	}
	if eq.Audit.DeletedAt != nil {
		dt := eq.Audit.DeletedAt.Format(time.RFC3339)
		resp.DeletedAt = &dt
	}
	resp.Categories = make([]models.EquipmentCategory, 0, len(eq.Categories))
	for _, cat := range eq.Categories {
		resp.Categories = append(resp.Categories, toCategoryResponse(cat))
	}
	return resp
}

func toCategoryResponse(cat domain.EquipmentCategory) models.EquipmentCategory {
	resp := models.EquipmentCategory{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		CreatedAt:   cat.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   cat.Audit.UpdatedAt.Format(time.RFC3339),
	}
	if cat.Audit.CreatedBy != nil {
		id := cat.Audit.CreatedBy.String()
		resp.CreatedBy = &id
	}
	if cat.Audit.UpdatedBy != nil {
		id := cat.Audit.UpdatedBy.String()
		resp.UpdatedBy = &id
	}
	if cat.Audit.DeletedAt != nil {
		dt := cat.Audit.DeletedAt.Format(time.RFC3339)
		resp.DeletedAt = &dt
	}
	return resp
}

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
			ID:        eq.EquipmentID,
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

func handleServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "forbidden"})
	case errors.Is(err, domain.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "not found"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal error", Details: err.Error()})
	}
}

func parseInt32Param(c *fiber.Ctx, name string) (int32, error) {
	v, err := strconv.ParseInt(c.Params(name), 10, 32)
	if err != nil {
		return 0, c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid id"})
	}
	return int32(v), nil
}
