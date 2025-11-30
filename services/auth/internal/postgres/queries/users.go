package queries

import (
	"log/slog"

	"github.com/google/uuid"

	"auth/internal/models"

	"utills/pointer"
)

func (u *User) ToDomain() models.User {
	id, err := uuid.Parse(u.ID.String())
	if err != nil {
		slog.Error("User ToDomain: parse id", "id", u.ID, "err", err)
	}
	return models.User{
		ID:           id,
		Username:     u.Username,
		Email:        pointer.From(u.Email),
		PasswordHash: u.PasswordHash,
		RoleID:       u.RoleID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
