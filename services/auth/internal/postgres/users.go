package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/libs/utills/pointer"

	"github.com/artmexbet/TechnoPlanner/services/auth/internal/models"
	"github.com/artmexbet/TechnoPlanner/services/auth/internal/postgres/queries"
)

func (p *Postgres) FindUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	u, err := p.q.FindUserByID(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("findUserByID: %w", err)
	}
	return u.ToDomain(), nil
}

func (p *Postgres) CreateUser(ctx context.Context, username, email, passwordHash string, roleID int32) (models.User, error) {
	u, err := p.q.CreateUser(ctx, queries.CreateUserParams{
		Username:     username,
		Email:        pointer.To(email),
		PasswordHash: passwordHash,
		RoleID:       roleID,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("createUser: %w", err)
	}
	return u.ToDomain(), nil
}

func (p *Postgres) FindUserByUsername(ctx context.Context, username string) (models.User, error) {
	u, err := p.q.FindUserByUsername(ctx, username)
	if err != nil {
		return models.User{}, fmt.Errorf("findUserByUsername: %w", err)
	}
	return u.ToDomain(), nil
}
