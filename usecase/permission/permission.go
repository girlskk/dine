package permission

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.PermissionInteractor = (*PermissionInteractor)(nil)

type PermissionInteractor struct {
	ds domain.DataStore
}

func NewPermissionInteractor(ds domain.DataStore) *PermissionInteractor {
	return &PermissionInteractor{ds: ds}
}

func (interactor *PermissionInteractor) CreatePermission(ctx context.Context, params *domain.CreatePermissionParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PermissionInteractor.CreatePermission")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	permission := &domain.Permission{
		ID:       uuid.New(),
		MenuID:   params.MenuID,
		PermCode: params.PermCode,
		Name:     params.Name,
		Method:   params.Method,
		Path:     params.Path,
		Enabled:  params.Enabled,
	}

	if err = interactor.checkExists(ctx, permission); err != nil {
		return err
	}

	if err = interactor.ds.PermissionRepo().Create(ctx, permission); err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return
}

func (interactor *PermissionInteractor) UpdatePermission(ctx context.Context, params *domain.UpdatePermissionParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PermissionInteractor.UpdatePermission")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	old, err := interactor.ds.PermissionRepo().FindByID(ctx, params.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrPermissionNotExists)
		}
		return fmt.Errorf("failed to fetch permission: %w", err)
	}

	permission := &domain.Permission{
		ID:       old.ID,
		MenuID:   params.MenuID,
		PermCode: params.PermCode,
		Name:     params.Name,
		Method:   params.Method,
		Path:     params.Path,
		Enabled:  params.Enabled,
	}

	if err = interactor.checkExists(ctx, permission); err != nil {
		return err
	}

	if err = interactor.ds.PermissionRepo().Update(ctx, permission); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrPermissionNotExists)
		}
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return
}

func (interactor *PermissionInteractor) DeletePermission(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PermissionInteractor.DeletePermission")
	defer func() { util.SpanErrFinish(span, err) }()

	if err = interactor.ds.PermissionRepo().Delete(ctx, id); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrPermissionNotExists)
		}
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return
}

func (interactor *PermissionInteractor) GetPermission(ctx context.Context, id uuid.UUID) (permission *domain.Permission, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PermissionInteractor.GetPermission")
	defer func() { util.SpanErrFinish(span, err) }()

	permission, err = interactor.ds.PermissionRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrPermissionNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch permission: %w", err)
		return
	}

	return
}

func (interactor *PermissionInteractor) GetPermissions(ctx context.Context, pager *upagination.Pagination, filter *domain.PermissionListFilter, orderBys ...domain.PermissionListOrderBy) (permissions []*domain.Permission, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "PermissionInteractor.GetPermissions")
	defer func() { util.SpanErrFinish(span, err) }()

	if pager == nil {
		err = domain.ParamsError(errors.New("pager is required"))
		return
	}
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
		return
	}

	permissions, total, err = interactor.ds.PermissionRepo().GetPermissions(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get permissions: %w", err)
		return
	}

	return
}

func (interactor *PermissionInteractor) checkExists(ctx context.Context, permission *domain.Permission) (err error) {
	if permission == nil {
		return fmt.Errorf("permission is nil")
	}
	if permission.PermCode == "" {
		return domain.ParamsError(errors.New("permission code is required"))
	}
	exists, existsErr := interactor.ds.PermissionRepo().Exists(ctx, domain.PermissionExistsParams{
		PermCode:  permission.PermCode,
		ExcludeID: permission.ID,
	})
	if existsErr != nil {
		return fmt.Errorf("failed to check permission exists: %w", existsErr)
	}
	if exists {
		return domain.ParamsError(domain.ErrPermissionCodeExists)
	}

	exists, existsErr = interactor.ds.PermissionRepo().Exists(ctx, domain.PermissionExistsParams{
		Method:    permission.Method,
		Path:      permission.Path,
		ExcludeID: permission.ID,
	})
	if existsErr != nil {
		return fmt.Errorf("failed to check permission exists: %w", existsErr)
	}
	if exists {
		return domain.ParamsError(domain.ErrPermissionCodeExists)
	}
	return nil
}
