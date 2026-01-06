package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrRouterMenuNotExists        = errors.New("菜单不存在")
	ErrRouterMenuNameExists       = errors.New("同级菜单名称已存在")
	ErrRouterMenuForbidenAddChild = errors.New("禁止添加子菜单")
)

// RouterMenuRepository 菜单仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/router_menu_repository.go -package=mock . RouterMenuRepository
type RouterMenuRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*RouterMenu, error)
	Create(ctx context.Context, menu *RouterMenu) error
	Update(ctx context.Context, menu *RouterMenu) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params RouterMenuExistsParams) (bool, error)
	GetRouterMenus(ctx context.Context, pager *upagination.Pagination, filter *RouterMenuListFilter, orderBys ...RouterMenuListOrderBy) ([]*RouterMenu, int, error)
}

// RouterMenuInteractor 菜单用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/router_menu_interactor.go -package=mock . RouterMenuInteractor
type RouterMenuInteractor interface {
	CreateRouterMenu(ctx context.Context, params *CreateRouterMenuParams) error
	UpdateRouterMenu(ctx context.Context, params *UpdateRouterMenuParams) error
	DeleteRouterMenu(ctx context.Context, id uuid.UUID) error
	GetRouterMenu(ctx context.Context, id uuid.UUID) (*RouterMenu, error)
	GetRouterMenus(ctx context.Context, pager *upagination.Pagination, filter *RouterMenuListFilter, orderBys ...RouterMenuListOrderBy) ([]*RouterMenu, int, error)
}

type RouterMenuListOrderByType int

const (
	_ RouterMenuListOrderByType = iota
	RouterMenuListOrderByID
	RouterMenuListOrderByCreatedAt
	RouterMenuListOrderBySort
)

type RouterMenuListOrderBy struct {
	OrderBy RouterMenuListOrderByType
	Desc    bool
}

func NewRouterMenuListOrderByID(desc bool) RouterMenuListOrderBy {
	return RouterMenuListOrderBy{OrderBy: RouterMenuListOrderByID, Desc: desc}
}

func NewRouterMenuListOrderByCreatedAt(desc bool) RouterMenuListOrderBy {
	return RouterMenuListOrderBy{OrderBy: RouterMenuListOrderByCreatedAt, Desc: desc}
}

func NewRouterMenuListOrderBySort(desc bool) RouterMenuListOrderBy {
	return RouterMenuListOrderBy{OrderBy: RouterMenuListOrderBySort, Desc: desc}
}

type RouterMenu struct {
	ID        uuid.UUID `json:"id"`
	UserType  UserType  `json:"user_type"`
	ParentID  uuid.UUID `json:"parent_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Icon      string    `json:"icon"`
	Sort      int       `json:"sort"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RouterMenuListFilter struct {
	UserType UserType  `json:"user_type"`
	ParentID uuid.UUID `json:"parent_id"`
	Name     string    `json:"name"`
	Enabled  *bool     `json:"enabled"`
}

type CreateRouterMenuParams struct {
	UserType  UserType  `json:"user_type"`
	ParentID  uuid.UUID `json:"parent_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Icon      string    `json:"icon"`
	Sort      int       `json:"sort"`
	Enabled   bool      `json:"enabled"`
}

type UpdateRouterMenuParams struct {
	ID        uuid.UUID `json:"id"`
	ParentID  uuid.UUID `json:"parent_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Icon      string    `json:"icon"`
	Sort      int       `json:"sort"`
	Enabled   bool      `json:"enabled"`
}

type RouterMenuExistsParams struct {
	UserType  UserType  `json:"user_type"`
	ParentID  uuid.UUID `json:"parent_id"`
	Name      string    `json:"name"`
	ExcludeID uuid.UUID `json:"exclude_id"`
}
