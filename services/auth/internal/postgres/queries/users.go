package queries

import (
	"auth/internal/models"

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
		Email:        u.Email.String,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}
