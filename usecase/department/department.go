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
	ds domain.DataStore
}

func NewDepartmentInteractor(ds domain.DataStore) *DepartmentInteractor {
	return &DepartmentInteractor{ds: ds}
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

	if err = interactor.ds.DepartmentRepo().Create(ctx, department); err != nil {
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

	old, err := interactor.ds.DepartmentRepo().FindByID(ctx, params.ID)
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

	if err = interactor.ds.DepartmentRepo().Update(ctx, department); err != nil {
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

	if err = interactor.ds.DepartmentRepo().Delete(ctx, id); err != nil {
		if domain.IsNotFound(err) {
			return domain.ParamsError(domain.ErrDepartmentNotExists)
		}
		return fmt.Errorf("failed to delete department: %w", err)
	}

	return
}

func (interactor *DepartmentInteractor) GetDepartment(ctx context.Context, id uuid.UUID) (department *domain.Department, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "DepartmentInteractor.GetDepartment")
	defer func() { util.SpanErrFinish(span, err) }()

	department, err = interactor.ds.DepartmentRepo().FindByID(ctx, id)
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

	departments, total, err = interactor.ds.DepartmentRepo().GetDepartments(ctx, pager, filter, orderBys...)
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

	exists, existsErr := interactor.ds.DepartmentRepo().Exists(ctx, domain.DepartmentExistsParams{
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
