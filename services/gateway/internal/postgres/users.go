package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/postgres/queries"
)

func (d *DB) ListPorters(ctx context.Context, roleID int32) ([]domain.User, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	rows, err := d.q.ListPorters(ctx, roleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("ListPorters: %w", err)
	}

	res := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		res = append(res, mapUser(row))
	}

	return res, nil
}

func (d *DB) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	row, err := d.q.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("GetUserByID: %w", err)
	}

	return mapUser(row), nil
}

func (d *DB) CreateUser(ctx context.Context, user domain.User) error {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	_, err := d.q.CreateUser(ctx, queries.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
	})
	if err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}

	return nil
}

func mapUser(row queries.User) domain.User {
	return domain.User{
		ID:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		RoleID:    row.RoleID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}
}
