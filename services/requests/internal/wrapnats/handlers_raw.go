package wrapnats

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

func (w *NatsWrapper) handleGatewayListRawRequests(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RawRequestListRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling list raw requests", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	requests, err := w.reqService.ListRawRequests(ctx, req.Status, limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "error listing raw requests", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	respDTO := make([]dto.RawRequest, len(requests))
	for i, r := range requests {
		respDTO[i] = mapRawRequestToDTO(r)
	}

	return respondSuccess(msg.Msg, "success", respDTO)
}

func (w *NatsWrapper) handleGatewayGetRawRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.UUIDRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling get raw request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	rawReq, err := w.reqService.GetRawRequest(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error getting raw request", "error", err)
		return respondError(msg.Msg, "not found", "raw request not found", statusNotFound)
	}

	return respondSuccess(msg.Msg, "success", mapRawRequestToDTO(*rawReq))
}

func (w *NatsWrapper) handleGatewayProcessRawRequest(msg *broker.Msg) error {
	ctx := msg.Context()

	var req dto.RawRequestProcessRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.ErrorContext(ctx, "error unmarshaling process raw request", "error", err)
		return respondError(msg.Msg, "invalid request format", err.Error(), statusBadRequest)
	}

	if err := w.validator.Struct(req); err != nil {
		slog.ErrorContext(ctx, "validation error", "error", err)
		return respondError(msg.Msg, "validation error", err.Error(), statusBadRequest)
	}

	rawID, err := uuid.Parse(req.RawRequestID)
	if err != nil {
		return respondError(msg.Msg, "invalid raw_request_id", err.Error(), statusBadRequest)
	}

	equipments := make([]domain.Equipment, len(req.Equipments))
	for i, eq := range req.Equipments {
		equipments[i] = domain.Equipment{
			ID:       eq.ID,
			Quantity: eq.Quantity,
		}
	}

	newRequest := domain.Request{
		RequestText:     req.RequestText,
		ScheduleTime:    req.ScheduleTime,
		Address:         req.Address,
		Equipments:      equipments,
		EquipmentString: req.EquipmentString,
	}

	if req.EndTime != nil {
		newRequest.EndTime = *req.EndTime
	}

	createdReq, updatedRaw, err := w.reqService.ProcessRawRequest(ctx, rawID, newRequest)
	if err != nil {
		slog.ErrorContext(ctx, "error processing raw request", "error", err)
		return respondError(msg.Msg, "internal server error", err.Error(), statusInternalServerError)
	}

	type processResponse struct {
		Request    dto.Request    `json:"request"`
		RawRequest dto.RawRequest `json:"raw_request"`
	}

	return respondSuccess(msg.Msg, "success", processResponse{
		Request:    mapRequestToDTO(*createdReq),
		RawRequest: mapRawRequestToDTO(*updatedRaw),
	})
}
