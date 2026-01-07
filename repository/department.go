package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/department"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.DepartmentRepository = (*DepartmentRepository)(nil)

// DepartmentRepository implements Department CRUD and pagination.
type DepartmentRepository struct {
	Client *ent.Client
}

func NewDepartmentRepository(client *ent.Client) *DepartmentRepository {
	return &DepartmentRepository{Client: client}
}

func (repo *DepartmentRepository) FindByID(ctx context.Context, id uuid.UUID) (dept *domain.Department, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DepartmentRepository.FindByID")
	defer func() { util.SpanErrFinish(span, err) }()

	ed, err := repo.Client.Department.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrDepartmentNotExists)
			return
		}
		return
	}
	dept = convertDepartmentToDomain(ed)
	return
}

func (repo *DepartmentRepository) Create(ctx context.Context, dept *domain.Department) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DepartmentRepository.Create")
	defer func() { util.SpanErrFinish(span, err) }()
	if dept == nil {
		return fmt.Errorf("department is nil")
	}

	builder := repo.Client.Department.Create().
		SetID(dept.ID).
		SetName(dept.Name).
		SetCode(dept.Code).
		SetDepartmentType(dept.DepartmentType).
		SetEnable(dept.Enable)
	if dept.MerchantID != uuid.Nil {
		builder = builder.SetMerchantID(dept.MerchantID)
	}
	if dept.StoreID != uuid.Nil {
		builder = builder.SetStoreID(dept.StoreID)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create department: %w", err)
		return
	}

	dept.ID = created.ID
	dept.CreatedAt = created.CreatedAt
	dept.UpdatedAt = created.UpdatedAt
	return
}

func (repo *DepartmentRepository) Update(ctx context.Context, dept *domain.Department) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DepartmentRepository.Update")
	defer func() { util.SpanErrFinish(span, err) }()
	if dept == nil {
		return fmt.Errorf("department is nil")
	}

	builder := repo.Client.Department.UpdateOneID(dept.ID).
		SetName(dept.Name).
		SetEnable(dept.Enable)

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(domain.ErrDepartmentNotExists)
			return
		}
		err = fmt.Errorf("failed to update department: %w", err)
		return
	}

	dept.UpdatedAt = updated.UpdatedAt
	return
}

func (repo *DepartmentRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DepartmentRepository.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	err = repo.Client.Department.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
			return
		}
		err = fmt.Errorf("failed to delete department: %w", err)
		return
	}
	return nil
}

func (repo *DepartmentRepository) GetDepartments(ctx context.Context, pager *upagination.Pagination, filter *domain.DepartmentListFilter, orderBys ...domain.DepartmentListOrderBy) (departments []*domain.Department, total int, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DepartmentRepository.GetDepartments")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.buildFilterQuery(filter)

	total, err = query.Clone().Count(ctx)
	if err != nil {
		err = fmt.Errorf("failed to count department: %w", err)
		return
	}

	list, err := query.
		Order(repo.orderBy(orderBys...)...).
		Offset(pager.Offset()).
		Limit(pager.Size).
		All(ctx)
	if err != nil {
		err = fmt.Errorf("failed to query department: %w", err)
		return
	}

	departments = lo.Map(list, func(item *ent.Department, _ int) *domain.Department {
		return convertDepartmentToDomain(item)
	})
	return
}

func (repo *DepartmentRepository) Exists(ctx context.Context, params domain.DepartmentExistsParams) (exists bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "DepartmentRepository.Exists")
	defer func() { util.SpanErrFinish(span, err) }()

	query := repo.Client.Department.Query().
		Where(department.Name(params.Name))
	if params.MerchantID != uuid.Nil {
		query = query.Where(department.MerchantID(params.MerchantID))
	}
	if params.StoreID != uuid.Nil {
		query = query.Where(department.StoreID(params.StoreID))
	}
	if params.ExcludeID != uuid.Nil {
		query = query.Where(department.IDNEQ(params.ExcludeID))
	}

	exists, err = query.Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check department existence: %w", err)
	}
	return
}

func (repo *DepartmentRepository) buildFilterQuery(filter *domain.DepartmentListFilter) *ent.DepartmentQuery {
	query := repo.Client.Department.Query()
	if filter == nil {
		return query
	}

	if filter.MerchantID != uuid.Nil {
		query = query.Where(department.MerchantID(filter.MerchantID))
	}
	if filter.StoreID != uuid.Nil {
		query = query.Where(department.StoreID(filter.StoreID))
	}
	if filter.DepartmentType != "" {
		query = query.Where(department.DepartmentTypeEQ(filter.DepartmentType))
	}
	if filter.Enable != nil {
		query = query.Where(department.EnableEQ(*filter.Enable))
	}
	if filter.Name != "" {
		query = query.Where(department.NameContains(filter.Name))
	}
	if filter.Code != "" {
		query = query.Where(department.CodeContains(filter.Code))
	}

	return query
}

func (repo *DepartmentRepository) orderBy(orderBys ...domain.DepartmentListOrderBy) []department.OrderOption {
	var opts []department.OrderOption
	for _, orderBy := range orderBys {
		rule := lo.TernaryF(orderBy.Desc, sql.OrderDesc, sql.OrderAsc)
		switch orderBy.OrderBy {
		case domain.DepartmentListOrderByID:
			opts = append(opts, department.ByID(rule))
		case domain.DepartmentListOrderByCreatedAt:
			opts = append(opts, department.ByCreatedAt(rule))
		}
	}
	if len(opts) == 0 {
		opts = append(opts, department.ByCreatedAt(sql.OrderDesc()))
	}
	return opts
}

func convertDepartmentToDomain(ed *ent.Department) *domain.Department {
	if ed == nil {
		return nil
	}
	return &domain.Department{
		ID:             ed.ID,
		Name:           ed.Name,
		Code:           ed.Code,
		DepartmentType: ed.DepartmentType,
		Enable:         ed.Enable,
		MerchantID:     ed.MerchantID,
		StoreID:        ed.StoreID,
		CreatedAt:      ed.CreatedAt,
		UpdatedAt:      ed.UpdatedAt,
	}
}
