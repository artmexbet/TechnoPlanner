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
	req := &domain.Request{
		ID:          r.ID,
		RequestText: r.RequestText,
		Status:      r.Status.ToDomain(),
		Equipments:  nil,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	// Parse ResponsibleInfo from JSONB
	if len(r.ResponsibleInfo) > 0 {
		var respInfo domain.ResponsibleInfo
		if err := json.Unmarshal(r.ResponsibleInfo, &respInfo); err == nil {
			req.ResponsibleInfo = &respInfo
		}
	}

	return req
}

func RequestFromDomain(r domain.Request) *Request {
	req := &Request{
		ID:             r.ID,
		TelegramUserID: r.Issuer.ID,
		RequestText:    r.RequestText,
		Status:         *RequestStatusFromDomain(r.Status),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}

	// Convert ResponsibleInfo to JSONB
	if r.ResponsibleInfo != nil {
		data, _ := json.Marshal(r.ResponsibleInfo)
		req.ResponsibleInfo = data
	}

	return req
}

func ResponsibleInfoFromDomain(r *domain.ResponsibleInfo) []byte {
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
		ID:          int32(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Quantity:    int32(t.Quantity),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
