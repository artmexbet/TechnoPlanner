package wrapnats

import (
	"encoding/json"
	"log/slog"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

// handleCreateRequest обрабатывает запрос от бота — сохраняет его как сырой запрос (raw request).
// Старый контракт dto.RequestCreateRequest сохраняется без изменений.
func (w *NatsWrapper) handleCreateRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestCreateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	// Формируем сырой текст из полей запроса
	rawText := ""
	if req.Text != nil {
		rawText = *req.Text
	}
	if req.EquipmentString != nil && *req.EquipmentString != "" {
		if rawText != "" {
			rawText += "\n"
		}
		rawText += *req.EquipmentString
	}

	username := ""
	if req.Username != nil {
		username = *req.Username
	}

	domainRaw := domain.RawRequest{
		TelegramID: req.TelegramUserID,
		Username:   username,
		FirstName:  username, // бот передаёт только username в этом контракте
		RawText:    rawText,
	}

	created, err := w.reqService.CreateRawRequest(ctx, domainRaw)
	if err != nil {
		slog.ErrorContext(ctx, "error creating raw request", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", mapRawRequestToDTO(*created))
}

func (w *NatsWrapper) handleUpdateStatus(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestStatusUpdateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling update status request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	if err := w.reqService.UpdateStatus(ctx, req.RequestID, domain.StatusType(req.Status)); err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", "status updated")
}

func (w *NatsWrapper) handleGetRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	reqDomain, err := w.reqService.Get(ctx, req.RequestID)
	if err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", mapRequestToDTO(*reqDomain))
}

func (w *NatsWrapper) handleListRequests(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestListByTelegramRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling list request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	user := domain.User{TelegramID: req.TelegramID}
	requests, err := w.reqService.List(ctx, user, req.Limit, req.Offset)
	if err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make([]dto.Request, len(requests))
	for i, r := range requests {
		respDTO[i] = mapRequestToDTO(r)
	}

	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleCancelRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling cancel request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.reqService.Cancel(ctx, req.RequestID); err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", "request canceled")
}

func (w *NatsWrapper) handleCancelRequestFromService(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling cancel request from service", "error", err)
		return err
	}

	if err := w.reqService.Cancel(ctx, req.RequestID); err != nil {
		slog.ErrorContext(ctx, "error canceling request from service", "error", err)
		return err
	}

	slog.InfoContext(ctx, "request canceled from service", "request_id", req.RequestID)
	return nil
}

func (w *NatsWrapper) handleAddEquipment(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.EquipmentAddRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling add equipment request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	domainEquipment := make([]domain.Equipment, len(req.Equipments))
	for i, eq := range req.Equipments {
		domainEquipment[i] = domain.Equipment{
			ID:       eq.ID,
			Name:     eq.Name,
			Quantity: eq.Quantity,
		}
	}

	if err := w.eqService.Add(ctx, domainEquipment); err != nil {
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", "equipment added")
}
