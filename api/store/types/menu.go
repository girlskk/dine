package types

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// MenuCreateReq 创建菜单请求
type MenuCreateReq struct {
	Name  string        `json:"name" binding:"required,max=255"`     // 菜单名称（必选）
	Items []MenuItemReq `json:"items" binding:"required,min=1,dive"` // 菜品列表（必选，至少一个）
}

// MenuItemReq 菜单项请求
type MenuItemReq struct {
	ProductID   uuid.UUID        `json:"product_id" binding:"required"` // 菜品ID（必选）
	BasePrice   *decimal.Decimal `json:"base_price,omitempty"`          // 基础价（可选，单位：分）
	MemberPrice *decimal.Decimal `json:"member_price,omitempty"`        // 会员价（可选，单位：分）
}

// MenuUpdateReq 更新菜单请求
type MenuUpdateReq struct {
	Name  string        `json:"name" binding:"required,max=255"`     // 菜单名称（必选）
	Items []MenuItemReq `json:"items" binding:"required,min=1,dive"` // 菜品列表（必选，至少一个）
}

// MenuListReq 菜单列表请求
type MenuListReq struct {
	upagination.RequestPagination
	Name string `json:"name" form:"name"` // 菜单名称（模糊匹配）
}
