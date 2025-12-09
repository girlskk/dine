package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AdminUserRepository.FindByUsername")
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

func (repo *AdminUserRepository) Find(ctx context.Context, id int) (u *domain.AdminUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AdminUserRepository.Find")
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
