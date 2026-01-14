package department

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// 确保实现了领域层定义的 DepartmentInteractor 接口
var _ domain.DepartmentInteractor = (*DepartmentInteractor)(nil)

type DepartmentInteractor struct {
	DS domain.DataStore
}

func NewDepartmentInteractor(ds domain.DataStore) *DepartmentInteractor {
	return &DepartmentInteractor{DS: ds}
}

func (interactor *DepartmentInteractor) CreateDepartment(ctx context.Context, params *domain.CreateDepartmentParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.CreateDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	department := &domain.Department{
		ID:             uuid.New(),
		Name:           params.Name,
		Code:           params.Code,
		DepartmentType: params.DepartmentType,
		Enabled:        params.Enabled,
		MerchantID:     params.MerchantID,
		StoreID:        params.StoreID,
	}

	if err = verifyDepartmentOwnership(user, department); err != nil {
		return err
	}

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.DepartmentRepo().Exists(ctx, domain.DepartmentExistsParams{
			Name:       department.Name,
			ExcludeID:  department.ID,
			MerchantID: department.MerchantID,
			StoreID:    department.StoreID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrDepartmentNameExists
		}
		err = ds.DepartmentRepo().Create(ctx, department)
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

func (interactor *DepartmentInteractor) UpdateDepartment(ctx context.Context, params *domain.UpdateDepartmentParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.UpdateDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		if !params.Enabled {
			hasUsers, err := ds.DepartmentRepo().CheckUserInDepartment(ctx, params.ID)
			if err != nil {
				return err
			}
			if hasUsers {
				return domain.ErrDepartmentHasUsersCannotDisable
			}
		}

		old, err := ds.DepartmentRepo().FindByID(ctx, params.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrDepartmentNotExists
			}
			return err
		}
		if err = verifyDepartmentOwnership(user, old); err != nil {
			return err
		}
		department := &domain.Department{
			ID:             old.ID,
			Name:           params.Name,
			Code:           old.Code,
			DepartmentType: old.DepartmentType,
			Enabled:        params.Enabled,
			MerchantID:     old.MerchantID,
			StoreID:        old.StoreID,
		}
		exists, err := ds.DepartmentRepo().Exists(ctx, domain.DepartmentExistsParams{
			Name:       department.Name,
			ExcludeID:  department.ID,
			MerchantID: department.MerchantID,
			StoreID:    department.StoreID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrDepartmentNameExists
		}
		err = ds.DepartmentRepo().Update(ctx, department)
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

func (interactor *DepartmentInteractor) DeleteDepartment(ctx context.Context, id uuid.UUID, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.DeleteDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		hasUsers, err := ds.DepartmentRepo().CheckUserInDepartment(ctx, id)
		if err != nil {
			return err
		}
		if hasUsers {
			return domain.ErrDepartmentHasUsersCannotDelete
		}
		dept, err := ds.DepartmentRepo().FindByID(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrDepartmentNotExists
			}
			return err
		}
		if err = verifyDepartmentOwnership(user, dept); err != nil {
			return err
		}
		err = ds.DepartmentRepo().Delete(ctx, id)
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

func (interactor *DepartmentInteractor) GetDepartment(ctx context.Context, id uuid.UUID, user domain.User) (department *domain.Department, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.GetDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	department, err = interactor.DS.DepartmentRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, domain.ErrDepartmentNotExists
		}
		return nil, err
	}
	if err = verifyDepartmentOwnership(user, department); err != nil {
		return nil, err
	}

	return department, nil
}

func (interactor *DepartmentInteractor) GetDepartments(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.DepartmentListFilter,
	orderBys ...domain.DepartmentListOrderBy,
) (departments []*domain.Department, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.GetDepartments")
	defer func() { util.SpanErrFinish(span, err) }()

	departments, total, err = interactor.DS.DepartmentRepo().GetDepartments(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get departments: %w", err)
		return
	}

	return departments, total, nil
}

func (interactor *DepartmentInteractor) SimpleUpdate(ctx context.Context, updateField domain.DepartmentSimpleUpdateField, params domain.DepartmentSimpleUpdateParams, user domain.User) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.SimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldDept, err := ds.DepartmentRepo().FindByID(ctx, params.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrDepartmentNotExists
			}
			return err
		}

		if err = verifyDepartmentOwnership(user, oldDept); err != nil {
			return err
		}

		switch updateField {
		case domain.DepartmentSimpleUpdateFieldEnabled:
			if oldDept.Enabled == params.Enabled {
				return nil
			}
			if !params.Enabled {
				// If disabling, ensure no users exist in department
				hasUsers, err := ds.DepartmentRepo().CheckUserInDepartment(ctx, params.ID)
				if err != nil {
					return fmt.Errorf("failed to check users in department: %w", err)
				}
				if hasUsers {
					return domain.ErrDepartmentHasUsersCannotDelete
				}
			}
			oldDept.Enabled = params.Enabled
		default:
			return fmt.Errorf("unsupported simple update field: %s", updateField)
		}
		err = ds.DepartmentRepo().Update(ctx, oldDept)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func verifyDepartmentOwnership(user domain.User, dept *domain.Department) error {
	switch user.GetUserType() {
	case domain.UserTypeAdmin:
	case domain.UserTypeBackend:
		if !domain.VerifyOwnerMerchant(user, dept.MerchantID) {
			return domain.ErrDepartmentNotExists
		}
	case domain.UserTypeStore:
		if !domain.VerifyOwnerShip(user, dept.MerchantID, dept.StoreID) {
			return domain.ErrDepartmentNotExists
		}
	}
	return nil
}
