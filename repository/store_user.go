package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
)

var _ domain.StoreUserRepository = (*StoreUserRepository)(nil)

type StoreUserRepository struct {
	Client *ent.Client
}

func NewStoreUserRepository(client *ent.Client) *StoreUserRepository {
	return &StoreUserRepository{
		Client: client,
	}
}

func (repo *StoreUserRepository) Create(ctx context.Context, user *domain.StoreUser) (err error) {
	return nil
}

func (repo *StoreUserRepository) FindByUsername(ctx context.Context, username string) (u *domain.StoreUser, err error) {
	return nil, nil
}

func (repo *StoreUserRepository) Find(ctx context.Context, id uuid.UUID) (u *domain.StoreUser, err error) {
	return nil, nil
}

func (repo *StoreUserRepository) Exists(ctx context.Context, username string) (exists bool, err error) {
	return false, nil
}

func (repo *StoreUserRepository) Update(ctx context.Context, user *domain.StoreUser) (err error) {
	return nil
}
