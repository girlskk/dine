package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/country.go -package=mock . CountryRepository
type CountryRepository interface {
	GetAll(ctx context.Context) (countryList []*Country, err error)
	FindByID(ctx context.Context, id uuid.UUID) (country *Country, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/country_interactor.go -package=mock . CountryInteractor
type CountryInteractor interface {
	GetAllCountries(ctx context.Context) (countryList []*Country, err error)
}

type Country struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Sort int       `json:"sort"`
}
