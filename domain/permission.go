package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrPermissionNotExists  = errors.New("权限不存在")
	ErrPermissionCodeExists = errors.New("权限编码已存在")
	ErrPermissionRouteDup   = errors.New("相同路由与方法的权限已存在")
)

// PermissionRepository 权限仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/permission_repository.go -package=mock . PermissionRepository
type PermissionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	FindByCode(ctx context.Context, permCode string) (*Permission, error)
	Create(ctx context.Context, permission *Permission) error
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params PermissionExistsParams) (bool, error)
	GetPermissions(ctx context.Context, pager *upagination.Pagination, filter *PermissionListFilter, orderBys ...PermissionListOrderBy) ([]*Permission, int, error)
}

// PermissionInteractor 权限用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/permission_interactor.go -package=mock . PermissionInteractor
type PermissionInteractor interface {
	CreatePermission(ctx context.Context, params *CreatePermissionParams) error
	UpdatePermission(ctx context.Context, params *UpdatePermissionParams) error
	DeletePermission(ctx context.Context, id uuid.UUID) error
	GetPermission(ctx context.Context, id uuid.UUID) (*Permission, error)
	GetPermissions(ctx context.Context, pager *upagination.Pagination, filter *PermissionListFilter, orderBys ...PermissionListOrderBy) ([]*Permission, int, error)
}

type PermissionListOrderByType int

const (
	_ PermissionListOrderByType = iota
	PermissionListOrderByID
	PermissionListOrderByCreatedAt
)

type PermissionListOrderBy struct {
	OrderBy PermissionListOrderByType
	Desc    bool
}

func NewPermissionListOrderByID(desc bool) PermissionListOrderBy {
	return PermissionListOrderBy{OrderBy: PermissionListOrderByID, Desc: desc}
}

func NewPermissionListOrderByCreatedAt(desc bool) PermissionListOrderBy {
	return PermissionListOrderBy{OrderBy: PermissionListOrderByCreatedAt, Desc: desc}
}

type Permission struct {
	ID        uuid.UUID `json:"id"`
	MenuID    uuid.UUID `json:"menu_id"`
	PermCode  string    `json:"perm_code"`
	Name      string    `json:"name"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PermissionListFilter struct {
	MenuID   uuid.UUID `json:"menu_id"`
	PermCode string    `json:"perm_code"`
	Name     string    `json:"name"`
	Method   string    `json:"method"`
	Path     string    `json:"path"`
	Enabled  *bool     `json:"enabled"`
}

type CreatePermissionParams struct {
	MenuID   uuid.UUID `json:"menu_id"`
	PermCode string    `json:"perm_code"`
	Name     string    `json:"name"`
	Method   string    `json:"method"`
	Path     string    `json:"path"`
	Enabled  bool      `json:"enabled"`
}

type UpdatePermissionParams struct {
	ID       uuid.UUID `json:"id"`
	MenuID   uuid.UUID `json:"menu_id"`
	PermCode string    `json:"perm_code"`
	Name     string    `json:"name"`
	Method   string    `json:"method"`
	Path     string    `json:"path"`
	Enabled  bool      `json:"enabled"`
}

type PermissionExistsParams struct {
	PermCode  string    `json:"perm_code"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	ExcludeID uuid.UUID `json:"exclude_id"`
}
