package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/province.go -package=mock . ProvinceRepository
type ProvinceRepository interface {
	GetAll(ctx context.Context, countryID uuid.UUID) (provinceList []*Province, err error)
	FindByID(ctx context.Context, id uuid.UUID) (province *Province, err error)
	Create(ctx context.Context, province *Province) (err error)
	Update(ctx context.Context, province *Province) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetByFilter(ctx context.Context, filter *ProvinceListFilter) (provinceList []*Province, err error)
}

type Province struct {
	ID        uuid.UUID `json:"id"`
	CountryID uuid.UUID `json:"country_id"`
	Name      string    `json:"name"`
	Sort      int       `json:"sort"`
}

type ProvinceListFilter struct {
	CountryID uuid.UUID
	Name      string
}
