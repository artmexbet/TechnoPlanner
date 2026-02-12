package wrapnats

import (
	"context"
	"log/slog"
	"time"

	"tech/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/mailru/easyjson"
	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/config/subjects"
	"github.com/artmexbet/TechnoPlanner/libs/dto"
)

type EquipmentService interface {
	AddEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error)
	DeleteEquipment(ctx context.Context, equipmentID int) error
	UpdateEquipment(ctx context.Context, equipment domain.Equipment) (*domain.Equipment, error)
	GetEquipmentByID(ctx context.Context, equipmentID int) (*domain.Equipment, error)
	GetAllEquipment(ctx context.Context) ([]domain.Equipment, error)
	GetEquipmentByCategory(ctx context.Context, categoryID int) ([]domain.Equipment, error)
}

type CategoryService interface {
	AddCategory(ctx context.Context, categoryName string, description string) (*domain.EquipmentCategory, error)
	UpdateCategory(ctx context.Context, category domain.EquipmentCategory) (*domain.EquipmentCategory, error)
	DeleteCategory(ctx context.Context, categoryID int) error
	GetCategoryByID(ctx context.Context, categoryID int) (*domain.EquipmentCategory, error)
	GetAllCategories(ctx context.Context) ([]domain.EquipmentCategory, error)
}

type Config struct {
	URL string `yaml:"url" env:"URL"`

	RequestTimeout time.Duration `yaml:"request_timeout" env:"REQUEST_TIMEOUT" env-default:"5s"`
}

type NatsWrapper struct {
	conn      *nats.Conn
	cfg       *Config
	validator *validator.Validate

	equipmentService EquipmentService
	categoryService  CategoryService
}

func New(cfg *Config, conn *nats.Conn, equipmentService EquipmentService, categoryService CategoryService) (*NatsWrapper, error) {
	return &NatsWrapper{
		conn:             conn,
		validator:        validator.New(validator.WithRequiredStructEnabled()),
		cfg:              cfg,
		equipmentService: equipmentService,
		categoryService:  categoryService,
	}, nil
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]nats.MsgHandler{
		// Equipment handlers
		subjects.GatewayEquipmentCreate:        w.handleCreateEquipment,
		subjects.GatewayEquipmentUpdate:        w.handleUpdateEquipment,
		subjects.GatewayEquipmentDelete:        w.handleDeleteEquipment,
		subjects.GatewayEquipmentGet:           w.handleGetEquipment,
		subjects.GatewayEquipmentList:          w.handleListEquipment,
		subjects.GatewayEquipmentGetByCategory: w.handleGetEquipmentByCategory,
		// Category handlers
		subjects.GatewayEquipmentCategoryCreate: w.handleCreateCategory,
		subjects.GatewayEquipmentCategoryUpdate: w.handleUpdateCategory,
		subjects.GatewayEquipmentCategoryDelete: w.handleDeleteCategory,
		subjects.GatewayEquipmentCategoryGet:    w.handleGetCategory,
		subjects.GatewayEquipmentCategoryList:   w.handleListCategories,
	}

	for subject, handler := range handlers {
		_, err := w.conn.Subscribe(subject, handler)
		if err != nil {
			slog.Error("cannot subscribe to subject", "subject", subject, "error", err)
		}
	}
	return w
}

// Equipment handlers

func (w *NatsWrapper) handleCreateEquipment(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentCreateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create equipment request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	domainReq := domain.Equipment{
		CategoryID:                req.CategoryID,
		Name:                      req.Name,
		Description:               req.Description,
		AdditionalCharacteristics: req.AdditionalCharacteristics,
	}

	createdReq, err := w.equipmentService.AddEquipment(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error creating equipment", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapEquipmentToDTO(*createdReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleUpdateEquipment(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentUpdateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update equipment request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	domainReq := domain.Equipment{
		ID:                        req.ID,
		CategoryID:                req.CategoryID,
		Name:                      req.Name,
		Description:               req.Description,
		AdditionalCharacteristics: req.AdditionalCharacteristics,
	}

	updatedReq, err := w.equipmentService.UpdateEquipment(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error updating equipment", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapEquipmentToDTO(*updatedReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetEquipment(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentGetByIDRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get equipment request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	getReq, err := w.equipmentService.GetEquipmentByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting equipment", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := mapEquipmentToDTO(*getReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleDeleteEquipment(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentDeleteRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling delete equipment request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err := w.equipmentService.DeleteEquipment(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error deleting equipment", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", "equipment deleted"); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleListEquipment(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	equipment, err := w.equipmentService.GetAllEquipment(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error listing equipment", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := make(dto.TechEquipmentList, len(equipment))
	for i, e := range equipment {
		respDTO[i] = mapEquipmentToDTO(e)
	}

	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	if err := respondSuccess(msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetEquipmentByCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentGetByCategoryRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get equipment by category request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	equipment, err := w.equipmentService.GetEquipmentByCategory(ctx, req.CategoryID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting equipment by category", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respDTO := make(dto.TechEquipmentList, len(equipment))
	for i, e := range equipment {
		respDTO[i] = mapEquipmentToDTO(e)
	}

	respData, err := easyjson.Marshal(respDTO)
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

	var req dto.TechEquipmentCategoryCreateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
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
	respData, err := easyjson.Marshal(respDTO)
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

	var req dto.TechEquipmentCategoryUpdateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update category request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
		return
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	updatedReq, err := w.categoryService.UpdateCategory(ctx, domain.EquipmentCategory{
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
	respData, err := easyjson.Marshal(respDTO)
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

	var req dto.TechEquipmentCategoryDeleteRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
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

	var req dto.TechEquipmentCategoryGetByIDRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
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
	respData, err := easyjson.Marshal(respDTO)
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

	respDTO := make(dto.TechEquipmentCategoryList, len(categories))
	for i, c := range categories {
		respDTO[i] = mapCategoryToDTO(c)
	}

	respData, err := easyjson.Marshal(respDTO)
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

func mapEquipmentToDTO(e domain.Equipment) dto.TechEquipment {
	return dto.TechEquipment{
		ID:                        e.ID,
		CategoryID:                e.CategoryID,
		Name:                      e.Name,
		Description:               e.Description,
		AdditionalCharacteristics: e.AdditionalCharacteristics,
		CreatedAt:                 e.CreatedAt,
		UpdatedAt:                 e.UpdatedAt,
	}
}

func mapCategoryToDTO(c domain.EquipmentCategory) dto.TechEquipmentCategory {
	return dto.TechEquipmentCategory{
		ID:          c.ID,
		Name:        c.Name,
		Description: &c.Description,
		Audit: dto.AuditFields{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		},
	}
}
