package queries

import (
	"log/slog"

	"technoBro/internal/domain"

	"github.com/google/uuid"
)

func (p *Place) ToDomain() domain.Place {
	var v uuid.UUID
	err := p.ID.Scan(&v)
	if err != nil {
		slog.Error("cannot convert id to domain", "err", err)
		return domain.Place{}
	}
	return domain.Place{
		ID:          v,
		Name:        p.Name,
		Description: p.Description.String,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}
