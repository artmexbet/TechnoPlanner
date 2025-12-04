package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/postgres/queries"
)

func (p *Postgres) GetUserByTelegramID(ctx context.Context, telegramID int64) (domain.User, error) {
	u, err := p.q.GetUserByTelegramID(ctx, telegramID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	return *u.ToDomain(), nil
}

func (p *Postgres) SaveTelegramUser(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := p.q.SaveTelegramUser(ctx, queries.SaveTelegramUserParams{
		TelegramID: user.TelegramID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	})
	if err != nil {
		return domain.User{}, err
	}
	return *u.ToDomain(), nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	u, err := p.q.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return *u.ToDomain(), nil
}
