package queries

import "gateway/internal/domain"

func (p Place) ToDomain() domain.Place {
	return domain.Place{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
