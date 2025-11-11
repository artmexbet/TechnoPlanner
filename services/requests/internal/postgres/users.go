package postgres

import (
	"context"

	"requests/internal/domain"
	"requests/internal/postgres/queries"
	"utills/pointer"

	"github.com/google/uuid"
)

func (p *Postgres) GetUserByTelegramID(ctx context.Context, telegramID int64) (domain.User, error) {
	u, err := p.q.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return domain.User{}, err
	}
	return *u.ToDomain(), nil
}

func (p *Postgres) SaveTelegramUser(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := p.q.SaveTelegramUser(ctx, queries.SaveTelegramUserParams{
		TelegramID: user.TelegramID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   pointer.To(user.LastName),
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
