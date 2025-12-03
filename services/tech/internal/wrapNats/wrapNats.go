package wrapnats

import (
	"context"
	"log/slog"
	"tech/internal/domain"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type ITechService interface {
	AddTechnic(ctx context.Context, technic domain.Technic) (domain.Technic, error)
	DeleteTechnic(ctx context.Context, techID uuid.UUID) error
	UpdateTechnic(ctx context.Context, technic domain.Technic) (domain.Technic, error)
	GetTechnicByID(ctx context.Context, techID uuid.UUID) (domain.Technic, error)
	GetTechnicByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Technic, error)
	AddCategory(ctx context.Context, categoryName string) (domain.TechnicCategory, error)
	UpdateCategoryName(ctx context.Context, category domain.TechnicCategory) (domain.TechnicCategory, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
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

	techService ITechService
}

func New(cfg *Config, conn *nats.Conn, techService ITechService) (*NatsWrapper, error) {
	return &NatsWrapper{
		conn:        conn,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
		cfg:         cfg,
		techService: techService,
	}, nil
}

func (w *NatsWrapper) HandleMsgs() *NatsWrapper {
	handlers := map[string]nats.MsgHandler{
		"tech.add":         w.handleAddTechnic,
		"tech.update":      w.handleUpdateTechnic,
		"tech.delete":      w.handleDeleteTechnic,
		"tech.get":         w.handleGetTechnic,
		"category.get":     w.handleCategoryGetTech,
		"category.add":     w.handleAddCategory,
		"category.update":  w.handleUpdateCategory,
		"category.delete":  w.handleDeleteCategory,
		"category.get.all": w.handleGetAllCategories,
	}

	for subject, handler := range handlers {
		_, err := w.conn.Subscribe(subject, handler)
		if err != nil {
			slog.Error("cannot subscribe to subject", "subject", subject, "error", err)
		}
	}
	return w
}

func (w *NatsWrapper) handleAddTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestUpdateTech
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	domainReq := req.ToDomain()
	createdReq, err := w.techService.AddTechnic(ctx, domainReq)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respData, err := createdReq.MarshalJSON()
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleUpdateTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestUpdateTech
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	updatedReq, err := w.techService.UpdateTechnic(ctx, req.ToDomain())
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respData, err := updatedReq.MarshalJSON()
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestSpecificTech
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	getReq, err := w.techService.GetTechnicByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respData, err := getReq.MarshalJSON()
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleDeleteTechnic(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestSpecificTech
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err = w.techService.DeleteTechnic(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", "tech deleted")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleCategoryGetTech(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestSpecificCategory
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	respData, err := w.techService.GetTechnicByCategory(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleAddCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestUpdateCategory
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	createdReq, err := w.techService.AddCategory(ctx, req.Name)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respData, err := createdReq.MarshalJSON()
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleDeleteCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestSpecificCategory
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	err = w.techService.DeleteCategory(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", "category deleted successfully")
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleUpdateCategory(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	var req requestUpdateCategory
	err := req.UnmarshalJSON(msg.Data)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		_ = respondError(msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err = w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		_ = respondError(msg, "validation error", err.Error(), statusBadRequest)
		return
	}

	createdReq, err := w.techService.UpdateCategoryName(ctx, domain.TechnicCategory{ID: req.ID, Name: req.Name})
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	respData, err := createdReq.MarshalJSON()
	if err != nil {
		slog.ErrorContext(ctx, "error marshaling response", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}

func (w *NatsWrapper) handleGetAllCategories(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), w.cfg.RequestTimeout)
	defer cancel()

	respData, err := w.techService.GetAllCategories(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", "error", err)
		_ = respondError(msg, "internal server error", err.Error(), statusInternalServerError)
		return
	}

	err = respondSuccess(msg, "success", respData)
	if err != nil {
		slog.ErrorContext(ctx, "error sending response", "error", err)
	}
}
