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
	ListByIDs(ctx context.Context, ids ...uuid.UUID) ([]*Role, error)
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
	SimpleUpdate(ctx context.Context, updateField RoleSimpleUpdateField, params RoleSimpleUpdateParams) error
}

type RoleSimpleUpdateField string

const (
	RoleSimpleUpdateFieldEnable RoleSimpleUpdateField = "enable"
)

type RoleDataScopeType string

const (
	RoleDataScopeAll        RoleDataScopeType = "all"        // 全部数据权限
	RoleDataScopeMerchant   RoleDataScopeType = "merchant"   // 品牌商数据权限
	RoleDataScopeStore      RoleDataScopeType = "store"      // 门店数据权限
	RoleDataScopeDepartment RoleDataScopeType = "department" // 部门数据权限
	RoleDataScopeSelf       RoleDataScopeType = "self"       // 仅本人数据权限
	RoleDataScopeCustom     RoleDataScopeType = "custom"     // 自定义数据权限
)

func (RoleDataScopeType) Values() []string {
	return []string{
		string(RoleDataScopeAll),
		string(RoleDataScopeMerchant),
		string(RoleDataScopeStore),
		string(RoleDataScopeDepartment),
		string(RoleDataScopeSelf),
		string(RoleDataScopeCustom),
	}
}

type LoginChannel string

const (
	LoginChannelPos    LoginChannel = "pos"    // pos
	LoginChannelMobile LoginChannel = "mobile" // 移动点餐
	LoginChannelStore  LoginChannel = "store"  // 门店管理后台
)

// RoleType 角色类型
type RoleType UserType

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
	ID            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`           // 角色名称
	Code          string            `json:"code"`           // 角色编码
	RoleType      RoleType          `json:"role_type"`      // 角色类型
	DataScope     RoleDataScopeType `json:"data_scope"`     // 数据权限范围
	Enable        bool              `json:"enable"`         // 是否启用
	MerchantID    uuid.UUID         `json:"merchant_id"`    // 所属商户 ID
	StoreID       uuid.UUID         `json:"store_id"`       // 所属门店 ID
	LoginChannels []LoginChannel    `json:"login_channels"` // 允许登录渠道，取自 login_channel，多选
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type RoleListFilter struct {
	Name       string    `json:"name"`
	RoleType   RoleType  `json:"role_type"`
	Enable     *bool     `json:"enable"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
}

type CreateRoleParams struct {
	Name          string            `json:"name"`
	Code          string            `json:"code"`
	RoleType      RoleType          `json:"role_type"`
	DataScope     RoleDataScopeType `json:"data_scope"`
	Enable        bool              `json:"enable"`
	LoginChannels []LoginChannel    `json:"login_channels"` // 允许登录渠道，取自 login_channel，多选
	MerchantID    uuid.UUID         `json:"merchant_id"`
	StoreID       uuid.UUID         `json:"store_id"`
}

type UpdateRoleParams struct {
	ID            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`
	RoleType      RoleType          `json:"role_type"`
	DataScope     RoleDataScopeType `json:"data_scope"`
	Enable        bool              `json:"enable"`
	LoginChannels []LoginChannel    `json:"login_channels"` // 允许登录渠道，取自 login_channel，多选
	MerchantID    uuid.UUID         `json:"merchant_id"`
	StoreID       uuid.UUID         `json:"store_id"`
}

type RoleExistsParams struct {
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	ExcludeID  uuid.UUID `json:"exclude_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
}

type RoleSimpleUpdateParams struct {
	ID     uuid.UUID `json:"id"`
	Enable bool      `json:"enable"`
}
