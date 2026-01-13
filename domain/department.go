package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrDepartmentNotExists            = errors.New("部门不存在")
	ErrDepartmentNameExists           = errors.New("部门名称已存在")
	ErrDepartmentCodeExists           = errors.New("部门编码已存在")
	ErrDepartmentHasUsersCannotDelete = errors.New("部门下存在用户，无法删除")
)

// DepartmentRepository 部门仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/department_repository.go -package=mock . DepartmentRepository
type DepartmentRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (department *Department, err error)
	Create(ctx context.Context, department *Department) (err error)
	Update(ctx context.Context, department *Department) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	Exists(ctx context.Context, params DepartmentExistsParams) (exists bool, err error)
	GetDepartments(ctx context.Context, pager *upagination.Pagination, filter *DepartmentListFilter, orderBys ...DepartmentListOrderBy) (departments []*Department, total int, err error)
	CheckUserInDepartment(ctx context.Context, departmentID uuid.UUID) (exists bool, err error)
}

// DepartmentInteractor 部门用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/department_interactor.go -package=mock . DepartmentInteractor
type DepartmentInteractor interface {
	CreateDepartment(ctx context.Context, params *CreateDepartmentParams, user User) (err error)
	UpdateDepartment(ctx context.Context, params *UpdateDepartmentParams, user User) (err error)
	DeleteDepartment(ctx context.Context, id uuid.UUID, user User) (err error)
	GetDepartment(ctx context.Context, id uuid.UUID, user User) (department *Department, err error)
	GetDepartments(ctx context.Context, pager *upagination.Pagination, filter *DepartmentListFilter, orderBys ...DepartmentListOrderBy) (departments []*Department, total int, err error)
	SimpleUpdate(ctx context.Context, updateField DepartmentSimpleUpdateField, params DepartmentSimpleUpdateParams, user User) error
}

type DepartmentType string

const (
	DepartmentAdmin   DepartmentType = "admin"
	DepartmentBackend DepartmentType = "backend"
	DepartmentStore   DepartmentType = "store"
)

func (DepartmentType) Values() []string {
	return []string{string(DepartmentAdmin), string(DepartmentBackend), string(DepartmentStore)}
}

type DepartmentListOrderByType int

const (
	_ DepartmentListOrderByType = iota
	DepartmentListOrderByID
	DepartmentListOrderByCreatedAt
)

type DepartmentListOrderBy struct {
	OrderBy DepartmentListOrderByType
	Desc    bool
}

func NewDepartmentListOrderByID(desc bool) DepartmentListOrderBy {
	return DepartmentListOrderBy{OrderBy: DepartmentListOrderByID, Desc: desc}
}

func NewDepartmentListOrderByCreatedAt(desc bool) DepartmentListOrderBy {
	return DepartmentListOrderBy{OrderBy: DepartmentListOrderByCreatedAt, Desc: desc}
}

type Department struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Code           string         `json:"code"`
	DepartmentType DepartmentType `json:"department_type"`
	Enabled        bool           `json:"enabled"`
	MerchantID     uuid.UUID      `json:"merchant_id"`
	StoreID        uuid.UUID      `json:"store_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type DepartmentListFilter struct {
	Name           string         `json:"name"`
	Code           string         `json:"code"`
	DepartmentType DepartmentType `json:"department_type"`
	Enabled        *bool          `json:"enabled"`
	MerchantID     uuid.UUID      `json:"merchant_id"`
	StoreID        uuid.UUID      `json:"store_id"`
}

type CreateDepartmentParams struct {
	Name           string         `json:"name"`
	Code           string         `json:"code"`
	DepartmentType DepartmentType `json:"department_type"`
	Enabled        bool           `json:"enabled"`
	MerchantID     uuid.UUID      `json:"merchant_id"`
	StoreID        uuid.UUID      `json:"store_id"`
}

type UpdateDepartmentParams struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`            // 模糊搜索
	Code           string         `json:"code"`            // 模糊搜索
	DepartmentType DepartmentType `json:"department_type"` // admin/backend/store
	Enabled        bool           `json:"enabled"`         // 是否启用
	MerchantID     uuid.UUID      `json:"merchant_id"`     // 商户 ID
	StoreID        uuid.UUID      `json:"store_id"`        // 门店 ID
}

type DepartmentExistsParams struct {
	Name       string    `json:"name"`
	ExcludeID  uuid.UUID `json:"exclude_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
}

// DepartmentSimpleUpdateField Simple update types for department
type DepartmentSimpleUpdateField string

const (
	DepartmentSimpleUpdateFieldEnabled DepartmentSimpleUpdateField = "enabled"
)

type DepartmentSimpleUpdateParams struct {
	ID      uuid.UUID `json:"id"`
	Enabled bool      `json:"enabled"`
}
