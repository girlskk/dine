package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/city.go -package=mock . CityRepository
type CityRepository interface {
	GetAll(ctx context.Context, provinceID uuid.UUID) (cityList []*City, err error)
	FindByID(ctx context.Context, id uuid.UUID) (city *City, err error)
	Create(ctx context.Context, city *City) (err error)
	Update(ctx context.Context, city *City) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetByFilter(ctx context.Context, filter *CityListFilter) (cityList []*City, err error)
}

type City struct {
	ID         uuid.UUID `json:"id"`
	CountryID  uuid.UUID `json:"country_id"`
	ProvinceID uuid.UUID `json:"province_id"`
	Name       string    `json:"name"`
	Sort       int       `json:"sort"`
}

type CityListFilter struct {
	CountryID  uuid.UUID
	ProvinceID uuid.UUID
	Name       string
}
