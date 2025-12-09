package domain

import (
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type (
	// 订单商品信息快照
	OrderProductInfoSnapshot struct {
		Price      decimal.Decimal `json:"price"`                 // 商品价格
		Images     []string        `json:"images"`                // 商品图片
		UnitID     int             `json:"unit_id"`               // 单位ID
		UnitName   string          `json:"unit_name"`             // 单位名称
		SpecID     int             `json:"spec_id,omitempty"`     // 规格ID
		SpecName   string          `json:"spec_name,omitempty"`   // 规格名称
		SpecPrice  decimal.Decimal `json:"spec_price,omitempty"`  // 规格价格
		AttrID     int             `json:"attr_id,omitempty"`     // 属性ID
		AttrName   string          `json:"attr_name,omitempty"`   // 属性名称
		RecipeID   int             `json:"recipe_id,omitempty"`   // 做法ID
		RecipeName string          `json:"recipe_name,omitempty"` // 做法名称
	}

	// 订单商品
	OrderItem struct {
		ID              int                      `json:"id"`               // 订单商品项ID
		OrderID         int                      `json:"order_id"`         // 订单ID
		ProductID       int                      `json:"product_id"`       // 商品ID
		Name            string                   `json:"name"`             // 商品名称
		Type            ProductType              `json:"type"`             // 商品类型：1-单规格 2-多规格 3-套餐商品
		AllowPointPay   bool                     `json:"allow_point_pay"`  // 是否支持积分支付
		Quantity        decimal.Decimal          `json:"quantity"`         // 商品数量
		Price           decimal.Decimal          `json:"price"`            // 商品实际单价
		Amount          decimal.Decimal          `json:"amount"`           // 商品总金额（商品数量 * 商品实际单价）
		CreatedAt       time.Time                `json:"created_at"`       // 创建时间
		UpdatedAt       time.Time                `json:"updated_at"`       // 更新时间
		ProductSnapshot OrderProductInfoSnapshot `json:"product_snapshot"` // 商品信息快照
		Remark          string                   `json:"remark"`           // 备注

		SetMealDetails []*OrderItemSetMealDetail `json:"set_meal_details,omitempty"` // 套餐商品详情（类型是套餐商品才有）
	}

	OrderItems []*OrderItem

	// 订单套餐商品详情
	OrderItemSetMealDetail struct {
		ID              int                      `json:"id"`
		OrderItemID     int                      `json:"order_item_id"`    // 订单商品项ID
		Name            string                   `json:"name"`             // 商品名称
		Type            ProductType              `json:"type"`             // 商品类型：1-单规格 2-多规格
		SetMealPrice    decimal.Decimal          `json:"set_meal_price"`   // 套餐内价格
		SetMealID       int                      `json:"set_meal_id"`      // 套餐ID
		ProductID       int                      `json:"product_id"`       // 商品ID
		Quantity        decimal.Decimal          `json:"quantity"`         // 数量（支持3位小数）
		ProductSnapshot OrderProductInfoSnapshot `json:"product_snapshot"` // 商品信息快照
		CreatedAt       time.Time                `json:"created_at"`       // 创建时间
		UpdatedAt       time.Time                `json:"updated_at"`       // 更新时间
	}
)

func (items OrderItems) TotalAmount() decimal.Decimal {
	return lo.Reduce(items, func(agg decimal.Decimal, item *OrderItem, _ int) decimal.Decimal {
		return agg.Add(item.Amount)
	}, decimal.Zero)
}

func (items OrderItems) PointsAvailable(storeID int) decimal.Decimal {
	amt := lo.Reduce(lo.Filter(items, func(item *OrderItem, _ int) bool {
		return item.AllowPointPay
	}), func(agg decimal.Decimal, item *OrderItem, _ int) decimal.Decimal {
		return agg.Add(item.Amount)
	}, decimal.Zero)

	// 边疆故人和极品牛庄改为0.5
	if storeID <= 2 {
		return amt.Mul(decimal.NewFromFloat(0.5))
	}

	return amt.Mul(decimal.NewFromFloat(0.3))
}
