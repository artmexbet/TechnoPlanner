package wrapnats

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"tech/internal/domain"

	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"
	"github.com/artmexbet/TechnoPlanner/libs/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type TechService interface {
	AddTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error)
	DeleteTechnic(ctx context.Context, techID uuid.UUID) error
	UpdateTechnic(ctx context.Context, technic domain.Technic) (*domain.Technic, error)
	GetTechnicByID(ctx context.Context, techID uuid.UUID) (*domain.Technic, error)
	GetAllTechnics(ctx context.Context) ([]domain.Technic, error)
	GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error)
}

type CategoryService interface {
	AddCategory(ctx context.Context, categoryName string, description string) (*domain.TechnicCategory, error)
	UpdateCategory(ctx context.Context, category domain.TechnicCategory) (*domain.TechnicCategory, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*domain.TechnicCategory, error)
	GetAllCategories(ctx context.Context) ([]domain.TechnicCategory, error)
}

type Config struct {
	URL string `yaml:"url" env:"URL"`

	RequestTimeout time.Duration `yaml:"request_timeout" env:"REQUEST_TIMEOUT" env-default:"5s"`
}

type NatsWrapper struct {
	conn      *nats.Conn
	cfg       *Config
	validator *validator.Validate

	techService     TechService
	categoryService CategoryService
}

func New(cfg *Config, conn *nats.Conn, techService TechService, categoryService CategoryService) (*NatsWrapper, error) {
	return &NatsWrapper{
		conn:            conn,
		validator:       validator.New(validator.WithRequiredStructEnabled()),
		cfg:             cfg,
		techService:     techService,
		categoryService: categoryService,
	}, nil
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]nats.MsgHandler{
		// Technic handlers
		subjects.GatewayTechnicCreate:        w.handleCreateTechnic,
		subjects.GatewayTechnicUpdate:        w.handleUpdateTechnic,
		subjects.GatewayTechnicDelete:        w.handleDeleteTechnic,
		subjects.GatewayTechnicGet:           w.handleGetTechnic,
		subjects.GatewayTechnicList:          w.handleListTechnics,
		subjects.GatewayTechnicGetByCategory: w.handleGetTechnicsByCategory,
		// Category handlers
		subjects.GatewayTechnicCategoryCreate: w.handleCreateCategory,
		subjects.GatewayTechnicCategoryUpdate: w.handleUpdateCategory,
		subjects.GatewayTechnicCategoryDelete: w.handleDeleteCategory,
		subjects.GatewayTechnicCategoryGet:    w.handleGetCategory,
		subjects.GatewayTechnicCategoryList:   w.handleListCategories,
	}

	for subject, handler := range handlers {
		_, err := w.conn.Subscribe(subject, handler)
		if err != nil {
			slog.Error("cannot subscribe to subject", "subject", subject, "error", err)
		}
	}
	return w
}

// Technic handlers

func (w *NatsWrapper) handleCreateTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicCreateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create technic request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	domainReq := domain.Technic{
		CategoryID:                req.CategoryID,
		Name:                      req.Name,
		Description:               req.Description,
		AdditionalCharacteristics: req.AdditionalCharacteristics,
	}

	createdReq, err := w.techService.AddTechnic(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error creating technic", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapTechnicToDTO(*createdReq)
	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleUpdateTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicUpdateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update technic request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	domainReq := domain.Technic{
		ID:                        req.ID,
		CategoryID:                req.CategoryID,
		Name:                      req.Name,
		Description:               req.Description,
		AdditionalCharacteristics: req.AdditionalCharacteristics,
	}

	updatedReq, err := w.techService.UpdateTechnic(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error updating technic", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapTechnicToDTO(*updatedReq)
	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicGetByIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get technic request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	getReq, err := w.techService.GetTechnicByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting technic", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapTechnicToDTO(*getReq)
	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleDeleteTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicDeleteRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling delete technic request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err := w.techService.DeleteTechnic(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error deleting technic", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", "technic deleted"); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleListTechnics(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	technics, err := w.techService.GetAllTechnics(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error listing technics", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := make([]dto.Technic, len(technics))
	for i, t := range technics {
		respDTO[i] = mapTechnicToDTO(t)
	}

	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetTechnicsByCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicGetByCategoryRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get technics by category request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	technics, err := w.techService.GetTechnicByCategory(ctx, req.CategoryID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting technics by category", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := make([]dto.Technic, len(technics))
	for i, t := range technics {
		respDTO[i] = mapTechnicToDTO(t)
	}

	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

// Category handlers

func (w *NatsWrapper) handleCreateCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicCategoryCreateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create category request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	createdReq, err := w.categoryService.AddCategory(ctx, req.Name, req.Description)
	if err != nil {
		slog.ErrorContext(ctx, "error creating category", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapCategoryToDTO(*createdReq)
	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleUpdateCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicCategoryUpdateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update category request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	updatedReq, err := w.categoryService.UpdateCategory(ctx, domain.TechnicCategory{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		slog.ErrorContext(ctx, "error updating category", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapCategoryToDTO(*updatedReq)
	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleDeleteCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicCategoryDeleteRequest
	if err := req.UnmarshalJSON(msg.Data); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling delete category request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err := w.categoryService.DeleteCategory(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error deleting category", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", "category deleted successfully"); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechnicCategoryGetByIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get category request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	category, err := w.categoryService.GetCategoryByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting category", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapCategoryToDTO(*category)
	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleListCategories(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	categories, err := w.categoryService.GetAllCategories(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error listing categories", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := make([]dto.TechnicCategory, len(categories))
	for i, c := range categories {
		respDTO[i] = mapCategoryToDTO(c)
	}

	respData, err := json.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

// Helper functions

func mapTechnicToDTO(t domain.Technic) dto.Technic {
	return dto.Technic{
		ID:                        t.ID,
		CategoryID:                t.CategoryID,
		Name:                      t.Name,
		Description:               t.Description,
		AdditionalCharacteristics: t.AdditionalCharacteristics,
		CreatedAt:                 t.CreatedAt,
		UpdatedAt:                 t.UpdatedAt,
	}
}

func mapCategoryToDTO(c domain.TechnicCategory) dto.TechnicCategory {
	return dto.TechnicCategory{
		ID:          c.ID,
		Name:        c.Name,
		Description: &c.Description,
		Audit: dto.AuditFields{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		},
	}
}
