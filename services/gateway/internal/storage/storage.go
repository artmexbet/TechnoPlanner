package storage

import (
	"context"

	"gateway/internal/domain"
)

type postgres interface {
	AddPlace(context.Context, domain.Place) (domain.Place, error)
}

type Storage struct {
	p postgres
}

func NewStorage(p postgres) *Storage {
	return &Storage{
		p: p,
	}
}
