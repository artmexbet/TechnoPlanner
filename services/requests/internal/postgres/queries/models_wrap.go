package queries

import (
	"requests/internal/domain"
	"utills/pointer"
)

func (r *Request) ToDomain() *domain.Request {
	var status domain.StatusType
	_ = status.Set(r.Status)
	return &domain.Request{
		ID:          r.ID,
		UserID:      r.TelegramUserID,
		RequestText: r.RequestText,
		Status:      status,
		Technics:    nil,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func RequestFromDomain(r domain.Request) *Request {
	return &Request{
		ID:             r.ID,
		TelegramUserID: r.UserID,
		RequestText:    r.RequestText,
		Status:         r.Status.String(),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func (u *TelegramUser) ToDomain() *domain.User {
	return &domain.User{
		ID:         u.ID,
		TelegramID: u.TelegramID,
		Username:   u.Username,
		FirstName:  u.FirstName,
		LastName:   pointer.From(u.LastName),
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
		LastName:   pointer.To(u.LastName),
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (t *Technic) ToDomain() *domain.Technic {
	return &domain.Technic{
		ID:          int(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Quantity:    int(t.Quantity),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func TechnicFromDomain(t domain.Technic) *Technic {
	return &Technic{
		ID:          int32(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Quantity:    int32(t.Quantity),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
