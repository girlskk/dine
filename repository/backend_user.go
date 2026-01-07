package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/backenduser"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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

	builder := repo.Client.BackendUser.Create().SetID(user.ID).
		SetUsername(user.Username).
		SetNickname(user.Nickname).
		SetHashedPassword(user.HashedPassword).
		SetMerchantID(user.MerchantID).
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
			return
		}
		err = fmt.Errorf("failed to create user: %w", err)
		return
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
		WithDepartment().
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

func (repo *BackendUserRepository) Exists(ctx context.Context, params domain.BackendUserExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.Exists")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.Client.BackendUser.Query().
		Where(backenduser.Username(params.Username))
	if params.ExcludeID != uuid.Nil {
		query = query.Where(backenduser.IDNEQ(params.ExcludeID))
	}
	return query.Exist(ctx)
}

func (repo *BackendUserRepository) Update(ctx context.Context, user *domain.BackendUser) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	builder := repo.Client.BackendUser.UpdateOneID(user.ID).
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
		}
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		err = fmt.Errorf("failed to update backend user: %w", err)
		return
	}
	return err
}
func (repo *BackendUserRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = repo.Client.BackendUser.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return err
		}
		err = fmt.Errorf("failed to delete backend user: %w", err)
		return err
	}
	return nil
}

func (repo *BackendUserRepository) GetUsers(ctx context.Context, pager *upagination.Pagination, filter *domain.BackendUserListFilter, orderBys ...domain.BackendUserOrderBy) (users []*domain.BackendUser, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "BackendUserRepository.GetUsers")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count backend users: %w", err)
		return
	}

	eUsers, err := query.
		Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query backend users: %w", err)
		return
	}

	users = lo.Map(eUsers, func(item *ent.BackendUser, _ int) *domain.BackendUser {
		return convertBackendUser(item)
	})
	return
}
func convertBackendUser(eu *ent.BackendUser) *domain.BackendUser {
	if eu == nil {
		return nil
	}
	du := &domain.BackendUser{
		ID:             eu.ID,
		MerchantID:     eu.MerchantID,
		DepartmentID:   eu.DepartmentID,
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

func (repo *BackendUserRepository) orderBy(orderBys ...domain.BackendUserOrderBy) []backenduser.OrderOption {
	var opts []backenduser.OrderOption

	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.BackendUserOrderByID:
			opts = append(opts, backenduser.ByID(rule))
		case domain.BackendUserOrderByCreatedAt:
			opts = append(opts, backenduser.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, backenduser.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func (repo *BackendUserRepository) buildFilterQuery(filter *domain.BackendUserListFilter) *ent.BackendUserQuery {
	query := repo.Client.BackendUser.Query()
	if filter == nil {
		return query
	}

	if len(filter.UserIDs) > 0 {
		query = query.Where(backenduser.IDIn(filter.UserIDs...))
	}
	if filter.Code != "" {
		query = query.Where(backenduser.CodeEQ(filter.Code))
	}
	if filter.RealName != "" {
		query = query.Where(backenduser.RealNameEQ(filter.RealName))
	}
	if filter.Gender != "" {
		query = query.Where(backenduser.GenderEQ(filter.Gender))
	}
	if filter.Email != "" {
		query = query.Where(backenduser.EmailEQ(filter.Email))
	}
	if filter.PhoneNumber != "" {
		query = query.Where(backenduser.PhoneNumberEQ(filter.PhoneNumber))
	}
	if filter.Enabled != nil {
		query = query.Where(backenduser.EnabledEQ(*filter.Enabled))
	}
	return query
}
