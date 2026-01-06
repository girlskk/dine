package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrRoleMenuNotExists = errors.New("角色菜单关系不存在")
	ErrRoleMenuExists    = errors.New("角色菜单关系已存在")
)

// RoleMenuRepository 角色菜单关系仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/role_menu_repository.go -package=mock . RoleMenuRepository
type RoleMenuRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*RoleMenu, error)
	Create(ctx context.Context, roleMenu *RoleMenu) error
	CreateBulk(ctx context.Context, roleMenus []*RoleMenu) error
	CreateBulkByRoleIDMenus(ctx context.Context, role *Role, menuIDs []uuid.UUID) error
	DeletesByRoleID(ctx context.Context, roleID uuid.UUID) error
	Deletes(ctx context.Context, ids []uuid.UUID) error
	GetRoleMenus(ctx context.Context, pager *upagination.Pagination, filter *RoleMenuListFilter, orderBys ...RoleMenuListOrderBy) ([]*RoleMenu, int, error)
}

type RoleMenuListOrderByType int

const (
	_ RoleMenuListOrderByType = iota
	RoleMenuListOrderByID
	RoleMenuListOrderByCreatedAt
)

type RoleMenuListOrderBy struct {
	OrderBy RoleMenuListOrderByType
	Desc    bool
}

func NewRoleMenuListOrderByID(desc bool) RoleMenuListOrderBy {
	return RoleMenuListOrderBy{OrderBy: RoleMenuListOrderByID, Desc: desc}
}

func NewRoleMenuListOrderByCreatedAt(desc bool) RoleMenuListOrderBy {
	return RoleMenuListOrderBy{OrderBy: RoleMenuListOrderByCreatedAt, Desc: desc}
}

type RoleMenu struct {
	ID         uuid.UUID `json:"id"`
	RoleType   RoleType  `json:"role_type"`
	RoleID     uuid.UUID `json:"role_id"`
	MenuID     uuid.UUID `json:"menu_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RoleMenuListFilter struct {
	RoleType   RoleType  `json:"role_type"`
	RoleID     uuid.UUID `json:"role_id"`
	MenuID     uuid.UUID `json:"menu_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
}
