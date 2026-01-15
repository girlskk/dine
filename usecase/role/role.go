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
	DS domain.DataStore
}

func NewRoleInteractor(ds domain.DataStore) *RoleInteractor {
	return &RoleInteractor{DS: ds}
}

func (interactor *RoleInteractor) CreateRole(ctx context.Context, params *domain.CreateRoleParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.CreateRole")
	defer func() { util.SpanErrFinish(span, err) }()
	role := &domain.Role{
		ID:            uuid.New(),
		Name:          params.Name,
		Code:          params.Code,
		RoleType:      params.RoleType,
		DataScope:     params.DataScope,
		Enabled:       params.Enabled,
		MerchantID:    params.MerchantID,
		StoreID:       params.StoreID,
		LoginChannels: params.LoginChannels,
	}

	if err = verifyRoleOwnership(user, role); err != nil {
		return err
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.RoleRepo().Exists(ctx, domain.RoleExistsParams{
			Name:       role.Name,
			MerchantID: role.MerchantID,
			StoreID:    role.StoreID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrRoleNameExists
		}
		err = ds.RoleRepo().Create(ctx, role)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (interactor *RoleInteractor) UpdateRole(ctx context.Context, params *domain.UpdateRoleParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.UpdateRole")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		old, err := ds.RoleRepo().FindByID(ctx, params.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrRoleNotExists
			}
			return err
		}

		if err = verifyRoleOwnership(user, old); err != nil {
			return err
		}

		userRoles, err := ds.UserRoleRepo().GetByRoleIDs(ctx, domain.UserType(old.RoleType), params.ID)
		if err != nil {
			return err
		}
		if len(userRoles) > 0 {
			return domain.ErrRoleAssignedCannotDisable
		}

		role := &domain.Role{
			ID:            old.ID,
			Name:          params.Name,
			Code:          old.Code,
			RoleType:      params.RoleType,
			DataScope:     params.DataScope,
			Enabled:       params.Enabled,
			MerchantID:    params.MerchantID,
			StoreID:       params.StoreID,
			LoginChannels: params.LoginChannels,
		}
		exists, err := ds.RoleRepo().Exists(ctx, domain.RoleExistsParams{
			Name:       role.Name,
			MerchantID: role.MerchantID,
			StoreID:    role.StoreID,
			ExcludeID:  role.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrRoleNameExists
		}
		err = ds.RoleRepo().Update(ctx, role)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (interactor *RoleInteractor) SimpleUpdate(ctx context.Context, updateField domain.RoleSimpleUpdateField, params domain.RoleSimpleUpdateParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.SimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldRole, err := ds.RoleRepo().FindByID(ctx, params.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrRoleNotExists
			}
			return err
		}

		if err = verifyRoleOwnership(user, oldRole); err != nil {
			return err
		}

		switch updateField {
		case domain.RoleSimpleUpdateFieldEnabled:
			if oldRole.Enabled == params.Enabled {
				return nil
			}
			if !params.Enabled {
				userRoles, err := ds.UserRoleRepo().GetByRoleIDs(ctx, domain.UserType(oldRole.RoleType), params.ID)
				if err != nil {
					return err
				}
				if len(userRoles) > 0 {
					return domain.ErrRoleAssignedCannotDisable
				}
			}
			oldRole.Enabled = params.Enabled
		default:
			return fmt.Errorf("unsupported simple update field: %s", updateField)
		}
		err = ds.RoleRepo().Update(ctx, oldRole)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (interactor *RoleInteractor) DeleteRole(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.DeleteRole")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldRole, err := ds.RoleRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrRoleNotExists
			}
			return err
		}

		if err = verifyRoleOwnership(user, oldRole); err != nil {
			return err
		}

		userRoles, err := ds.UserRoleRepo().GetByRoleIDs(ctx, domain.UserType(oldRole.RoleType), id)
		if err != nil {
			return err
		}
		if len(userRoles) > 0 {
			return domain.ErrRoleAssignedCannotDelete
		}
		err = interactor.DS.RoleRepo().Delete(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrRoleNotExists
			}
			return fmt.Errorf("failed to delete role: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return
}

func (interactor *RoleInteractor) GetRole(ctx context.Context, id uuid.UUID, user domain.User) (role *domain.Role, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleInteractor.GetRole")
	defer func() { util.SpanErrFinish(span, err) }()

	role, err = interactor.DS.RoleRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ErrRoleNotExists
		}
		return
	}

	if err = verifyRoleOwnership(user, role); err != nil {
		return nil, err
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

	roles, total, err = interactor.DS.RoleRepo().GetRoles(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get roles: %w", err)
		return
	}

	return
}

func verifyRoleOwnership(user domain.User, role *domain.Role) error {
	switch user.GetUserType() {
	case domain.UserTypeAdmin:
	case domain.UserTypeBackend:
		if !domain.VerifyOwnerMerchant(user, role.MerchantID) {
			return domain.ErrRoleNotExists
		}
	case domain.UserTypeStore:
		if !domain.VerifyOwnerShip(user, role.MerchantID, role.StoreID) {
			return domain.ErrRoleNotExists
		}
	}

	return nil
}
