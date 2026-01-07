package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/routermenu"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RouterMenuRepository = (*RouterMenuRepository)(nil)

// RouterMenuRepository implements RouterMenu CRUD and pagination.
type RouterMenuRepository struct {
	Client *ent.Client
}

func NewRouterMenuRepository(client *ent.Client) *RouterMenuRepository {
	return &RouterMenuRepository{Client: client}
}

func (repo *RouterMenuRepository) FindByID(ctx context.Context, id uuid.UUID) (m *domain.RouterMenu, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RouterMenuRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	em, err := repo.Client.RouterMenu.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrRouterMenuNotExists)
			return
		}
		return
	}
	m = convertRouterMenuToDomain(em)
	return
}

func (repo *RouterMenuRepository) Create(ctx context.Context, m *domain.RouterMenu) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RouterMenuRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if m == nil {
		return fmt.Errorf("router menu is nil")
	}

	builder := repo.Client.RouterMenu.Create().
		SetID(m.ID).
		SetUserType(m.UserType).
		SetParentID(m.ParentID).
		SetName(m.Name).
		SetPath(m.Path).
		SetLayer(m.Layer).
		SetComponent(m.Component).
		SetIcon(m.Icon).
		SetSort(m.Sort).
		SetEnabled(m.Enabled)

	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create router menu: %w", err)
		return
	}

	m.ID = created.ID
	m.CreatedAt = created.CreatedAt
	m.UpdatedAt = created.UpdatedAt
	return
}

func (repo *RouterMenuRepository) Update(ctx context.Context, m *domain.RouterMenu) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RouterMenuRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if m == nil {
		return fmt.Errorf("router menu is nil")
	}

	builder := repo.Client.RouterMenu.UpdateOneID(m.ID).
		SetParentID(m.ParentID).
		SetName(m.Name).
		SetPath(m.Path).
		SetLayer(m.Layer).
		SetComponent(m.Component).
		SetIcon(m.Icon).
		SetSort(m.Sort).
		SetEnabled(m.Enabled)

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrRouterMenuNotExists)
			return
		}
		err = fmt.Errorf("failed to update router menu: %w", err)
		return
	}

	m.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *RouterMenuRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RouterMenuRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.RouterMenu.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete router menu: %w", err)
		return
	}
	return nil
}

func (repo *RouterMenuRepository) GetRouterMenus(ctx context.Context, pager *upagination.Pagination, filter *domain.RouterMenuListFilter, orderBys ...domain.RouterMenuListOrderBy) (menus []*domain.RouterMenu, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RouterMenuRepository.GetRouterMenus")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count router menu: %w", err)
		return
	}

	list, err := query.
		Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query router menu: %w", err)
		return
	}

	menus = lo.Map(list, func(item *ent.RouterMenu, _ int) *domain.RouterMenu {
		return convertRouterMenuToDomain(item)
	})
	return
}

func (repo *RouterMenuRepository) Exists(ctx context.Context, params domain.RouterMenuExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "RouterMenuRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.RouterMenu.Query()
	if params.ParentID != uuid.Nil {
		query = query.Where(routermenu.ParentID(params.ParentID))
	} else {
		query = query.Where(routermenu.ParentIDIsNil())
	}
	if params.Name != "" {
		query = query.Where(routermenu.Name(params.Name))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(routermenu.IDNEQ(params.ExcludeID))
	}
	if params.UserType != "" {
		query = query.Where(routermenu.UserTypeEQ(params.UserType))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check router menu existence: %w", err)
	}
	return
}

func (repo *RouterMenuRepository) buildFilterQuery(filter *domain.RouterMenuListFilter) *ent.RouterMenuQuery {
	query := repo.Client.RouterMenu.Query()
	if filter == nil {
		return query
	}

	if filter.ParentID != uuid.Nil {
		query = query.Where(routermenu.ParentID(filter.ParentID))
	}
	if filter.Name != "" {
		query = query.Where(routermenu.NameContains(filter.Name))
	}
	if filter.Enabled != nil {
		query = query.Where(routermenu.EnabledEQ(*filter.Enabled))
	}
	if filter.UserType != "" {
		query = query.Where(routermenu.UserTypeEQ(filter.UserType))
	}
	if filter.Layer > 0 {
		query = query.Where(routermenu.LayerEQ(filter.Layer))
	}

	return query
}

func (repo *RouterMenuRepository) orderBy(orderBys ...domain.RouterMenuListOrderBy) []routermenu.OrderOption {
	var opts []routermenu.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.RouterMenuListOrderByID:
			opts = append(opts, routermenu.ByID(rule))
		case domain.RouterMenuListOrderByCreatedAt:
			opts = append(opts, routermenu.ByCreatedAt(rule))
		case domain.RouterMenuListOrderBySort:
			opts = append(opts, routermenu.BySort(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, routermenu.BySort(sql.OrderAsc()))
	}
	return opts
}

func convertRouterMenuToDomain(em *ent.RouterMenu) *domain.RouterMenu {
	if em == nil {
		return nil
	}
	return &domain.RouterMenu{
		ID:        em.ID,
		UserType:  em.UserType,
		ParentID:  em.ParentID,
		Name:      em.Name,
		Path:      em.Path,
		Layer:     em.Layer,
		Component: em.Component,
		Icon:      em.Icon,
		Sort:      em.Sort,
		Enabled:   em.Enabled,
		CreatedAt: em.CreatedAt,
		UpdatedAt: em.UpdatedAt,
	}
}
