package domain

import (
	"context"

	"github.com/google/uuid"
)

type DistrictRepository interface {
	GetAll(ctx context.Context, cityID uuid.UUID) (districtList []*District, err error)
	FindByID(ctx context.Context, id uuid.UUID) (district *District, err error)
	Create(ctx context.Context, district *District) (err error)
	Update(ctx context.Context, district *District) (err error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByFilter(ctx context.Context, filter *DistrictListFilter) (districtList []*District, err error)
}

type District struct {
	ID         uuid.UUID `json:"id"`
	CountryID  uuid.UUID `json:"country_id"`
	ProvinceID uuid.UUID `json:"province_id"`
	CityID     uuid.UUID `json:"city_id"`
	Name       string    `json:"name"`
	Sort       int       `json:"sort"`
}

type DistrictListFilter struct {
	CountryID  uuid.UUID
	ProvinceID uuid.UUID
	CityID     uuid.UUID
	Name       string
}
