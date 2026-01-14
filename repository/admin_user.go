package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/adminuser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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
		WithDepartment().
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

	builder := repo.Client.AdminUser.Create().SetID(user.ID).
		SetUsername(user.Username).
		SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword).
		SetCode(user.Code).
		SetRealName(user.RealName).
		SetEmail(user.Email).
		SetPhoneNumber(user.PhoneNumber).
		SetEnabled(user.Enabled).
		SetIsSuperadmin(user.IsSuperAdmin)

	if user.DepartmentID != uuid.Nil {
		builder = builder.SetDepartmentID(user.DepartmentID)
	}
	if user.Gender != "" {
		builder = builder.SetGender(user.Gender)
	} else {
		builder = builder.SetGender(domain.GenderUnknown)
	}
	_, err = builder.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = domain.ConflictError(err)
			return
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

	builder := repo.Client.AdminUser.UpdateOneID(user.ID).
		SetUsername(user.Username).
		SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword).
		SetRealName(user.RealName).
		SetEmail(user.Email).
		SetPhoneNumber(user.PhoneNumber).
		SetEnabled(user.Enabled).
		SetIsSuperadmin(user.IsSuperAdmin)
	if user.DepartmentID != uuid.Nil {
		builder = builder.SetDepartmentID(user.DepartmentID)
	}
	if user.Gender != "" {
		builder = builder.SetGender(user.Gender)
	} else {
		builder = builder.SetGender(domain.GenderUnknown)
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
			return
		}
		err = fmt.Errorf("failed to delete user: %w", err)
		return
	}
	return nil
}

func (repo *AdminUserRepository) Exists(ctx context.Context, params domain.AdminUserExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "AdminUserRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.AdminUser.Query().Where(adminuser.DeletedAt(0))
	if params.Username != "" {
		query = query.Where(adminuser.Username(params.Username))
	}
	if params.Code != "" {
		query = query.Where(adminuser.Code(params.Code))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(adminuser.IDNEQ(params.ExcludeID))
	}
	exists, err = query.Exist(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check admin user exists: %w", err)
	}
	return
}

func (repo *AdminUserRepository) GetUsers(
	ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.AdminUserListFilter,
	orderBys ...domain.AdminUserOrderBy,
) (users []*domain.AdminUser, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "AdminUserRepository.GetUsers")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count admin users: %w", err)
		return
	}

	list, err := query.
		Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query admin users: %w", err)
		return
	}

	users = lo.Map(list, func(item *ent.AdminUser, _ int) *domain.AdminUser {
		return convertAdminUser(item)
	})
	return
}

func (repo *AdminUserRepository) buildFilterQuery(filter *domain.AdminUserListFilter) *ent.AdminUserQuery {
	query := repo.Client.AdminUser.Query().Where(adminuser.DeletedAt(0))
	if filter == nil {
		return query
	}

	if len(filter.UserIDs) > 0 {
		query = query.Where(adminuser.IDIn(filter.UserIDs...))
	}
	if filter.Code != "" {
		query = query.Where(adminuser.CodeContains(filter.Code))
	}
	if filter.RealName != "" {
		query = query.Where(adminuser.RealNameContains(filter.RealName))
	}
	if filter.Gender != "" {
		query = query.Where(adminuser.GenderEQ(filter.Gender))
	}
	if filter.Email != "" {
		query = query.Where(adminuser.EmailContains(filter.Email))
	}
	if filter.PhoneNumber != "" {
		query = query.Where(adminuser.PhoneNumberContains(filter.PhoneNumber))
	}
	if filter.Enabled != nil {
		query = query.Where(adminuser.EnabledEQ(*filter.Enabled))
	}
	return query
}

func (repo *AdminUserRepository) orderBy(orderBys ...domain.AdminUserOrderBy) []adminuser.OrderOption {
	var opts []adminuser.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.Ternary(orderBy.Desc, sql.OrderDesc(), sql.OrderAsc())
		switch orderBy.OrderBy {
		case domain.AdminUserOrderByID:
			opts = append(opts, adminuser.ByID(rule))
		case domain.AdminUserOrderByCreatedAt:
			opts = append(opts, adminuser.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, adminuser.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertAdminUser(eu *ent.AdminUser) *domain.AdminUser {
	if eu == nil {
		return nil
	}
	du := &domain.AdminUser{
		ID:             eu.ID,
		Username:       eu.Username,
		HashedPassword: eu.HashedPassword,
		Nickname:       eu.Nickname,
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
	if eu.Edges.Department != nil {
		du.Department = convertDepartmentToDomain(eu.Edges.Department)
	}
	return du
}
