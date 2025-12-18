package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/backenduser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
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
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.BackendUser.Create().
		SetUsername(user.Username).
		SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword).
		Save(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (repo *BackendUserRepository) FindByUsername(ctx context.Context, username string) (u *domain.BackendUser, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.FindByUsername")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.BackendUser.Query().
		Where(backenduser.Username(username)).
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

func (repo *BackendUserRepository) Find(ctx context.Context, id uuid.UUID) (u *domain.BackendUser, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.BackendUser.Query().
		Where(backenduser.ID(id)).
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
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.BackendUser.Query().Where(backenduser.Username(username))
	return query.Exist(ctx)
}

func (repo *BackendUserRepository) Update(ctx context.Context, user *domain.BackendUser) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	update := repo.Client.BackendUser.Update().
		SetHashedPassword(user.HashedPassword)

	if user.ID != uuid.Nil {
		update.Where(backenduser.ID(user.ID))
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
		MerchantID:     eu.MerchantID,
		Username:       eu.Username,
		HashedPassword: eu.HashedPassword,
		Nickname:       eu.Nickname,
	}
}
