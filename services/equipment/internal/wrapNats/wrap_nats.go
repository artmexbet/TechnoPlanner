package wrapnats

import (
	"context"
	"log/slog"
	"time"

	"tech/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/mailru/easyjson"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
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
	ReserveEquipment(ctx context.Context, items []domain.ReserveItem) error
	ReleaseEquipment(ctx context.Context, items []domain.ReserveItem) error
	CheckAvailability(ctx context.Context, items []domain.ReserveItem) (bool, []int, error)
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
	conn      *broker.NATS
	cfg       *Config
	validator *validator.Validate
	publisher *NatsPublisher

	equipmentService EquipmentService
	categoryService  CategoryService
}

func New(cfg *Config, conn *broker.NATS, equipmentService EquipmentService, categoryService CategoryService) (*NatsWrapper, error) {
	return &NatsWrapper{
		conn:             conn,
		validator:        validator.New(validator.WithRequiredStructEnabled()),
		cfg:              cfg,
		publisher:        NewNatsPublisher(conn),
		equipmentService: equipmentService,
		categoryService:  categoryService,
	}, nil
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]broker.MsgHandler{
		// Equipment handlers
		subjects.GatewayEquipmentCreate:        w.handleCreateEquipment,
		subjects.GatewayEquipmentUpdate:        w.handleUpdateEquipment,
		subjects.GatewayEquipmentDelete:        w.handleDeleteEquipment,
		subjects.GatewayEquipmentGet:           w.handleGetEquipment,
		subjects.GatewayEquipmentList:          w.handleListEquipment,
		subjects.GatewayEquipmentGetByCategory: w.handleGetEquipmentByCategory,
		// Reservation handlers
		subjects.GatewayEquipmentReserve:           w.handleReserveEquipment,
		subjects.GatewayEquipmentRelease:           w.handleReleaseEquipment,
		subjects.GatewayEquipmentCheckAvailability: w.handleCheckAvailability,
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

// ─── Equipment handlers ──────────────────────────────────────────────────────

func (w *NatsWrapper) handleCreateEquipment(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentCreateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create equipment request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	domainReq := domain.Equipment{
		CategoryID:                req.CategoryID,
		Name:                      req.Name,
		Description:               req.Description,
		AdditionalCharacteristics: req.AdditionalCharacteristics,
		Quantity:                  req.Quantity,
	}

	createdReq, err := w.equipmentService.AddEquipment(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error creating equipment", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapEquipmentToDTO(*createdReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}

	w.publisher.PublishEquipmentCreated(dto.EquipmentSyncEvent{
		ID:          createdReq.ID,
		Name:        createdReq.Name,
		Description: createdReq.Description,
		Quantity:    createdReq.Quantity,
	})
	return nil
}

func (w *NatsWrapper) handleUpdateEquipment(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentUpdateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update equipment request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	domainReq := domain.Equipment{
		ID:                        req.ID,
		CategoryID:                req.CategoryID,
		Name:                      req.Name,
		Description:               req.Description,
		AdditionalCharacteristics: req.AdditionalCharacteristics,
		Quantity:                  req.Quantity,
	}

	updatedReq, err := w.equipmentService.UpdateEquipment(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error updating equipment", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapEquipmentToDTO(*updatedReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}

	w.publisher.PublishEquipmentUpdated(dto.EquipmentSyncEvent{
		ID:          updatedReq.ID,
		Name:        updatedReq.Name,
		Description: updatedReq.Description,
		Quantity:    updatedReq.Quantity,
	})
	return nil
}

func (w *NatsWrapper) handleGetEquipment(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentGetByIDRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get equipment request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	getReq, err := w.equipmentService.GetEquipmentByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting equipment", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapEquipmentToDTO(*getReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleDeleteEquipment(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentDeleteRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling delete equipment request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	if err := w.equipmentService.DeleteEquipment(ctx, req.ID); err != nil {
		slog.ErrorContext(ctx, "error deleting equipment", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", "equipment deleted"); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}

	w.publisher.PublishEquipmentDeleted(req.ID)
	return nil
}

func (w *NatsWrapper) handleListEquipment(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	equipment, err := w.equipmentService.GetAllEquipment(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error listing equipment", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make(dto.TechEquipmentList, len(equipment))
	for i, e := range equipment {
		respDTO[i] = mapEquipmentToDTO(e)
	}

	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleGetEquipmentByCategory(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentGetByCategoryRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get equipment by category request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	equipment, err := w.equipmentService.GetEquipmentByCategory(ctx, req.CategoryID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting equipment by category", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make(dto.TechEquipmentList, len(equipment))
	for i, e := range equipment {
		respDTO[i] = mapEquipmentToDTO(e)
	}

	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

// ─── Category handlers ───────────────────────────────────────────────────────

func (w *NatsWrapper) handleCreateCategory(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentCategoryCreateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create category request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	createdReq, err := w.categoryService.AddCategory(ctx, req.Name, req.Description)
	if err != nil {
		slog.ErrorContext(ctx, "error creating category", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapCategoryToDTO(*createdReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleUpdateCategory(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentCategoryUpdateRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update category request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	updatedReq, err := w.categoryService.UpdateCategory(ctx, domain.EquipmentCategory{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		slog.ErrorContext(ctx, "error updating category", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapCategoryToDTO(*updatedReq)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleDeleteCategory(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentCategoryDeleteRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling delete category request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	if err := w.categoryService.DeleteCategory(ctx, req.ID); err != nil {
		slog.ErrorContext(ctx, "error deleting category", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", "category deleted successfully"); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleGetCategory(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentCategoryGetByIDRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get category request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	category, err := w.categoryService.GetCategoryByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting category", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := mapCategoryToDTO(*category)
	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleListCategories(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	categories, err := w.categoryService.GetAllCategories(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error listing categories", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make(dto.TechEquipmentCategoryList, len(categories))
	for i, c := range categories {
		respDTO[i] = mapCategoryToDTO(c)
	}

	respData, err := easyjson.Marshal(respDTO)
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "success", respData); err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
	return nil
}

// ─── Helper functions ────────────────────────────────────────────────────────

func mapEquipmentToDTO(e domain.Equipment) dto.TechEquipment {
	return dto.TechEquipment{
		ID:                        e.ID,
		CategoryID:                e.CategoryID,
		Name:                      e.Name,
		Description:               e.Description,
		AdditionalCharacteristics: e.AdditionalCharacteristics,
		Quantity:                  e.Quantity,
		ReservedQuantity:          e.ReservedQuantity,
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

// ─── Reservation handlers ────────────────────────────────────────────────────

func (w *NatsWrapper) handleReserveEquipment(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentReserveRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "handleReserveEquipment: unmarshal error", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	items := make([]domain.ReserveItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = domain.ReserveItem{EquipmentID: it.EquipmentID, Quantity: it.Quantity}
	}

	if err := w.equipmentService.ReserveEquipment(ctx, items); err != nil {
		slog.ErrorContext(ctx, "handleReserveEquipment: service error", "error", err)
		return respondError(msg.Msg, "reserve failed", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "reserved", nil); err != nil {
		slog.ErrorContext(ctx, "handleReserveEquipment: respond error", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleReleaseEquipment(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentReleaseRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "handleReleaseEquipment: unmarshal error", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	items := make([]domain.ReserveItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = domain.ReserveItem{EquipmentID: it.EquipmentID, Quantity: it.Quantity}
	}

	if err := w.equipmentService.ReleaseEquipment(ctx, items); err != nil {
		slog.ErrorContext(ctx, "handleReleaseEquipment: service error", "error", err)
		return respondError(msg.Msg, "release failed", err.Error(), statusInternalServerError)
	}

	if err := respondSuccess(msg.Msg, "released", nil); err != nil {
		slog.ErrorContext(ctx, "handleReleaseEquipment: respond error", "error", err)
	}
	return nil
}

func (w *NatsWrapper) handleCheckAvailability(msg *broker.Msg) error {
	ctx, cancel := context.WithTimeout(msg.Context(), w.cfg.RequestTimeout)
	defer cancel()

	var req dto.TechEquipmentCheckRequest
	if err := easyjson.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "handleCheckAvailability: unmarshal error", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	items := make([]domain.ReserveItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = domain.ReserveItem{EquipmentID: it.EquipmentID, Quantity: it.Quantity}
	}

	ok, unavailable, err := w.equipmentService.CheckAvailability(ctx, items)
	if err != nil {
		slog.ErrorContext(ctx, "handleCheckAvailability: service error", "error", err)
		return respondError(msg.Msg, "check failed", err.Error(), statusInternalServerError)
	}

	resp := dto.TechEquipmentCheckResponse{
		Available:      ok,
		UnavailableIDs: unavailable,
	}
	respData, _ := easyjson.Marshal(resp)
	if err := respondSuccess(msg.Msg, "checked", respData); err != nil {
		slog.ErrorContext(ctx, "handleCheckAvailability: respond error", "error", err)
	}
	return nil
}
