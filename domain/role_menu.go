package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// RoleMenuRepository 角色菜单关系仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/role_menu_repository.go -package=mock . RoleMenuRepository
type RoleMenuRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*RoleMenu, error)
	Create(ctx context.Context, roleMenu *RoleMenu) error
	CreateBulk(ctx context.Context, roleMenus []*RoleMenu) error
	CreateBulkByRoleIDPaths(ctx context.Context, role *Role, paths []string) error
	DeletesByRoleID(ctx context.Context, roleID uuid.UUID) error
	Deletes(ctx context.Context, ids []uuid.UUID) error
	GetRoleMenus(ctx context.Context, pager *upagination.Pagination, filter *RoleMenuListFilter, orderBys ...RoleMenuListOrderBy) ([]*RoleMenu, int, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]*RoleMenu, error)
}

type RoleMenuInteractor interface {
	SetRoleMenu(ctx context.Context, roleID uuid.UUID, paths []string, user User) error
	RoleMenuList(ctx context.Context, roleID uuid.UUID) (paths []string, err error)
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
	Path       string    `json:"path"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RoleMenuListFilter struct {
	RoleType   RoleType  `json:"role_type"`
	RoleID     uuid.UUID `json:"role_id"`
	Path       string    `json:"path"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
}
