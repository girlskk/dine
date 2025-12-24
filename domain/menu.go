package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ------------------------------------------------------------
// 错误定义
// ------------------------------------------------------------

var (
	ErrMenuNotExists          = errors.New("菜单不存在")
	ErrMenuNameExists         = errors.New("菜单名称已存在")
	ErrMenuStoreBound         = errors.New("门店已绑定其他菜单")
	ErrMenuHasStores          = errors.New("菜单下有关联门店，不能删除")
	ErrMenuItemProductInvalid = errors.New("菜品无效，必须属于当前品牌商")
)

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// MenuDistributionRule 下发规则
type MenuDistributionRule string

const (
	MenuDistributionRuleOverride MenuDistributionRule = "override" // 新增并覆盖同名菜品
	MenuDistributionRuleKeep     MenuDistributionRule = "keep"     // 对同名菜品不做修改
)

func (MenuDistributionRule) Values() []string {
	return []string{
		string(MenuDistributionRuleOverride),
		string(MenuDistributionRuleKeep),
	}
}

// MenuItemSaleRule 下发售卖规则
type MenuItemSaleRule string

const (
	MenuItemSaleRuleKeepBrandStatus MenuItemSaleRule = "keep_brand_status" // 保留品牌状态
	MenuItemSaleRuleKeepStoreStatus MenuItemSaleRule = "keep_store_status" // 保留门店状态
)

func (MenuItemSaleRule) Values() []string {
	return []string{
		string(MenuItemSaleRuleKeepBrandStatus),
		string(MenuItemSaleRuleKeepStoreStatus),
	}
}

// ------------------------------------------------------------
// 实体定义
// ------------------------------------------------------------

// Menu 菜单实体
type Menu struct {
	ID               uuid.UUID            `json:"id"`                // 菜单ID
	MerchantID       uuid.UUID            `json:"merchant_id"`       // 品牌商ID
	StoreID          uuid.UUID            `json:"store_id"`          // 门店ID
	Name             string               `json:"name"`              // 菜单名称
	DistributionRule MenuDistributionRule `json:"distribution_rule"` // 下发规则
	StoreCount       int                  `json:"store_count"`       // 适用门店数量
	ItemCount        int                  `json:"item_count"`        // 菜单项数量
	CreatedAt        time.Time            `json:"created_at"`        // 创建时间
	UpdatedAt        time.Time            `json:"updated_at"`        // 更新时间

	// 关联信息
	Stores []*StoreSimple `json:"stores,omitempty"` // 关联门店列表
	Items  MenuItems      `json:"items,omitempty"`  // 菜单项列表
}

// Menus 菜单集合
type Menus []*Menu

// MenuItem 菜单项实体
type MenuItem struct {
	ID          uuid.UUID        `json:"id"`           // 菜单项ID
	MenuID      uuid.UUID        `json:"menu_id"`      // 菜单ID
	ProductID   uuid.UUID        `json:"product_id"`   // 菜品ID
	SaleRule    MenuItemSaleRule `json:"sale_rule"`    // 下发售卖规则
	BasePrice   *decimal.Decimal `json:"base_price"`   // 基础价（可选，单位：分）
	MemberPrice *decimal.Decimal `json:"member_price"` // 会员价（可选，单位：分）
	CreatedAt   time.Time        `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time        `json:"updated_at"`   // 更新时间

	// 关联信息
	Product *Product `json:"product,omitempty"` // 关联商品
}

// MenuItems 菜单项集合
type MenuItems []*MenuItem

// ------------------------------------------------------------
// 仓储和用例接口
// ------------------------------------------------------------

// MenuRepository 菜单仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/menu_repository.go -package=mock . MenuRepository
type MenuRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Menu, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*Menu, error)
	Create(ctx context.Context, menu *Menu) error
	Update(ctx context.Context, menu *Menu) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params MenuExistsParams) (bool, error)
	CheckStoreBound(ctx context.Context, storeIDs []uuid.UUID, excludeMenuID uuid.UUID) (bool, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params MenuSearchParams) (*MenuSearchRes, error)
}

// MenuInteractor 菜单用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/menu_interactor.go -package=mock . MenuInteractor
type MenuInteractor interface {
	Create(ctx context.Context, menu *Menu) error
	Update(ctx context.Context, menu *Menu) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDetail(ctx context.Context, id uuid.UUID) (*Menu, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params MenuSearchParams) (*MenuSearchRes, error)
}

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// MenuExistsParams 存在性检查参数
type MenuExistsParams struct {
	MerchantID uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// MenuSearchParams 查询参数
type MenuSearchParams struct {
	MerchantID uuid.UUID
	Name       string // 菜单名称（模糊匹配）
}

// MenuSearchRes 查询结果
type MenuSearchRes struct {
	*upagination.Pagination
	Items Menus `json:"items"`
}
