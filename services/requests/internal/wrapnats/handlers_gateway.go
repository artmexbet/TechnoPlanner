package wrapnats

import (
	"encoding/json"
	"log/slog"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

func (w *NatsWrapper) handleGatewayListRequests(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestListRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling gateway list request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	requests, err := w.reqService.ListByResponsible(ctx, req.ResponsibleID)
	if err != nil {
		slog.ErrorContext(ctx, "error listing requests", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make([]dto.Request, len(requests))
	for i, r := range requests {
		respDTO[i] = mapRequestToDTO(r)
	}

	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleGatewayGetRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestByIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling gateway get request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	reqDomain, err := w.reqService.Get(ctx, req.RequestID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting request", "error", err)
		return respondError(msg.Msg, "not found", "request not found", statusNotFound)
	}

	return respondSuccess(msg.Msg, "success", mapRequestToDTO(*reqDomain))
}

func (w *NatsWrapper) handleGatewayAssignResponsible(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.AssignResponsibleRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling gateway assign request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	updatedReq, err := w.reqService.AssignResponsible(ctx, req.RequestID, req.ResponsibleID)
	if err != nil {
		slog.ErrorContext(ctx, "error assigning responsible", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", mapRequestToDTO(*updatedReq))
}

func (w *NatsWrapper) handleGatewayUpdateRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RequestUpdateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling gateway update request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	updates := domain.RequestUpdate{
		RequestText:   req.RequestText,
		Status:        (*domain.StatusType)(req.Status),
		ScheduleTime:  req.ScheduleTime,
		Address:       req.Address,
		ResponsibleID: req.ResponsibleID,
	}

	updatedReq, err := w.reqService.UpdateRequest(ctx, req.RequestID, updates)
	if err != nil {
		slog.ErrorContext(ctx, "error updating request", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", mapRequestToDTO(*updatedReq))
}

func (w *NatsWrapper) handleGatewayListResponsibles(msg *broker.Msg) error {
	ctx := msg.Context()

	responsibles, err := w.reqService.ListResponsibles(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error listing responsibles", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make([]dto.Responsible, len(responsibles))
	for i, r := range responsibles {
		respDTO[i] = dto.Responsible{
			ID:       r.ID,
			Username: r.Username,
		}
	}

	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleGatewayCreateResponsible(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.ResponsibleCreateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling create responsible request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.reqService.SaveResponsible(ctx, req.ID, req.Username); err != nil {
		slog.ErrorContext(ctx, "error saving responsible", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", dto.Responsible{
		ID:       req.ID,
		Username: req.Username,
	})
}

func (w *NatsWrapper) handleGatewayGetResponsible(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.UUIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get responsible request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	responsible, err := w.reqService.GetResponsible(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting responsible", "error", err)
		return respondError(msg.Msg, "not found", "responsible not found", statusNotFound)
	}

	return respondSuccess(msg.Msg, "success", dto.Responsible{
		ID:       responsible.ID,
		Username: responsible.Username,
	})
}

func (w *NatsWrapper) handleGatewayDeleteResponsible(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.UUIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling delete responsible request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.reqService.DeleteResponsible(ctx, req.ID); err != nil {
		slog.ErrorContext(ctx, "error deleting responsible", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	return respondSuccess(msg.Msg, "success", nil)
}
