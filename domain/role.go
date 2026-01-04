package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrRoleNotExists  = errors.New("角色不存在")
	ErrRoleNameExists = errors.New("角色名称已存在")
	ErrRoleCodeExists = errors.New("角色编码已存在")
)

// RoleRepository 角色仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/role_repository.go -package=mock . RoleRepository
type RoleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params RoleExistsParams) (bool, error)
	GetRoles(ctx context.Context, pager *upagination.Pagination, filter *RoleListFilter, orderBys ...RoleListOrderBy) ([]*Role, int, error)
}

// RoleInteractor 角色用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/role_interactor.go -package=mock . RoleInteractor
type RoleInteractor interface {
	CreateRole(ctx context.Context, params *CreateRoleParams) error
	UpdateRole(ctx context.Context, params *UpdateRoleParams) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	GetRole(ctx context.Context, id uuid.UUID) (*Role, error)
	GetRoles(ctx context.Context, pager *upagination.Pagination, filter *RoleListFilter, orderBys ...RoleListOrderBy) ([]*Role, int, error)
}

type RoleType string

const (
	RoleTypeAdmin   RoleType = "admin"
	RoleTypeBackend RoleType = "backend"
	RoleTypeStore   RoleType = "store"
)

func (RoleType) Values() []string {
	return []string{string(RoleTypeAdmin), string(RoleTypeBackend), string(RoleTypeStore)}
}

// RoleListOrderByType defines the allowed ordering columns for roles.
type RoleListOrderByType int

const (
	_ RoleListOrderByType = iota
	RoleListOrderByID
	RoleListOrderByCreatedAt
)

type RoleListOrderBy struct {
	OrderBy RoleListOrderByType
	Desc    bool
}

func NewRoleListOrderByID(desc bool) RoleListOrderBy {
	return RoleListOrderBy{OrderBy: RoleListOrderByID, Desc: desc}
}

func NewRoleListOrderByCreatedAt(desc bool) RoleListOrderBy {
	return RoleListOrderBy{OrderBy: RoleListOrderByCreatedAt, Desc: desc}
}

type Role struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	RoleType  RoleType  `json:"role_type"`
	Enable    bool      `json:"enable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleListFilter struct {
	Name     string   `json:"name"`
	RoleType RoleType `json:"role_type"`
	Enable   *bool    `json:"enable"`
}

type CreateRoleParams struct {
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	RoleType RoleType `json:"role_type"`
	Enable   bool     `json:"enable"`
}

type UpdateRoleParams struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	RoleType RoleType  `json:"role_type"`
	Enable   bool      `json:"enable"`
}

type RoleExistsParams struct {
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	ExcludeID uuid.UUID `json:"exclude_id"`
}
