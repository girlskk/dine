package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/userrole"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.UserRoleRepository = (*UserRoleRepository)(nil)

// UserRoleRepository implements UserRole relation persistence.
type UserRoleRepository struct {
	Client *ent.Client
}

func NewUserRoleRepository(client *ent.Client) *UserRoleRepository {
	return &UserRoleRepository{Client: client}
}

func (repo *UserRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (ur *domain.UserRole, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	eur, err := repo.Client.UserRole.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrUserRoleNotExists)
			return
		}
		return
	}
	ur = convertUserRoleToDomain(eur)
	return
}

func (repo *UserRoleRepository) FindOneByUser(ctx context.Context, user domain.User) (ur *domain.UserRole, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.FindOneByUser")
	defer func() { util.SpanErrFinish(span, err) }()

	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	eur, err := repo.Client.UserRole.Query().
		Where(
			userrole.UserTypeEQ(user.GetUserType()),
			userrole.UserID(user.GetUserID()),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NotFoundError(domain.ErrUserRoleNotExists)
		}
		return nil, err
	}
	ur = convertUserRoleToDomain(eur)
	return
}

func (repo *UserRoleRepository) Create(ctx context.Context, ur *domain.UserRole) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if ur == nil {
		return fmt.Errorf("user role is nil")
	}

	builder := repo.Client.UserRole.Create().
		SetID(ur.ID).
		SetUserType(ur.UserType).
		SetUserID(ur.UserID).
		SetRoleID(ur.RoleID)

	if ur.MerchantID != uuid.Nil {
		builder = builder.SetMerchantID(ur.MerchantID)
	}
	if ur.StoreID != uuid.Nil {
		builder = builder.SetStoreID(ur.StoreID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create user role: %w", err)
		return
	}

	ur.ID = created.ID
	ur.CreatedAt = created.CreatedAt
	ur.UpdatedAt = created.UpdatedAt
	return
}

func (repo *UserRoleRepository) CreateBulk(ctx context.Context, relations []*domain.UserRole) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.CreateBulk")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(relations) == 0 {
		return nil
	}

	builders := make([]*ent.UserRoleCreate, 0, len(relations))
	for _, ur := range relations {
		if ur == nil {
			continue
		}
		builder := repo.Client.UserRole.Create().
			SetID(ur.ID).
			SetUserType(ur.UserType).
			SetUserID(ur.UserID).
			SetRoleID(ur.RoleID)
		if ur.MerchantID != uuid.Nil {
			builder = builder.SetMerchantID(ur.MerchantID)
		}
		if ur.StoreID != uuid.Nil {
			builder = builder.SetStoreID(ur.StoreID)
		}
		builders = append(builders, builder)
	}

	if len(builders) == 0 {
		return nil
	}

	_, err = repo.Client.UserRole.CreateBulk(builders...).Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to bulk create user roles: %w", err)
	}
	return
}

func (repo *UserRoleRepository) CreateBulkByRoleIDUsers(ctx context.Context, roleID uuid.UUID, users []domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.CreateBulkByRoleIDUsers")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(users) == 0 {
		return nil
	}

	builders := make([]*ent.UserRoleCreate, 0, len(users))
	for _, u := range users {
		if u == nil {
			continue
		}

		builder := repo.Client.UserRole.Create().
			SetID(uuid.New()).
			SetUserType(u.GetUserType()).
			SetUserID(u.GetUserID()).
			SetRoleID(roleID)

		if mid := u.GetMerchantID(); mid != uuid.Nil {
			builder = builder.SetMerchantID(mid)
		}
		if sid := u.GetStoreID(); sid != uuid.Nil {
			builder = builder.SetStoreID(sid)
		}

		builders = append(builders, builder)
	}

	if len(builders) == 0 {
		return nil
	}

	_, err = repo.Client.UserRole.CreateBulk(builders...).Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to bulk create user roles by role/users: %w", err)
	}
	return
}

func (repo *UserRoleRepository) CreateBulkByUserIDRoles(ctx context.Context, user domain.User, roles []uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.CreateBulkByUserIDRoles")
	defer func() { util.SpanErrFinish(span, err) }()

	if user == nil || len(roles) == 0 {
		return nil
	}

	builders := make([]*ent.UserRoleCreate, 0, len(roles))
	for _, roleID := range roles {
		builder := repo.Client.UserRole.Create().
			SetID(uuid.New()).
			SetUserType(user.GetUserType()).
			SetUserID(user.GetUserID()).
			SetRoleID(roleID)

		if mid := user.GetMerchantID(); mid != uuid.Nil {
			builder = builder.SetMerchantID(mid)
		}
		if sid := user.GetStoreID(); sid != uuid.Nil {
			builder = builder.SetStoreID(sid)
		}
		builders = append(builders, builder)
	}

	_, err = repo.Client.UserRole.CreateBulk(builders...).Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to bulk create user roles by user/roles: %w", err)
	}
	return
}
func (repo *UserRoleRepository) Update(ctx context.Context, userRole *domain.UserRole) error {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.Update")
	defer func() { util.SpanErrFinish(span, nil) }()

	_, err := repo.Client.UserRole.UpdateOneID(userRole.ID).
		SetUserType(userRole.UserType).
		SetRoleID(userRole.RoleID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}

func (repo *UserRoleRepository) Deletes(ctx context.Context, ids ...uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.Deletes")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(ids) == 0 {
		return nil
	}

	_, err = repo.Client.UserRole.Delete().Where(userrole.IDIn(ids...)).Exec(ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete user roles: %w", err)
	}
	return
}

func (repo *UserRoleRepository) DeleteByRoles(ctx context.Context, roleIDs ...uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.DeleteByRoles")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(roleIDs) == 0 {
		return nil
	}

	_, err = repo.Client.UserRole.Delete().Where(userrole.RoleIDIn(roleIDs...)).Exec(ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete user roles by roles: %w", err)
	}
	return
}

func (repo *UserRoleRepository) DeleteByUsers(ctx context.Context, userType domain.UserType, userIDs ...uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.DeleteByUsers")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(userIDs) == 0 {
		return nil
	}

	_, err = repo.Client.UserRole.Delete().Where(
		userrole.UserTypeEQ(userType),
		userrole.UserIDIn(userIDs...),
	).Exec(ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete user roles by users: %w", err)
	}
	return
}

func (repo *UserRoleRepository) GetUserRoles(ctx context.Context, pager *upagination.Pagination, filter *domain.UserRoleListFilter, orderBys ...domain.UserRoleListOrderBy) (list []*domain.UserRole, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.GetUserRoles")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count user roles: %w", err)
		return
	}

	items, err := query.Order(repo.orderBy(orderBys...)...).Offset(pager.Offset()).Limit(pager.Size).All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query user roles: %w", err)
		return
	}

	list = lo.Map(items, func(eur *ent.UserRole, _ int) *domain.UserRole { return convertUserRoleToDomain(eur) })
	return
}

func (repo *UserRoleRepository) GetByRoleIDs(ctx context.Context, userType domain.UserType, roleIDs ...uuid.UUID) (userRoles []*domain.UserRole, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.GetUsersByRoleID")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(roleIDs) == 0 {
		return
	}
	eUserRoles, err := repo.Client.UserRole.Query().
		Where(
			userrole.RoleIDIn(roleIDs...),
			userrole.UserTypeEQ(userType),
		).All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query users by role: %w", err)
		return
	}
	userRoles = lo.Map(eUserRoles, func(item *ent.UserRole, _ int) *domain.UserRole {
		return convertUserRoleToDomain(item)
	})
	return
}

func (repo *UserRoleRepository) GetByUserIDs(ctx context.Context, userType domain.UserType, userIDs ...uuid.UUID) (userRoles []*domain.UserRole, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "UserRoleRepository.GetRolesByUserID")
	defer func() { util.SpanErrFinish(span, err) }()

	eUserRoles, err := repo.Client.UserRole.Query().
		Where(
			userrole.UserIDIn(userIDs...),
			userrole.UserTypeEQ(userType),
		).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query roles by user: %w", err)
		return
	}
	userRoles = lo.Map(eUserRoles, func(item *ent.UserRole, _ int) *domain.UserRole {
		return convertUserRoleToDomain(item)
	})
	return
}

func (repo *UserRoleRepository) buildFilterQuery(filter *domain.UserRoleListFilter) *ent.UserRoleQuery {
	query := repo.Client.UserRole.Query()
	if filter == nil {
		return query
	}
	if filter.UserType != "" {
		query = query.Where(userrole.UserTypeEQ(filter.UserType))
	}
	if filter.UserID != uuid.Nil {
		query = query.Where(userrole.UserID(filter.UserID))
	}
	if filter.RoleID != uuid.Nil {
		query = query.Where(userrole.RoleID(filter.RoleID))
	}
	if filter.MerchantID != uuid.Nil {
		query = query.Where(userrole.MerchantID(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(userrole.StoreID(filter.StoreID))
	}
	return query
}

func (repo *UserRoleRepository) orderBy(orderBys ...domain.UserRoleListOrderBy) []userrole.OrderOption {
	var opts []userrole.OrderOption
	for _, ob := range orderBys {
		rule := lo.TernaryF(ob.Desc, sql.OrderDesc, sql.OrderAsc)
		switch ob.OrderBy {
		case domain.UserRoleListOrderByID:
			opts = append(opts, userrole.ByID(rule))
		case domain.UserRoleListOrderByCreatedAt:
			opts = append(opts, userrole.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, userrole.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertUserRoleToDomain(eur *ent.UserRole) *domain.UserRole {
	if eur == nil {
		return nil
	}
	return &domain.UserRole{
		ID:         eur.ID,
		UserType:   eur.UserType,
		UserID:     eur.UserID,
		RoleID:     eur.RoleID,
		MerchantID: eur.MerchantID,
		StoreID:    eur.StoreID,
		CreatedAt:  eur.CreatedAt,
		UpdatedAt:  eur.UpdatedAt,
	}
}
