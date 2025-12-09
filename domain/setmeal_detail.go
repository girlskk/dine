package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// SetMealDetailRepository 仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/setmeal_detail_repository.go -package=mock . SetMealDetailRepository
type SetMealDetailRepository interface {
	ListBySetMealID(ctx context.Context, id int) (SetMealDetails, error)
	BatchCreate(ctx context.Context, details SetMealDetails) error
	DeleteBySetMealID(ctx context.Context, setMealID int) error
}

// 套餐商品详情
type SetMealDetail struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`           // 商品名称
	Price        decimal.Decimal `json:"price"`          // 商品价格
	SetMealPrice decimal.Decimal `json:"set_meal_price"` // 套餐内价格
	SetMealID    int             `json:"set_meal_id"`    // 套餐商品ID
	ProductID    int             `json:"product_id"`     // 商品详情ID
	Quantity     decimal.Decimal `json:"quantity"`       // 数量（支持3位小数）
	CreatedAt    time.Time       `json:"created_at"`     // 创建时间
	UpdatedAt    time.Time       `json:"updated_at"`     // 更新时间
	UnitID       int             `json:"unit_id"`        // 单位ID
	CategoryID   int             `json:"category_id"`    // 分类ID
	ProductType  ProductType     `json:"product_type"`   // 商品类型
	Images       []string        `json:"images"`         // 商品图片

	Unit  *ProductUnit    `json:"unit"`  // 商品单位
	Spec  *ProductSpecRel `json:"spec"`  // 套餐商品具体规格
	Specs ProductSpecRels `json:"specs"` // 商品规格
}

type SetMealDetails []*SetMealDetail
