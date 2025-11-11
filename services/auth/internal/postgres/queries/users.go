package queries

import (
	"auth/internal/models"
	"utills/pointer"

	"log/slog"

	"github.com/google/uuid"
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
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
