package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/rolemenu"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RoleMenuRepository = (*RoleMenuRepository)(nil)

// RoleMenuRepository implements role-menu relation persistence.
type RoleMenuRepository struct {
	Client *ent.Client
}

func NewRoleMenuRepository(client *ent.Client) *RoleMenuRepository {
	return &RoleMenuRepository{Client: client}
}

func (repo *RoleMenuRepository) FindByID(ctx context.Context, id uuid.UUID) (rm *domain.RoleMenu, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	erm, err := repo.Client.RoleMenu.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrRoleMenuNotExists)
			return
		}
		return
	}
	rm = convertRoleMenuToDomain(erm)
	return
}

func (repo *RoleMenuRepository) Create(ctx context.Context, rm *domain.RoleMenu) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if rm == nil {
		return fmt.Errorf("role menu is nil")
	}

	builder := repo.Client.RoleMenu.Create().
		SetID(rm.ID).
		SetRoleType(rm.RoleType).
		SetRoleID(rm.RoleID).
		SetPath(rm.Path)

	if rm.MerchantID != uuid.Nil {
		builder = builder.SetMerchantID(rm.MerchantID)
	}
	if rm.StoreID != uuid.Nil {
		builder = builder.SetStoreID(rm.StoreID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create role menu: %w", err)
		return
	}

	rm.ID = created.ID
	rm.CreatedAt = created.CreatedAt
	rm.UpdatedAt = created.UpdatedAt
	return
}

func (repo *RoleMenuRepository) CreateBulk(ctx context.Context, relations []*domain.RoleMenu) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.CreateBulk")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(relations) == 0 {
		return nil
	}

	builders := make([]*ent.RoleMenuCreate, 0, len(relations))
	for _, rm := range relations {
		if rm == nil {
			continue
		}
		builder := repo.Client.RoleMenu.Create().
			SetID(rm.ID).
			SetRoleType(rm.RoleType).
			SetRoleID(rm.RoleID).
			SetPath(rm.Path)
		if rm.MerchantID != uuid.Nil {
			builder = builder.SetMerchantID(rm.MerchantID)
		}
		if rm.StoreID != uuid.Nil {
			builder = builder.SetStoreID(rm.StoreID)
		}
		builders = append(builders, builder)
	}

	if len(builders) == 0 {
		return nil
	}

	_, err = repo.Client.RoleMenu.CreateBulk(builders...).Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to bulk create role menus: %w", err)
	}
	return
}

func (repo *RoleMenuRepository) CreateBulkByRoleIDPaths(ctx context.Context, role *domain.Role, paths []string) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.CreateBulkByRoleIDPaths")
	defer func() { util.SpanErrFinish(span, err) }()

	if role == nil || len(paths) == 0 {
		return nil
	}

	builders := make([]*ent.RoleMenuCreate, 0, len(paths))
	for _, path := range paths {
		builder := repo.Client.RoleMenu.Create().
			SetRoleType(role.RoleType).
			SetRoleID(role.ID).
			SetPath(path)

		if role.MerchantID != uuid.Nil {
			builder = builder.SetMerchantID(role.MerchantID)
		}
		if role.StoreID != uuid.Nil {
			builder = builder.SetStoreID(role.StoreID)
		}
		builders = append(builders, builder)
	}

	_, err = repo.Client.RoleMenu.CreateBulk(builders...).Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to bulk create role menus by role/paths: %w", err)
	}
	return
}

func (repo *RoleMenuRepository) DeletesByRoleID(ctx context.Context, roleID uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.DeletesByRoleID")
	defer func() { util.SpanErrFinish(span, err) }()

	_, err = repo.Client.RoleMenu.Delete().Where(rolemenu.RoleID(roleID)).Exec(ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete role menus by role: %w", err)
	}
	return
}

func (repo *RoleMenuRepository) Deletes(ctx context.Context, ids []uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.Deletes")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(ids) == 0 {
		return nil
	}

	_, err = repo.Client.RoleMenu.Delete().Where(rolemenu.IDIn(ids...)).Exec(ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete role menus: %w", err)
	}
	return
}

func (repo *RoleMenuRepository) GetRoleMenus(ctx context.Context, pager *upagination.Pagination, filter *domain.RoleMenuListFilter, orderBys ...domain.RoleMenuListOrderBy) (list []*domain.RoleMenu, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.GetRoleMenus")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count role menus: %w", err)
		return
	}

	items, err := query.Order(repo.orderBy(orderBys...)...).Offset(pager.Offset()).Limit(pager.Size).All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query role menus: %w", err)
		return
	}

	list = lo.Map(items, func(erm *ent.RoleMenu, _ int) *domain.RoleMenu { return convertRoleMenuToDomain(erm) })
	return
}

func (repo *RoleMenuRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) (list []*domain.RoleMenu, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RoleMenuRepository.GetByRoleID")
	defer func() { util.SpanErrFinish(span, err) }()

	items, err := repo.Client.RoleMenu.Query().Where(rolemenu.RoleID(roleID)).All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query role menus by role id: %w", err)
		return
	}

	list = lo.Map(items, func(erm *ent.RoleMenu, _ int) *domain.RoleMenu { return convertRoleMenuToDomain(erm) })
	return
}

func (repo *RoleMenuRepository) buildFilterQuery(filter *domain.RoleMenuListFilter) *ent.RoleMenuQuery {
	query := repo.Client.RoleMenu.Query()
	if filter == nil {
		return query
	}
	if filter.RoleType != "" {
		query = query.Where(rolemenu.RoleTypeEQ(filter.RoleType))
	}
	if filter.RoleID != uuid.Nil {
		query = query.Where(rolemenu.RoleID(filter.RoleID))
	}
	if filter.Path != "" {
		query = query.Where(rolemenu.PathEQ(filter.Path))
	}
	if filter.MerchantID != uuid.Nil {
		query = query.Where(rolemenu.MerchantID(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(rolemenu.StoreID(filter.StoreID))
	}
	return query
}

func (repo *RoleMenuRepository) orderBy(orderBys ...domain.RoleMenuListOrderBy) []rolemenu.OrderOption {
	var opts []rolemenu.OrderOption
	for _, ob := range orderBys {
		rule := lo.TernaryF(ob.Desc, sql.OrderDesc, sql.OrderAsc)
		switch ob.OrderBy {
		case domain.RoleMenuListOrderByID:
			opts = append(opts, rolemenu.ByID(rule))
		case domain.RoleMenuListOrderByCreatedAt:
			opts = append(opts, rolemenu.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, rolemenu.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertRoleMenuToDomain(erm *ent.RoleMenu) *domain.RoleMenu {
	if erm == nil {
		return nil
	}
	return &domain.RoleMenu{
		ID:         erm.ID,
		RoleType:   erm.RoleType,
		RoleID:     erm.RoleID,
		Path:       erm.Path,
		MerchantID: erm.MerchantID,
		StoreID:    erm.StoreID,
		CreatedAt:  erm.CreatedAt,
		UpdatedAt:  erm.UpdatedAt,
	}
}
