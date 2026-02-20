package wrapnats

import (
	"github.com/artmexbet/TechnoPlanner/libs/dto"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

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

func mapRawRequestToDTO(r domain.RawRequest) dto.RawRequest {
	result := dto.RawRequest{
		ID:         r.ID.String(),
		TelegramID: r.TelegramID,
		Username:   r.Username,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		RawText:    r.RawText,
		Status:     string(r.Status),
		CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if r.ProcessedRequestID != nil {
		id := r.ProcessedRequestID.String()
		result.ProcessedRequestID = &id
	}
	return result
}
