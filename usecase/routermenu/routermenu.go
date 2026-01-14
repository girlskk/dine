package routermenu

import (
	"context"
	"errors"
	"fmt"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RouterMenuInteractor = (*RouterMenuInteractor)(nil)

type RouterMenuInteractor struct {
	DS domain.DataStore
}

func NewRouterMenuInteractor(ds domain.DataStore) *RouterMenuInteractor {
	return &RouterMenuInteractor{DS: ds}
}

func (interactor *RouterMenuInteractor) GetRouterMenus(ctx context.Context,
	filter *domain.RouterMenuListFilter,
	orderBys ...domain.RouterMenuListOrderBy,
) (menus []*domain.RouterMenu, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RouterMenuInteractor.GetRouterMenus")
	defer func() { util.SpanErrFinish(span, err) }()

	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
		return
	}

	menus, total, err = interactor.DS.RouterMenuRepo().GetRouterMenus(ctx, filter, orderBys...)
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

	exists, existsErr := interactor.DS.RouterMenuRepo().Exists(ctx, domain.RouterMenuExistsParams{
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
