package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RoleInteractor = (*RoleInteractor)(nil)

type RoleInteractor struct {
	ds domain.DataStore
}

func NewRoleInteractor(ds domain.DataStore) *RoleInteractor {
	return &RoleInteractor{ds: ds}
}

func (interactor *RoleInteractor) CreateRole(ctx context.Context, params *domain.CreateRoleParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.CreateRole")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	role := &domain.Role{
		ID:         uuid.New(),
		Name:       params.Name,
		Code:       params.Code,
		RoleType:   params.RoleType,
		DataScope:  params.DataScope,
		Enable:     params.Enable,
		MerchantID: params.MerchantID,
		StoreID:    params.StoreID,
	}

	if err = interactor.checkExists(ctx, role); err != nil {
		return err
	}

	if err = interactor.ds.RoleRepo().Create(ctx, role); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return
}

func (interactor *RoleInteractor) UpdateRole(ctx context.Context, params *domain.UpdateRoleParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.UpdateRole")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	old, err := interactor.ds.RoleRepo().FindByID(ctx, params.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRoleNotExists)
		}
		return fmt.Errorf("failed to fetch role: %w", err)
	}

	role := &domain.Role{
		ID:         old.ID,
		Name:       params.Name,
		Code:       params.Code,
		RoleType:   params.RoleType,
		DataScope:  params.DataScope,
		Enable:     params.Enable,
		MerchantID: params.MerchantID,
		StoreID:    params.StoreID,
	}

	if err = interactor.checkExists(ctx, role); err != nil {
		return err
	}

	if err = interactor.ds.RoleRepo().Update(ctx, role); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRoleNotExists)
		}
		return fmt.Errorf("failed to update role: %w", err)
	}

	return
}

func (interactor *RoleInteractor) DeleteRole(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.DeleteRole")
	defer func() { util.SpanErrFinish(span, err) }()

	if err = interactor.ds.RoleRepo().Delete(ctx, id); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrRoleNotExists)
		}
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return
}

func (interactor *RoleInteractor) GetRole(ctx context.Context, id uuid.UUID) (role *domain.Role, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.GetRole")
	defer func() { util.SpanErrFinish(span, err) }()

	role, err = interactor.ds.RoleRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrRoleNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch role: %w", err)
		return
	}

	return
}

func (interactor *RoleInteractor) GetRoles(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.RoleListFilter,
	orderBys ...domain.RoleListOrderBy,
) (roles []*domain.Role, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.GetRoles")
	defer func() { util.SpanErrFinish(span, err) }()

	if pager == nil {
		err = domain.ParamsError(errors.New("pager is required"))
		return
	}
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
		return
	}

	roles, total, err = interactor.ds.RoleRepo().GetRoles(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get roles: %w", err)
		return
	}

	return
}

func (interactor *RoleInteractor) checkExists(ctx context.Context, role *domain.Role) (err error) {
	if role == nil {
		return fmt.Errorf("role is nil")
	}

	exists, existsErr := interactor.ds.RoleRepo().Exists(ctx, domain.RoleExistsParams{
		Name:       role.Name,
		Code:       role.Code,
		MerchantID: role.MerchantID,
		StoreID:    role.StoreID,
		ExcludeID:  role.ID,
	})
	if existsErr != nil {
		return fmt.Errorf("failed to check role exists: %w", existsErr)
	}
	if exists {
		// Prefer name conflict message first; code uniqueness uses same check
		return domain.ParamsError(domain.ErrRoleNameExists)
	}

	return nil
}
