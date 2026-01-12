package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/storeuser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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

	builder := repo.Client.StoreUser.Create().SetID(user.ID).
		SetUsername(user.Username).
		SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword).
		SetMerchantID(user.MerchantID).
		SetStoreID(user.StoreID).
		SetCode(user.Code).
		SetRealName(user.RealName).
		SetGender(user.Gender).
		SetEmail(user.Email).
		SetPhoneNumber(user.PhoneNumber).
		SetEnabled(user.Enabled).
		SetIsSuperadmin(user.IsSuperAdmin)

	if user.DepartmentID != uuid.Nil {
		builder = builder.SetDepartmentID(user.DepartmentID)
	}

	_, err = builder.Save(ctx)
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
		WithDepartment().
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

func (repo *StoreUserRepository) Exists(ctx context.Context, params domain.StoreUserExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.StoreUser.Query().
		Where(storeuser.Username(params.Username))
	if params.ExcludeID != uuid.Nil {
		builder = builder.Where(storeuser.IDNEQ(params.ExcludeID))
	}
	exists, err = builder.Exist(ctx)
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

	builder := repo.Client.StoreUser.UpdateOneID(user.ID).
		SetUsername(user.Username).
		SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword).
		SetRealName(user.RealName).
		SetGender(user.Gender).
		SetEmail(user.Email).
		SetPhoneNumber(user.PhoneNumber).
		SetEnabled(user.Enabled).
		SetIsSuperadmin(user.IsSuperAdmin)

	if user.DepartmentID != uuid.Nil {
		builder = builder.SetDepartmentID(user.DepartmentID)
	}
	_, err = builder.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = domain.ConflictError(err)
			return
		}
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to update store user: %w", err)
		return
	}

	return nil
}

func (repo *StoreUserRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.StoreUser.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to delete store user: %w", err)
		return
	}
	return nil
}

func (repo *StoreUserRepository) GetUsers(
	ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.StoreUserListFilter,
	orderBys ...domain.StoreUserOrderBy,
) (users []*domain.StoreUser, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "StoreUserRepository.GetUsers")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count store users: %w", err)
		return
	}

	list, err := query.
		Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query store users: %w", err)
		return
	}

	users = lo.Map(list, func(item *ent.StoreUser, _ int) *domain.StoreUser {
		return convertStoreUser(item)
	})
	return
}

func (repo *StoreUserRepository) buildFilterQuery(filter *domain.StoreUserListFilter) *ent.StoreUserQuery {
	query := repo.Client.StoreUser.Query()
	if filter == nil {
		return query
	}

	if len(filter.UserIDs) > 0 {
		query = query.Where(storeuser.IDIn(filter.UserIDs...))
	}
	if filter.Code != "" {
		query = query.Where(storeuser.CodeContains(filter.Code))
	}
	if filter.RealName != "" {
		query = query.Where(storeuser.RealNameContains(filter.RealName))
	}
	if filter.Gender != "" {
		query = query.Where(storeuser.GenderEQ(filter.Gender))
	}
	if filter.Email != "" {
		query = query.Where(storeuser.EmailContains(filter.Email))
	}
	if filter.PhoneNumber != "" {
		query = query.Where(storeuser.PhoneNumberContains(filter.PhoneNumber))
	}
	if filter.Enabled != nil {
		query = query.Where(storeuser.EnabledEQ(*filter.Enabled))
	}

	return query
}

func (repo *StoreUserRepository) orderBy(orderBys ...domain.StoreUserOrderBy) []storeuser.OrderOption {
	var opts []storeuser.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.Ternary(orderBy.Desc, sql.OrderDesc(), sql.OrderAsc())
		switch orderBy.OrderBy {
		case domain.StoreUserOrderByID:
			opts = append(opts, storeuser.ByID(rule))
		case domain.StoreUserOrderByCreatedAt:
			opts = append(opts, storeuser.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, storeuser.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
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
		Code:           eu.Code,
		RealName:       eu.RealName,
		Gender:         eu.Gender,
		Email:          eu.Email,
		PhoneNumber:    eu.PhoneNumber,
		Enabled:        eu.Enabled,
		IsSuperAdmin:   eu.IsSuperadmin,
		CreatedAt:      eu.CreatedAt,
		UpdatedAt:      eu.UpdatedAt,
	}

	if eu.Edges.Store != nil {
		su.Store = convertStore(eu.Edges.Store)
	}
	if eu.Edges.Department != nil {
		su.Department = convertDepartmentToDomain(eu.Edges.Department)
	}
	return su
}
