package wrapnats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

type statusCode int

const (
	statusOK                  statusCode = 200
	statusBadRequest          statusCode = 400
	statusUnauthorized        statusCode = 401
	statusForbidden           statusCode = 403
	statusNotFound            statusCode = 404
	statusInternalServerError statusCode = 500
)

func respondError(msg *nats.Msg, message string, payload interface{}, statusCode statusCode) error { //nolint:revive //todo: удалить статус код, он не нужен, так как в GatewayResponse уже есть поле Message для описания ошибки
	resp := dto.GatewayResponse{
		Success: false,
		Message: message,
	}

	// Если payload - это string, просто используем его как сообщение
	if _, ok := payload.(string); ok {
		resp.Data = json.RawMessage(fmt.Sprintf(`"%s"`, payload))
	} else {
		data, _ := json.Marshal(payload)
		resp.Data = data
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshaling error response: %w", err)
	}
	err = msg.Respond(data)
	if err != nil {
		return fmt.Errorf("respondError: %w", err)
	}
	slog.Info("responded with error", "message", resp.Message)
	return nil
}

func respondSuccess(msg *nats.Msg, message string, payload interface{}) error {
	resp := dto.GatewayResponse{
		Success: true,
		Message: message,
	}

	// Если payload уже []byte, используем его напрямую
	if data, ok := payload.([]byte); ok {
		resp.Data = data
	} else {
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshaling success response payload: %w", err)
		}
		resp.Data = data
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshaling success response: %w", err)
	}
	err = msg.Respond(data)
	if err != nil {
		return fmt.Errorf("respondSuccess: %w", err)
	}
	slog.Info("responded with success", "message", resp.Message)
	return nil
}

// Helper functions for mapping

func mapRequestCreateToDomain(req dto.RequestCreateRequest) domain.Request {
	equipments := make([]domain.Equipment, len(req.Equipments))
	for i, eq := range req.Equipments {
		equipments[i] = domain.Equipment{
			ID:       eq.ID,
			Quantity: eq.Quantity,
		}
	}

	var username string
	if req.Username != nil {
		username = *req.Username
	}

	return domain.Request{
		RequestText:     req.Text,
		ScheduleTime:    req.ScheduleTime,
		Equipments:      equipments,
		EquipmentString: req.EquipmentString,
		Address:         req.Address,
		Issuer: domain.User{
			TelegramID: req.TelegramUserID,
			Username:   username,
		},
	}
}

func mapRequestToDTO(req domain.Request) dto.Request {
	equipment := make([]dto.RequestEquipment, len(req.Equipments))
	for i, eq := range req.Equipments {
		equipment[i] = dto.RequestEquipment{
			RequestID:   req.ID,
			EquipmentID: int32(eq.ID),
			Quantity:    int32(eq.Quantity),
			CreatedAt:   req.CreatedAt,
			UpdatedAt:   req.UpdatedAt,
		}
	}

	return dto.Request{
		ID:           req.ID,
		RequestText:  req.RequestText,
		Status:       dto.RequestStatus(req.Status),
		ScheduleTime: req.ScheduleTime,
		EndTime:      req.EndTime,
		Address:      req.Address,
		Equipment:    equipment,
		Audit: dto.AuditFields{
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
		},
	}
}
