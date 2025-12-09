package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrOrderCartNotFound = errors.New("购物车商品不存在")
)

// OrderCartRepository 购物车仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/order_cart_repository.go -package=mock . OrderCartRepository
type OrderCartRepository interface {
	// 获取指定桌台的购物车商品
	ListByTable(ctx context.Context, tableID int, withExtra bool) (OrderCarts, error)
	// 根据唯一键查找购物车商品
	FindByUniqueKey(ctx context.Context, key OrderCartItemUniqueKey) (item *OrderCart, err error)
	// 根据ID查找购物车商品
	FindByID(ctx context.Context, id int) (item *OrderCart, err error)
	// 创建购物车商品
	Create(ctx context.Context, item *OrderCart) (err error)
	// 原子减少数量，如果数量为1，则删除记录
	DecrementQuantity(ctx context.Context, id int) (err error)
	// 原子增加数量
	IncrementQuantity(ctx context.Context, id int) (err error)
	// 清空购物车
	ClearByTable(ctx context.Context, tableID int) (err error)
}

// OrderCartInteractor 购物车用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/order_cart_interactor.go -package=mock . OrderCartInteractor
type OrderCartInteractor interface {
	// 添加商品到购物车
	AddItem(ctx context.Context, params OrderCartAddParams) (res OrderCarts, err error)
	// 从购物车移除商品
	RemoveItem(ctx context.Context, id int, tableID int) (res OrderCarts, err error)
	// 获取指定桌台的购物车商品
	ListByTable(ctx context.Context, tableID int) (res OrderCarts, err error)
}

// OrderCartItemUniqueKey 购物车项唯一键
type OrderCartItemUniqueKey struct {
	TableID       int // 台桌ID
	ProductID     int // 商品ID
	ProductSpecID int // 商品规格ID（可选）
	AttrID        int // 属性ID（可选）
	RecipeID      int // 做法ID（可选）
}

type OrderCartAddParams OrderCartItemUniqueKey

type (
	OrderCart struct {
		ID            int             `json:"id"`                        // 购物车ID
		TableID       int             `json:"table_id"`                  // 台桌ID
		ProductID     int             `json:"product_id"`                // 商品ID
		ProductSpecID int             `json:"product_spec_id,omitempty"` // 商品规格ID
		AttrID        int             `json:"attr_id,omitempty"`         // 属性ID
		RecipeID      int             `json:"recipe_id,omitempty"`       // 做法ID
		Quantity      decimal.Decimal `json:"quantity"`                  // 商品数量
		CreatedAt     time.Time       `json:"created_at"`                // 创建时间
		UpdatedAt     time.Time       `json:"updated_at"`                // 更新时间
		// 关联信息
		Name           string          `json:"name"`                       // 商品名称
		Price          decimal.Decimal `json:"price"`                      // 商品价格
		Images         []string        `json:"images,omitempty"`           // 商品图片
		SpecName       string          `json:"spec_name,omitempty"`        // 商品规格
		AttrName       string          `json:"attr_name,omitempty"`        // 商品属性
		RecipeName     string          `json:"recipe_name,omitempty"`      // 商品做法
		CategoryID     int             `json:"category_id,omitempty"`      // 商品分类ID
		CategoryName   string          `json:"category_name,omitempty"`    // 商品分类名称
		SetMealDetails SetMealDetails  `json:"set_meal_details,omitempty"` // 套餐详情
	}
)

// OrderCarts 购物车项目集合
type OrderCarts []*OrderCart
