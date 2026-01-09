package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrRolePermissionNotExists = errors.New("角色权限关系不存在")
	ErrRolePermissionExists    = errors.New("角色权限关系已存在")
)

// RolePermissionRepository 角色权限关系仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/role_permission_repository.go -package=mock . RolePermissionRepository
type RolePermissionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*RolePermission, error)
	Create(ctx context.Context, rolePermission *RolePermission) error
	CreateBulk(ctx context.Context, rolePermissions []*RolePermission) error
	CreateBulkByRoleIDPermissions(ctx context.Context, role *Role, permissionIDs []uuid.UUID) error
	DeletesByRoleID(ctx context.Context, roleID uuid.UUID) error
	Deletes(ctx context.Context, ids []uuid.UUID) error
	GetRolePermissions(ctx context.Context, pager *upagination.Pagination, filter *RolePermissionListFilter, orderBys ...RolePermissionListOrderBy) ([]*RolePermission, int, error)
}

type RolePermissionListOrderByType int

const (
	_ RolePermissionListOrderByType = iota
	RolePermissionListOrderByID
	RolePermissionListOrderByCreatedAt
)

type RolePermissionListOrderBy struct {
	OrderBy RolePermissionListOrderByType
	Desc    bool
}

func NewRolePermissionListOrderByID(desc bool) RolePermissionListOrderBy {
	return RolePermissionListOrderBy{OrderBy: RolePermissionListOrderByID, Desc: desc}
}

func NewRolePermissionListOrderByCreatedAt(desc bool) RolePermissionListOrderBy {
	return RolePermissionListOrderBy{OrderBy: RolePermissionListOrderByCreatedAt, Desc: desc}
}

type RolePermission struct {
	ID           uuid.UUID `json:"id"`
	RoleType     RoleType  `json:"role_type"`
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	MerchantID   uuid.UUID `json:"merchant_id"`
	StoreID      uuid.UUID `json:"store_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RolePermissionListFilter struct {
	RoleType     RoleType  `json:"role_type"`
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	MerchantID   uuid.UUID `json:"merchant_id"`
	StoreID      uuid.UUID `json:"store_id"`
}
