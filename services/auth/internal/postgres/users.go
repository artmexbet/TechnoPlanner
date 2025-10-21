package postgres

import (
	"auth/internal/models"
	"auth/internal/postgres/queries"

	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *Postgres) FindUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	var pgID pgtype.UUID
	err := pgID.Scan(id)
	if err != nil {
		return models.User{}, fmt.Errorf("findUserByID: %w", err)
	}
	u, err := p.q.FindUserByID(ctx, pgID)
	if err != nil {
		return models.User{}, fmt.Errorf("findUserByID: %w", err)
	}
	return u.ToDomain(), nil
}

func (p *Postgres) CreateUser(ctx context.Context, username, email, passwordHash string) (models.User, error) {
	var pgEmail pgtype.Text
	err := pgEmail.Scan(email)
	if err != nil {
		return models.User{}, fmt.Errorf("createUser: %w", err)
	}
	u, err := p.q.CreateUser(ctx, queries.CreateUserParams{
		Username:     username,
		Email:        pgEmail,
		PasswordHash: passwordHash,
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
