package department

import (
	"context"
	"errors"
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

func (interactor *DepartmentInteractor) CreateDepartment(ctx context.Context, params *domain.CreateDepartmentParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.CreateDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	department := &domain.Department{
		ID:             uuid.New(),
		Name:           params.Name,
		Code:           params.Code,
		DepartmentType: params.DepartmentType,
		Enable:         params.Enable,
		MerchantID:     params.MerchantID,
		StoreID:        params.StoreID,
	}

	if err = interactor.checkExists(ctx, department); err != nil {
		return err
	}

	if err = interactor.DS.DepartmentRepo().Create(ctx, department); err != nil {
		return fmt.Errorf("failed to create department: %w", err)
	}

	return
}

func (interactor *DepartmentInteractor) UpdateDepartment(ctx context.Context, params *domain.UpdateDepartmentParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.UpdateDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	if params == nil {
		return fmt.Errorf("params is nil")
	}

	old, err := interactor.DS.DepartmentRepo().FindByID(ctx, params.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrDepartmentNotExists)
		}
		return fmt.Errorf("failed to fetch department: %w", err)
	}

	department := &domain.Department{
		ID:             old.ID,
		Name:           params.Name,
		Code:           old.Code,
		DepartmentType: old.DepartmentType,
		Enable:         params.Enable,
		MerchantID:     old.MerchantID,
		StoreID:        old.StoreID,
	}

	if err = interactor.checkExists(ctx, department); err != nil {
		return err
	}

	if err = interactor.DS.DepartmentRepo().Update(ctx, department); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrDepartmentNotExists)
		}
		return fmt.Errorf("failed to update department: %w", err)
	}

	return
}

func (interactor *DepartmentInteractor) DeleteDepartment(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.DeleteDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		hasUsers, err := ds.DepartmentRepo().CheckUserInDepartment(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to check user in department: %w", err)
		}
		if hasUsers {
			return domain.ParamsError(domain.ErrDepartmentHasUsersCannotDelete)
		}
		if err = ds.DepartmentRepo().Delete(ctx, id); err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrDepartmentNotExists)
			}
			return fmt.Errorf("failed to delete department: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (interactor *DepartmentInteractor) GetDepartment(ctx context.Context, id uuid.UUID) (department *domain.Department, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.GetDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	department, err = interactor.DS.DepartmentRepo().FindByID(ctx, id)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrDepartmentNotExists)
			return
		}
		err = fmt.Errorf("failed to fetch department: %w", err)
		return
	}

	return
}

func (interactor *DepartmentInteractor) GetDepartments(ctx context.Context,
	pager *upagination.Pagination,
	filter *domain.DepartmentListFilter,
	orderBys ...domain.DepartmentListOrderBy,
) (departments []*domain.Department, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.GetDepartments")
	defer func() { util.SpanErrFinish(span, err) }()

	if pager == nil {
		err = domain.ParamsError(errors.New("pager is required"))
		return
	}
	if filter == nil {
		err = domain.ParamsError(errors.New("filter is required"))
		return
	}

	departments, total, err = interactor.DS.DepartmentRepo().GetDepartments(ctx, pager, filter, orderBys...)
	if err != nil {
		err = fmt.Errorf("failed to get departments: %w", err)
		return
	}

	return
}

func (interactor *DepartmentInteractor) checkExists(ctx context.Context, department *domain.Department) (err error) {
	if department == nil {
		return fmt.Errorf("department is nil")
	}

	exists, existsErr := interactor.DS.DepartmentRepo().Exists(ctx, domain.DepartmentExistsParams{
		Name:       department.Name,
		ExcludeID:  department.ID,
		MerchantID: department.MerchantID,
		StoreID:    department.StoreID,
	})
	if existsErr != nil {
		return fmt.Errorf("failed to check department exists: %w", existsErr)
	}
	if exists {
		return domain.ParamsError(domain.ErrDepartmentNameExists)
	}

	return nil
}

func (interactor *DepartmentInteractor) SimpleUpdate(ctx context.Context, updateField domain.DepartmentSimpleUpdateField, params domain.DepartmentSimpleUpdateParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.SimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		oldDept, err := ds.DepartmentRepo().FindByID(ctx, params.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrDepartmentNotExists)
			}
			return fmt.Errorf("failed to fetch department: %w", err)
		}

		switch updateField {
		case domain.DepartmentSimpleUpdateFieldEnable:
			if oldDept.Enable == params.Enable {
				return nil
			}
			if !params.Enable {
				// If disabling, ensure no users exist in department
				hasUsers, err := ds.DepartmentRepo().CheckUserInDepartment(ctx, params.ID)
				if err != nil {
					return fmt.Errorf("failed to check users in department: %w", err)
				}
				if hasUsers {
					return domain.ParamsError(domain.ErrDepartmentHasUsersCannotDelete)
				}
			}
			oldDept.Enable = params.Enable
		default:
			return domain.ParamsError(fmt.Errorf("unsupported simple update field: %s", updateField))
		}

		if err = ds.DepartmentRepo().Update(ctx, oldDept); err != nil {
			return err
		}
		return nil
	})

	return
}
