package types

import (
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// AdditionalFeeCreateReq 创建附加费
type AdditionalFeeCreateReq struct {
	Name                string                            `json:"name" binding:"required,max=50"`                                               // 名称
	FeeCategory         domain.AdditionalCategory         `json:"fee_category" binding:"required,oneof=service_fee additional_fee packing_fee"` // 附加费分类 (service_fee/additional_fee/packing_fee)
	ChargeMode          domain.AdditionalFeeChargeMode    `json:"charge_mode" binding:"required,oneof=percent fixed"`                           // 收费模式 (percent/fixed)
	FeeValue            decimal.Decimal                   `json:"fee_value" binding:"required"`                                                 // 费用值，percent 模式示例 0.06 表示 6%
	IncludeInReceivable bool                              `json:"include_in_receivable"`                                                        // 是否计入实收
	Taxable             bool                              `json:"taxable"`                                                                      // 附加费是否收税
	DiscountScope       domain.AdditionalFeeDiscountScope `json:"discount_scope" binding:"required,oneof=before_discount after_discount"`       // 折扣作用范围 (before_discount/after_discount)
	OrderChannels       []domain.OrderChannel             `json:"order_channels" binding:"required"`                                            // 允许的下单渠道
	DiningWays          []domain.DiningWay                `json:"dining_ways" binding:"required"`                                               // 适用用餐方式
	Enabled             bool                              `json:"enabled"`                                                                      // 是否启用
	SortOrder           int                               `json:"sort_order" binding:"omitempty,gte=0"`                                         // 排序值
}

// AdditionalFeeUpdateReq 更新附加费
type AdditionalFeeUpdateReq struct {
	Name                string                            `json:"name" binding:"required,max=50"`                                               // 名称
	FeeCategory         domain.AdditionalCategory         `json:"fee_category" binding:"required,oneof=service_fee additional_fee packing_fee"` // 附加费分类 (service_fee/additional_fee/packing_fee)
	ChargeMode          domain.AdditionalFeeChargeMode    `json:"charge_mode" binding:"required,oneof=percent fixed"`                           // 收费模式 (percent/fixed)
	FeeValue            decimal.Decimal                   `json:"fee_value" binding:"required"`                                                 // 费用值，percent 模式示例 0.06 表示 6%
	IncludeInReceivable bool                              `json:"include_in_receivable"`                                                        // 是否计入实收
	Taxable             bool                              `json:"taxable"`                                                                      // 附加费是否收税
	DiscountScope       domain.AdditionalFeeDiscountScope `json:"discount_scope" binding:"required,oneof=before_discount after_discount"`       // 折扣作用范围 (before_discount/after_discount)
	OrderChannels       []domain.OrderChannel             `json:"order_channels" binding:"required"`                                            // 允许的下单渠道
	DiningWays          []domain.DiningWay                `json:"dining_ways" binding:"required"`                                               // 适用用餐方式
	Enabled             bool                              `json:"enabled"`                                                                      // 是否启用
	SortOrder           int                               `json:"sort_order" binding:"omitempty,gte=0"`                                         // 排序值
}

// AdditionalFeeListReq 列表查询
type AdditionalFeeListReq struct {
	upagination.RequestPagination
	Name    string `form:"name"`    // 名称
	Enabled *bool  `form:"enabled"` // 是否启用
}

// AdditionalFeeListResp 附加费列表响应
type AdditionalFeeListResp struct {
	AdditionalFees []*domain.AdditionalFee `json:"additional_fees"` // 附加费列表
	Total          int                     `json:"total"`           // 总数
}
