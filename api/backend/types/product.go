package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// ProductCreateReq 创建商品请求
type ProductCreateReq struct {
	// 基础信息
	Name         string                      `json:"name" binding:"required,max=255"` // 商品名称（必选）
	CategoryID   uuid.UUID                   `json:"category_id" binding:"required"`  // 分类ID（必选）
	MenuID       *uuid.UUID                  `json:"menu_id,omitempty"`               // 菜单ID（可选）
	Mnemonic     string                      `json:"mnemonic,omitempty"`              // 助记词（可选）
	ShelfLife    int                         `json:"shelf_life,omitempty"`            // 保质期（可选，单位：天）
	SupportTypes []domain.ProductSupportType `json:"support_types,omitempty"`         // 支持类型（可选，堂食、外带）

	// 属性关联
	UnitID uuid.UUID `json:"unit_id" binding:"required"` // 单位ID（必选）

	// 售卖信息
	SaleStatus         domain.ProductSaleStatus `json:"sale_status" binding:"required,oneof=on_sale off_sale"` // 售卖状态（必选）
	SaleChannels       []domain.SaleChannel     `json:"sale_channels,omitempty"`                               // 售卖渠道（可选，可多选）
	EffectiveDateType  domain.EffectiveDateType `json:"effective_date_type,omitempty"`                         // 生效日期类型（可选）
	EffectiveStartTime *time.Time               `json:"effective_start_time,omitempty"`                        // 生效开始时间（可选）
	EffectiveEndTime   *time.Time               `json:"effective_end_time,omitempty"`                          // 生效结束时间（可选）
	MinSaleQuantity    int                      `json:"min_sale_quantity" binding:"required,min=1"`            // 起售份数（必选）
	AddSaleQuantity    int                      `json:"add_sale_quantity" binding:"required,min=1"`            // 加售份数（必选）

	// 其他信息
	InheritTaxRate bool       `json:"inherit_tax_rate"`      // 是否继承原分类税率（必选）
	TaxRateID      *uuid.UUID `json:"tax_rate_id,omitempty"` // 指定税率ID（可选）
	InheritStall   bool       `json:"inherit_stall"`         // 是否继承原出品部门（必选）
	StallID        *uuid.UUID `json:"stall_id,omitempty"`    // 指定出品部门ID（可选）

	// 展示信息
	MainImage    string   `json:"main_image,omitempty"`    // 主图（可选）
	DetailImages []string `json:"detail_images,omitempty"` // 详情图片（可选，多张）
	Description  string   `json:"description,omitempty"`   // 菜品描述（可选）

	// 关联信息
	SpecRelations []ProductSpecRelationReq `json:"spec_relations" binding:"required,min=1,dive"` // 商品规格关联列表（必选，至少一个）
	AttrRelations []ProductAttrRelationReq `json:"attr_relations,omitempty"`                     // 商品口味做法关联列表（可选）
	TagIDs        []uuid.UUID              `json:"tag_ids,omitempty"`                            // 商品标签ID列表（可选）
}

// ProductSpecRelationReq 商品规格关联请求
type ProductSpecRelationReq struct {
	SpecID             uuid.UUID        `json:"spec_id" binding:"required"`     // 规格ID（必选）
	BasePrice          decimal.Decimal  `json:"base_price" binding:"required"`  // 基础价格（必选，单位：分）
	MemberPrice        *decimal.Decimal `json:"member_price,omitempty"`         // 会员价（可选，单位：分）
	PackingFeeID       uuid.UUID        `json:"packing_fee_id"`                 // 打包费ID（引用费用配置）
	EstimatedCostPrice *decimal.Decimal `json:"estimated_cost_price,omitempty"` // 预估成本价（可选，单位：分）
	OtherPrice1        *decimal.Decimal `json:"other_price1,omitempty"`         // 其他价格1（可选，单位：分）
	OtherPrice2        *decimal.Decimal `json:"other_price2,omitempty"`         // 其他价格2（可选，单位：分）
	OtherPrice3        *decimal.Decimal `json:"other_price3,omitempty"`         // 其他价格3（可选，单位：分）
	Barcode            string           `json:"barcode,omitempty"`              // 条形码（可选）
	IsDefault          bool             `json:"is_default"`                     // 是否默认项
}

// ProductAttrRelationReq 商品口味做法关联请求
type ProductAttrRelationReq struct {
	AttrID     uuid.UUID `json:"attr_id" binding:"required"`      // 口味做法ID（必选）
	AttrItemID uuid.UUID `json:"attr_item_id" binding:"required"` // 口味做法项ID（必选）
	IsDefault  bool      `json:"is_default"`                      // 是否默认项
}
