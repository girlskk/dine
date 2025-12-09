package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/ent/backenduser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
)

var _ domain.BackendUserRepository = (*BackendUserRepository)(nil)

type BackendUserRepository struct {
	Client *ent.Client
}

func NewBackendUserRepository(client *ent.Client) *BackendUserRepository {
	return &BackendUserRepository{
		Client: client,
	}
}

func (repo *BackendUserRepository) Create(ctx context.Context, user *domain.BackendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.BackendUser.Create().
		SetUsername(user.Username).
		SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword).
		SetStoreID(user.StoreID).
		Save(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (repo *BackendUserRepository) FindByUsername(ctx context.Context, username string) (u *domain.BackendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserRepository.FindByUsername")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.BackendUser.Query().
		Where(backenduser.Username(username)).
		WithStore().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrUserNotFound)
		}
		return
	}

	u = convertBackendUser(eu)

	return
}

func (repo *BackendUserRepository) FindByStoreID(ctx context.Context, storeID int) (u *domain.BackendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserRepository.FindByStoreID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.BackendUser.Query().
		Where(backenduser.StoreID(storeID)).
		WithStore().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrUserNotFound)
		}
		return
	}

	u = convertBackendUser(eu)

	return
}

func (repo *BackendUserRepository) Find(ctx context.Context, id int) (u *domain.BackendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.BackendUser.Query().
		Where(backenduser.ID(id)).
		WithStore().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return
	}

	u = convertBackendUser(eu)

	return
}

func (repo *BackendUserRepository) Exists(ctx context.Context, username string) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.BackendUser.Query().Where(backenduser.Username(username))
	return query.Exist(ctx)
}

func (repo *BackendUserRepository) Update(ctx context.Context, user *domain.BackendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	update := repo.Client.BackendUser.Update().
		//SetUsername(user.Username).
		//SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword)

	if user.ID > 0 {
		update.Where(backenduser.ID(user.ID))
	} else {
		update.Where(backenduser.StoreID(user.StoreID))
	}
	_, err = update.Save(ctx)
	return err
}

func convertBackendUser(eu *ent.BackendUser) *domain.BackendUser {
	if eu == nil {
		return nil
	}

	return &domain.BackendUser{
		ID:             eu.ID,
		Username:       eu.Username,
		HashedPassword: eu.HashedPassword,
		Nickname:       eu.Nickname,
		StoreID:        eu.StoreID,
		Store:          convertStore(eu.Edges.Store),
	}
}
