package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/permission"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PermissionRepository = (*PermissionRepository)(nil)

// PermissionRepository implements Permission CRUD and pagination.
type PermissionRepository struct {
	Client *ent.Client
}

func NewPermissionRepository(client *ent.Client) *PermissionRepository {
	return &PermissionRepository{Client: client}
}

func (repo *PermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (domainPermission *domain.Permission, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PermissionRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	ep, err := repo.Client.Permission.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrPermissionNotExists)
			return
		}
		return
	}
	domainPermission = convertPermissionToDomain(ep)
	return
}

func (repo *PermissionRepository) FindByCode(ctx context.Context, code string) (domainPermission *domain.Permission, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PermissionRepository.FindByCode")
	defer func() { util.SpanErrFinish(span, err) }()

	ep, err := repo.Client.Permission.Query().
		Where(permission.PermCode(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrPermissionNotExists)
			return
		}
		return
	}
	domainPermission = convertPermissionToDomain(ep)
	return
}

func (repo *PermissionRepository) Create(ctx context.Context, domainPermission *domain.Permission) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PermissionRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if domainPermission == nil {
		return fmt.Errorf("permission is nil")
	}

	builder := repo.Client.Permission.Create().
		SetID(domainPermission.ID).
		SetMenuID(domainPermission.MenuID).
		SetPermCode(domainPermission.PermCode).
		SetName(domainPermission.Name).
		SetMethod(domainPermission.Method).
		SetPath(domainPermission.Path).
		SetEnabled(domainPermission.Enabled)

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create permission: %w", err)
		return
	}

	domainPermission.ID = created.ID
	domainPermission.CreatedAt = created.CreatedAt
	domainPermission.UpdatedAt = created.UpdatedAt
	return
}

func (repo *PermissionRepository) Update(ctx context.Context, p *domain.Permission) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PermissionRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if p == nil {
		return fmt.Errorf("permission is nil")
	}

	builder := repo.Client.Permission.UpdateOneID(p.ID).
		SetMenuID(p.MenuID).
		SetPermCode(p.PermCode).
		SetName(p.Name).
		SetMethod(p.Method).
		SetPath(p.Path).
		SetEnabled(p.Enabled)

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrPermissionNotExists)
			return
		}
		err = fmt.Errorf("failed to update permission: %w", err)
		return
	}

	p.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PermissionRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Permission.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete permission: %w", err)
		return
	}
	return nil
}

func (repo *PermissionRepository) GetPermissions(ctx context.Context, pager *upagination.Pagination, filter *domain.PermissionListFilter, orderBys ...domain.PermissionListOrderBy) (permissions []*domain.Permission, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PermissionRepository.GetPermissions")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count permission: %w", err)
		return
	}

	list, err := query.
		Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query permission: %w", err)
		return
	}

	permissions = lo.Map(list, func(item *ent.Permission, _ int) *domain.Permission {
		return convertPermissionToDomain(item)
	})
	return
}

func (repo *PermissionRepository) Exists(ctx context.Context, params domain.PermissionExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "PermissionRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Permission.Query()

	if params.PermCode != "" {
		query = query.Where(permission.PermCode(params.PermCode))
	}
	if params.Method != "" {
		query = query.Where(permission.Method(params.Method))
	}
	if params.Path != "" {
		query = query.Where(permission.Path(params.Path))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(permission.IDNEQ(params.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check permission existence: %w", err)
	}
	return
}

func (repo *PermissionRepository) buildFilterQuery(filter *domain.PermissionListFilter) *ent.PermissionQuery {
	query := repo.Client.Permission.Query()
	if filter == nil {
		return query
	}

	if filter.MenuID != uuid.Nil {
		query = query.Where(permission.MenuID(filter.MenuID))
	}
	if filter.PermCode != "" {
		query = query.Where(permission.PermCodeContains(filter.PermCode))
	}
	if filter.Name != "" {
		query = query.Where(permission.NameContains(filter.Name))
	}
	if filter.Method != "" {
		query = query.Where(permission.Method(filter.Method))
	}
	if filter.Path != "" {
		query = query.Where(permission.PathContains(filter.Path))
	}
	if filter.Enabled != nil {
		query = query.Where(permission.EnabledEQ(*filter.Enabled))
	}

	return query
}

func (repo *PermissionRepository) orderBy(orderBys ...domain.PermissionListOrderBy) []permission.OrderOption {
	var opts []permission.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.PermissionListOrderByID:
			opts = append(opts, permission.ByID(rule))
		case domain.PermissionListOrderByCreatedAt:
			opts = append(opts, permission.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, permission.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertPermissionToDomain(ep *ent.Permission) *domain.Permission {
	if ep == nil {
		return nil
	}
	return &domain.Permission{
		ID:        ep.ID,
		MenuID:    ep.MenuID,
		PermCode:  ep.PermCode,
		Name:      ep.Name,
		Method:    ep.Method,
		Path:      ep.Path,
		Enabled:   ep.Enabled,
		CreatedAt: ep.CreatedAt,
		UpdatedAt: ep.UpdatedAt,
	}
}
