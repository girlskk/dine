package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/role"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RoleRepository = (*RoleRepository)(nil)

// RoleRepository implements Role CRUD and pagination.
type RoleRepository struct {
	Client *ent.Client
}

func NewRoleRepository(client *ent.Client) *RoleRepository {
	return &RoleRepository{Client: client}
}

func (repo *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (domainRole *domain.Role, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	er, err := repo.Client.Role.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrRoleNotExists)
			return
		}
		return
	}
	domainRole = convertRoleToDomain(er)
	return
}

func (repo *RoleRepository) Create(ctx context.Context, domainRole *domain.Role) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainRole == nil {
		return fmt.Errorf("role is nil")
	}

	builder := repo.Client.Role.Create().
		SetID(domainRole.ID).
		SetName(domainRole.Name).
		SetCode(domainRole.Code).
		SetRoleType(domainRole.RoleType).
		SetDataScope(domainRole.DataScope).
		SetEnable(domainRole.Enable).
		SetMerchantID(domainRole.MerchantID).
		SetStoreID(domainRole.StoreID)

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create role: %w", err)
		return
	}

	domainRole.ID = created.ID
	domainRole.CreatedAt = created.CreatedAt
	domainRole.UpdatedAt = created.UpdatedAt
	return
}

func (repo *RoleRepository) Update(ctx context.Context, domainRole *domain.Role) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainRole == nil {
		return fmt.Errorf("role is nil")
	}

	builder := repo.Client.Role.UpdateOneID(domainRole.ID).
		SetName(domainRole.Name).
		SetDataScope(domainRole.DataScope).
		SetEnable(domainRole.Enable)

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrRoleNotExists)
			return
		}
		err = fmt.Errorf("failed to update role: %w", err)
		return
	}

	domainRole.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *RoleRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Role.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete role: %w", err)
		return
	}
	return nil
}

func (repo *RoleRepository) GetRoles(ctx context.Context, pager *upagination.Pagination, filter *domain.RoleListFilter, orderBys ...domain.RoleListOrderBy) (roles []*domain.Role, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleRepository.GetRoles")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count role: %w", err)
		return
	}

	list, err := query.
		Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query role: %w", err)
		return
	}

	roles = lo.Map(list, func(item *ent.Role, _ int) *domain.Role {
		return convertRoleToDomain(item)
	})
	return
}

func (repo *RoleRepository) Exists(ctx context.Context, params domain.RoleExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Role.Query().Where(role.Name(params.Name))
	if params.Code != "" {
		query = query.Where(role.Code(params.Code))
	}
	if params.MerchantID != uuid.Nil {
		query = query.Where(role.MerchantID(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(role.StoreID(params.StoreID))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(role.IDNEQ(params.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check role existence: %w", err)
	}
	return
}

func (repo *RoleRepository) buildFilterQuery(filter *domain.RoleListFilter) *ent.RoleQuery {
	query := repo.Client.Role.Query()
	if filter == nil {
		return query
	}

	if filter.MerchantID != uuid.Nil {
		query = query.Where(role.MerchantID(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(role.StoreID(filter.StoreID))
	}
	if filter.RoleType != "" {
		query = query.Where(role.RoleTypeEQ(filter.RoleType))
	}
	if filter.Enable != nil {
		query = query.Where(role.EnableEQ(*filter.Enable))
	}
	if filter.Name != "" {
		query = query.Where(role.NameContains(filter.Name))
	}

	return query
}

func (repo *RoleRepository) orderBy(orderBys ...domain.RoleListOrderBy) []role.OrderOption {
	var opts []role.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.RoleListOrderByID:
			opts = append(opts, role.ByID(rule))
		case domain.RoleListOrderByCreatedAt:
			opts = append(opts, role.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, role.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertRoleToDomain(er *ent.Role) *domain.Role {
	if er == nil {
		return nil
	}
	return &domain.Role{
		ID:         er.ID,
		Name:       er.Name,
		Code:       er.Code,
		RoleType:   er.RoleType,
		DataScope:  er.DataScope,
		Enable:     er.Enable,
		MerchantID: er.MerchantID,
		StoreID:    er.StoreID,
		CreatedAt:  er.CreatedAt,
		UpdatedAt:  er.UpdatedAt,
	}
}
