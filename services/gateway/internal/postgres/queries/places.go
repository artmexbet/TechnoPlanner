package queries

import (
	"log/slog"

	"github.com/google/uuid"

	"gateway/internal/domain"
)

func (p *Place) ToDomain() domain.Place {
	var v uuid.UUID
	err := p.ID.Scan(&v)
	if err != nil {
		slog.Error("failed to scan Place ID to uuid.UUID in ToDomain", "place_id", p.ID, "place_name", p.Name, "err", err)
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
