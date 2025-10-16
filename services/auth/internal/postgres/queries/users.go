package queries

import (
	"auth/internal/models"

	"github.com/google/uuid"
)

func (u *User) ToDomain() models.User {
	var id uuid.UUID
	_ = id.Scan(u.ID)
	return models.User{
		ID:           id,
		Username:     u.Username,
		Email:        u.Email.String,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}
