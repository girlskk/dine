package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/storeuser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
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
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.StoreUser.Create().SetID(user.ID).
		SetUsername(user.Username).
		SetHashedPassword(user.HashedPassword).
		SetNickname(user.Nickname).
		SetMerchantID(user.MerchantID).
		SetStoreID(user.StoreID).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = domain.ConflictError(err)
		}
		err = fmt.Errorf("failed to create store user: %w", err)
		return
	}

	return nil
}

func (repo *StoreUserRepository) FindByUsername(ctx context.Context, username string) (u *domain.StoreUser, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.FindByUsername")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.StoreUser.Query().
		Where(storeuser.Username(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return
	}

	u = convertStoreUser(eu)
	return
}

func (repo *StoreUserRepository) Find(ctx context.Context, id uuid.UUID) (u *domain.StoreUser, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.StoreUser.Query().
		Where(storeuser.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return
	}

	u = convertStoreUser(eu)
	return
}

func (repo *StoreUserRepository) Exists(ctx context.Context, username string) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	exists, err = repo.Client.StoreUser.Query().
		Where(storeuser.Username(username)).
		Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check store user existence: %w", err)
		return
	}
	return
}

func (repo *StoreUserRepository) Update(ctx context.Context, user *domain.StoreUser) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.StoreUser.UpdateOneID(user.ID).
		SetUsername(user.Username).
		SetHashedPassword(user.HashedPassword).
		SetNickname(user.Nickname).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = domain.ConflictError(err)
		}
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update store user: %w", err)
		return
	}

	return nil
}

func convertStoreUser(eu *ent.StoreUser) *domain.StoreUser {
	if eu == nil {
		return nil
	}

	su := &domain.StoreUser{
		ID:             eu.ID,
		Username:       eu.Username,
		HashedPassword: eu.HashedPassword,
		Nickname:       eu.Nickname,
		MerchantID:     eu.MerchantID,
		StoreID:        eu.StoreID,
	}

	if eu.Edges.Store != nil {
		su.Store = convertStore(eu.Edges.Store)
	}

	return su
}
