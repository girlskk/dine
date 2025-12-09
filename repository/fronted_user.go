package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/frontenduser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.FrontendUserRepository = (*FrontendUserRepository)(nil)

type FrontendUserRepository struct {
	Client *ent.Client
}

func NewFrontendUserRepository(client *ent.Client) *FrontendUserRepository {
	return &FrontendUserRepository{
		Client: client,
	}
}

func (repo *FrontendUserRepository) FindByUsername(ctx context.Context, username string) (u *domain.FrontendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserRepository.FindByUsername")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.FrontendUser.Query().
		Where(frontenduser.Username(username)).
		WithStore().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return
	}

	u = convertFrontendUser(eu)

	return
}

func (repo *FrontendUserRepository) Find(ctx context.Context, id int) (u *domain.FrontendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	eu, err := repo.Client.FrontendUser.Query().
		Where(frontenduser.ID(id)).
		WithStore().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return
	}

	u = convertFrontendUser(eu)

	return
}

func convertFrontendUser(eu *ent.FrontendUser) *domain.FrontendUser {
	if eu == nil {
		return nil
	}

	return &domain.FrontendUser{
		ID:             eu.ID,
		Username:       eu.Username,
		HashedPassword: eu.HashedPassword,
		Nickname:       eu.Nickname,
		StoreID:        eu.StoreID,
		Store:          convertStore(eu.Edges.Store),
	}
}

func (repo *FrontendUserRepository) Exists(ctx context.Context, username string) (exists bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.FrontendUser.Query().Where(frontenduser.Username(username))
	return query.Exist(ctx)
}

func (repo *FrontendUserRepository) Create(ctx context.Context, user *domain.FrontendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	_, err = repo.Client.FrontendUser.Create().
		SetUsername(user.Username).
		SetHashedPassword(user.HashedPassword).
		SetNickname(user.Nickname).
		SetStoreID(user.StoreID).
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

func (repo *FrontendUserRepository) Update(ctx context.Context, user *domain.FrontendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	pwd := lo.Ternary(user.HashedPassword != "", &user.HashedPassword, nil)

	_, err = repo.Client.FrontendUser.UpdateOneID(user.ID).
		SetUsername(user.Username).
		SetNillableHashedPassword(pwd).
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

func (repo *FrontendUserRepository) List(ctx context.Context, pager *upagination.Pagination, filter *domain.FrontendUserListFilter) (dusers []*domain.FrontendUser, total int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserRepository.List")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	query := repo.filterBuildQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count: %w", err)
		return
	}

	users, err := query.Order(frontenduser.ByCreatedAt(sql.OrderDesc()), frontenduser.ByID(sql.OrderDesc())).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get users: %w", err)
		return
	}

	dusers = lo.Map(users, func(user *ent.FrontendUser, _ int) *domain.FrontendUser {
		return convertFrontendUser(user)
	})

	return
}

func (repo *FrontendUserRepository) filterBuildQuery(filter *domain.FrontendUserListFilter) *ent.FrontendUserQuery {
	query := repo.Client.FrontendUser.Query()

	if filter.StoreID > 0 {
		query = query.Where(frontenduser.StoreID(filter.StoreID))
	}

	return query
}
