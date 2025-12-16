package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/adminuser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.AdminUserRepository = (*AdminUserRepository)(nil)

type AdminUserRepository struct {
	Client *ent.Client
}

func NewAdminUserRepository(client *ent.Client) *AdminUserRepository {
	return &AdminUserRepository{
		Client: client,
	}
}

func (repo *AdminUserRepository) FindByUsername(ctx context.Context, username string) (u *domain.AdminUser, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "AdminUserRepository.FindByUsername")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.AdminUser.Query().
		Where(adminuser.Username(username)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return
	}

	u = convertAdminUser(eu)

	return
}

func (repo *AdminUserRepository) Find(ctx context.Context, id uuid.UUID) (u *domain.AdminUser, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "AdminUserRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.AdminUser.Query().
		Where(adminuser.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return
	}

	u = convertAdminUser(eu)

	return
}

func (repo *AdminUserRepository) Create(ctx context.Context, user *domain.AdminUser) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "AdminUserRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.AdminUser.Create().SetID(user.ID).
		SetUsername(user.Username).
		SetHashedPassword(user.HashedPassword).
		SetNickname(user.Nickname).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			err = domain.ConflictError(err)
		}
		err = fmt.Errorf("failed to create user: %w", err)
		return
	}

	return nil
}

func (repo *AdminUserRepository) Update(ctx context.Context, user *domain.AdminUser) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "AdminUserRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	_, err = repo.Client.AdminUser.UpdateOneID(user.ID).
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
		err = fmt.Errorf("failed to update user: %w", err)
		return
	}

	return nil
}

func (repo *AdminUserRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "AdminUserRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.AdminUser.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to delete user: %w", err)
		return
	}
	return nil
}

func convertAdminUser(eu *ent.AdminUser) *domain.AdminUser {
	if eu == nil {
		return nil
	}

	return &domain.AdminUser{
		ID:             eu.ID,
		Username:       eu.Username,
		HashedPassword: eu.HashedPassword,
		Nickname:       eu.Nickname,
	}
}
