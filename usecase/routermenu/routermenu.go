package routermenu

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RouterMenuInteractor = (*RouterMenuInteractor)(nil)

type RouterMenuInteractor struct {
	ds domain.DataStore
}

func NewRouterMenuInteractor(ds domain.DataStore) *RouterMenuInteractor {
	return &RouterMenuInteractor{ds: ds}
}

func (interactor *RouterMenuInteractor) CreateRouterMenu(ctx context.Context, params *domain.CreateRouterMenuParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RouterMenuInteractor.CreateRouterMenu")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	menu := &domain.RouterMenu{
		ID:        uuid.New(),
		UserType:  params.UserType,
		ParentID:  params.ParentID,
		Name:      params.Name,
		Path:      params.Path,
		Component: params.Component,
		Icon:      params.Icon,
		Sort:      params.Sort,
		Enabled:   params.Enabled,
		Layer:     1,
	}

	if err = interactor.checkExists(ctx, menu); err != nil {
		return err
	}

	if params.ParentID != uuid.Nil {
		parent, err := interactor.ds.RouterMenuRepo().FindByID(ctx, params.ParentID)
		if err != nil {
			return fmt.Errorf("failed to fetch parent router menu: %w", err)
		}
		menu.Layer = parent.Layer + 1
	}
	if err = interactor.ds.RouterMenuRepo().Create(ctx, menu); err != nil {
		return fmt.Errorf("failed to create router menu: %w", err)
	}

	return
}

func (interactor *RouterMenuInteractor) UpdateRouterMenu(ctx context.Context, params *domain.UpdateRouterMenuParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RouterMenuInteractor.UpdateRouterMenu")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	old, err := interactor.ds.RouterMenuRepo().FindByID(ctx, params.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRouterMenuNotExists)
		}
		return fmt.Errorf("failed to fetch router menu: %w", err)
	}

	menu := &domain.RouterMenu{
		ID:        old.ID,
		UserType:  old.UserType,
		ParentID:  params.ParentID,
		Name:      params.Name,
		Path:      params.Path,
		Component: params.Component,
		Icon:      params.Icon,
		Sort:      params.Sort,
		Enabled:   params.Enabled,
		Layer:     old.Layer,
	}

	if err = interactor.checkExists(ctx, menu); err != nil {
		return err
	}

	if params.ParentID != uuid.Nil {
		parent, err := interactor.ds.RouterMenuRepo().FindByID(ctx, params.ParentID)
		if err != nil {
			return fmt.Errorf("failed to fetch parent router menu: %w", err)
		}
		menu.Layer = parent.Layer + 1
	}

	if err = interactor.ds.RouterMenuRepo().Update(ctx, menu); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRouterMenuNotExists)
		}
		return fmt.Errorf("failed to update router menu: %w", err)
	}

	return
}

func (interactor *RouterMenuInteractor) DeleteRouterMenu(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RouterMenuInteractor.DeleteRouterMenu")
	defer func() { util.SpanErrFinish(span, err) }()

	if err = interactor.ds.RouterMenuRepo().Delete(ctx, id); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRouterMenuNotExists)
		}
		return fmt.Errorf("failed to delete router menu: %w", err)
	}

	return
}

func (interactor *RouterMenuInteractor) GetRouterMenu(ctx context.Context, id uuid.UUID) (menu *domain.RouterMenu, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RouterMenuInteractor.GetRouterMenu")
	defer func() { util.SpanErrFinish(span, err) }()

	menu, err = interactor.ds.RouterMenuRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrRouterMenuNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch router menu: %w", err)
		return
	}

	return
}

func (interactor *RouterMenuInteractor) GetRouterMenus(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.RouterMenuListFilter,
	orderBys ...domain.RouterMenuListOrderBy,
) (menus []*domain.RouterMenu, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RouterMenuInteractor.GetRouterMenus")
	defer func() { util.SpanErrFinish(span, err) }()

	if pager == nil {
		err = domain.ParamsError(errors.New("pager is required"))
		return
	}
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
		return
	}

	menus, total, err = interactor.ds.RouterMenuRepo().GetRouterMenus(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get router menus: %w", err)
		return
	}

	return
}

func (interactor *RouterMenuInteractor) checkExists(ctx context.Context, menu *domain.RouterMenu) (err error) {
	if menu == nil {
		return fmt.Errorf("router menu is nil")
	}

	exists, existsErr := interactor.ds.RouterMenuRepo().Exists(ctx, domain.RouterMenuExistsParams{
		ParentID:  menu.ParentID,
		Name:      menu.Name,
		ExcludeID: menu.ID,
		UserType:  menu.UserType,
	})
	if existsErr != nil {
		return fmt.Errorf("failed to check router menu exists: %w", existsErr)
	}
	if exists {
		return domain.ParamsError(domain.ErrRouterMenuNameExists)
	}

	return nil
}
