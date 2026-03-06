package queries

import (
	"encoding/json"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

func (t *RequestStatus) ToDomain() domain.StatusType {
	var status domain.StatusType
	_ = status.Set(string(*t))
	return status
}

func RequestStatusFromDomain(s domain.StatusType) *RequestStatus {
	str := RequestStatus(s.String())
	return &str
}

func (r *Request) ToDomain() *domain.Request {
	pInfo := &domain.PorterInfo{}
	if r.PorterID != nil {
		pInfo.UserID = *r.PorterID
	}
	req := &domain.Request{
		ID:           r.ID,
		RequestText:  r.RequestText,
		Status:       r.Status.ToDomain(),
		ScheduleTime: r.ScheduleTime,
		EndTime:      r.EndTime,
		Address:      r.Address,
		Equipments:   nil,
		PorterInfo:   pInfo,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}

	return req
}

func PorterInfoFromDomain(r *domain.PorterInfo) []byte {
	if r == nil {
		return nil
	}
	data, _ := json.Marshal(r)
	return data
}

func (u *TelegramUser) ToDomain() *domain.User {
	return &domain.User{
		ID:         u.ID,
		TelegramID: u.TelegramID,
		Username:   u.Username,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func TelegramUserFromDomain(u domain.User) *TelegramUser {
	return &TelegramUser{
		ID:         u.ID,
		TelegramID: u.TelegramID,
		Username:   u.Username,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (t *Equipment) ToDomain() *domain.Equipment {
	return &domain.Equipment{
		ID:               int(t.ID),
		Name:             t.Name,
		Description:      t.Description,
		Quantity:         int(t.Quantity),
		ReservedQuantity: int(t.ReservedQuantity),
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

func (t *GetEquipmentByRequestIDRow) ToDomain() *domain.Equipment {
	return &domain.Equipment{
		ID:          int(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Quantity:    int(t.Quantity),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func (t *BatchGetEquipmentByRequestIDRow) ToDomain() *domain.Equipment {
	return &domain.Equipment{
		ID:          int(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Quantity:    int(t.Quantity),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func EquipmentFromDomain(t domain.Equipment) *Equipment {
	return &Equipment{
		ID:               int32(t.ID),
		Name:             t.Name,
		Description:      t.Description,
		Quantity:         int32(t.Quantity),
		ReservedQuantity: int32(t.ReservedQuantity),
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

func (r *RawRequest) ToDomain() *domain.RawRequest {
	return &domain.RawRequest{
		ID:                 r.ID,
		TelegramID:         r.TelegramID,
		Username:           r.Username,
		FirstName:          r.FirstName,
		LastName:           r.LastName,
		RawText:            r.RawText,
		Status:             domain.RawRequestStatus(r.Status),
		ProcessedRequestID: r.ProcessedRequestID,
		CreatedAt:          r.CreatedAt,
	}
}
