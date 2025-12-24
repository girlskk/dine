package domain

import (
	"context"

	"github.com/google/uuid"
)

type CountryRepository interface {
	GetAll(ctx context.Context) (countryList []*Country, err error)
	FindByID(ctx context.Context, id uuid.UUID) (country *Country, err error)
	Create(ctx context.Context, country *Country) (err error)
	Update(ctx context.Context, country *Country) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
}

type Country struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Sort int       `json:"sort"`
}
