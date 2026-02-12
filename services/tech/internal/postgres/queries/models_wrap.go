package queries

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"tech/internal/domain"
)

func DBTechSliceToDomain(dbTech []Technic) []domain.Technic {
	techToRet := make([]domain.Technic, 0, len(dbTech))
	for _, v := range dbTech {
		var chars AdditionalCharacteristics
		chars.UnmarshalJSON(v.AdditionalCharacteristics) //nolint:errcheck
		techToRet = append(techToRet, domain.Technic{
			ID:                        v.ID,
			CategoryID:                v.CategoryID.Bytes,
			Name:                      v.Name,
			Description:               *v.Description,
			AdditionalCharacteristics: chars.Chars,
			CreatedAt:                 v.CreatedAt,
			UpdatedAt:                 v.UpdatedAt,
		})
	}
	return techToRet
}

func DBCategorySliceToDomain(dbCategory []Category) []domain.TechnicCategory {
	domainCategories := make([]domain.TechnicCategory, 0, len(dbCategory))
	for _, v := range dbCategory {
		domainCategories = append(domainCategories, domain.TechnicCategory{Name: v.Name, ID: v.ID})
	}
	return domainCategories
}

func ConvertTechToQueryAdd(tech domain.Technic) (*AddTechnicParams, error) {
	chars := AdditionalCharacteristics{Chars: tech.AdditionalCharacteristics}
	marshalledChars, err := chars.MarshalJSON()
	if err != nil {
		slog.Error("Tech: DB: ConverTechToQueryAdd:", "error", err)
		return nil, err
	}
	var params AddTechnicParams
	params.CategoryID = pgtype.UUID{Bytes: tech.CategoryID}
	params.Name = tech.Name
	params.Description = &tech.Description
	params.AdditionalCharacteristics = marshalledChars
	return &params, nil
}

func ConvertTechToQueryUpdate(tech domain.Technic) (*UpdateTechnicParams, error) {
	chars := AdditionalCharacteristics{Chars: tech.AdditionalCharacteristics}
	marshalledChars, err := chars.MarshalJSON()
	if err != nil {
		slog.Error("Tech: DB: ConverTechToQueryUpdate:", "error", err)
		return nil, err
	}
	var params UpdateTechnicParams
	params.CategoryID = pgtype.UUID{Bytes: tech.CategoryID}
	params.Name = tech.Name
	params.Description = &tech.Description
	params.AdditionalCharacteristics = marshalledChars
	return &params, nil
}

func (q *Technic) ConvertToDomain() (*domain.Technic, error) {
	var chars AdditionalCharacteristics
	err := chars.UnmarshalJSON(q.AdditionalCharacteristics)
	if err != nil {
		slog.Error("Tech: DB: ConvertToDomain:", "error", err)
		return nil, err
	}
	tech := domain.Technic{
		ID:                        q.ID,
		CategoryID:                q.CategoryID.Bytes,
		Name:                      q.Name,
		Description:               *q.Description,
		AdditionalCharacteristics: chars.Chars,
	}
	return &tech, nil
}
