package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_spec_rel_repository.go -package=mock . ProductSpecRelRepository
type ProductSpecRelRepository interface {
	CreateBulk(ctx context.Context, relations ProductSpecRelations) error
}

// ProductSpecRelation 商品-规格关联实体
type ProductSpecRelation struct {
	ID                 uuid.UUID        `json:"id"`                   // 关联ID
	ProductID          uuid.UUID        `json:"product_id"`           // 商品ID
	SpecID             uuid.UUID        `json:"spec_id"`              // 规格ID
	BasePrice          decimal.Decimal  `json:"base_price"`           // 基础价格（单位：分）
	MemberPrice        *decimal.Decimal `json:"member_price"`         // 会员价（单位：分，可选）
	PackingFeeID       uuid.UUID        `json:"packing_fee_id"`       // 打包费ID（引用费用配置）
	EstimatedCostPrice *decimal.Decimal `json:"estimated_cost_price"` // 预估成本价（单位：分，可选）
	OtherPrice1        *decimal.Decimal `json:"other_price1"`         // 其他价格1（单位：分，可选）
	OtherPrice2        *decimal.Decimal `json:"other_price2"`         // 其他价格2（单位：分，可选）
	OtherPrice3        *decimal.Decimal `json:"other_price3"`         // 其他价格3（单位：分，可选）
	Barcode            string           `json:"barcode"`              // 条形码
	IsDefault          bool             `json:"is_default"`           // 是否默认项
	CreatedAt          time.Time        `json:"created_at"`           // 创建时间
	UpdatedAt          time.Time        `json:"updated_at"`           // 更新时间

	// 关联信息
	SpecName   string      `json:"spec_name"`   // 规格名称
	PackingFee *PackingFee `json:"packing_fee"` // 打包费
}

// @TODO 关联信息
type PackingFee struct {
	ID    uuid.UUID       `json:"id"`    // 打包费ID
	Name  string          `json:"name"`  // 打包费名称
	Price decimal.Decimal `json:"price"` // 打包费价格
}

type ProductSpecRelations []*ProductSpecRelation
