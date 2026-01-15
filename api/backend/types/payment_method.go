package types

import (
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// PaymentMethodCreateReq 创建结算方式请求
type PaymentMethodCreateReq struct {
	Name             string                               `json:"name" binding:"required,max=255"`                                                             // 结算方式名称（必选）
	AccountingRule   domain.PaymentMethodAccountingRule   `json:"accounting_rule" binding:"required,oneof=income discount"`                                    // 计入规则（必选）
	PaymentType      domain.PaymentMethodPayType          `json:"payment_type" binding:"required,oneof=other cash offline_card custom_coupon partner_coupon"`  // // 结算类型（必选）
	FeeRate          *decimal.Decimal                     `json:"fee_rate"`                                                                                    // 手续费率,百分比
	InvoiceRule      domain.PaymentMethodInvoiceRule      `json:"invoice_rule" binding:"required,oneof=no_invoice actual_amount"`                              // 实收部分开票规则（必选）
	CashDrawerStatus bool                                 `json:"cash_drawer_status"`                                                                          // 开钱箱状态
	Status           bool                                 `json:"status"`                                                                                      // 启用/停用状态
	DisplayChannels  []domain.PaymentMethodDisplayChannel `json:"display_channels" binding:"required,min=1,dive,oneof=POS Mobile Scan SelfService ThirdParty"` // 收银终端显示渠道（可选，可多选）
	Source           domain.PaymentMethodSource           `json:"source" binding:"required,oneof=brand store system"`                                          // 来源:brand-品牌,store-门店,system-系统
}

// PaymentMethodUpdateReq 更新结算方式请求
type PaymentMethodUpdateReq struct {
	Name             string                               `json:"name" binding:"required,max=255"`                                                             // 结算方式名称（必选）
	AccountingRule   domain.PaymentMethodAccountingRule   `json:"accounting_rule" binding:"required,oneof=income discount"`                                    // 计入规则（必选）
	PaymentType      domain.PaymentMethodPayType          `json:"payment_type" binding:"required,oneof=other cash offline_card custom_coupon partner_coupon"`  // // 结算类型（必选）
	FeeRate          *decimal.Decimal                     `json:"fee_rate"`                                                                                    // 手续费率,百分比
	InvoiceRule      domain.PaymentMethodInvoiceRule      `json:"invoice_rule" binding:"required,oneof=no_invoice actual_amount"`                              // 实收部分开票规则（必选）
	CashDrawerStatus bool                                 `json:"cash_drawer_status" binding:"oneof=true false"`                                               // 开钱箱状态
	Status           bool                                 `json:"status" binding:"oneof=true false"`                                                           // 启用/停用状态
	DisplayChannels  []domain.PaymentMethodDisplayChannel `json:"display_channels" binding:"required,min=1,dive,oneof=POS Mobile Scan SelfService ThirdParty"` // 收银终端显示渠道（可选，可多选）
}

// PaymentMethodListReq 结算方式列表请求
type PaymentMethodListReq struct {
	upagination.RequestPagination
	Name   string                     `json:"name" form:"name"`     // 结算方式名称（模糊匹配）
	Source domain.PaymentMethodSource `json:"source" form:"source"` // 来源:brand-品牌,store-门店,system-系统
}

// PaymentMethodStatReq 统计各个结算分类对应的结算方式数量
type PaymentMethodStatReq struct {
	Name   string                     `json:"name" form:"name"`     // 结算方式名称（模糊匹配）
	Source domain.PaymentMethodSource `json:"source" form:"source"` // 来源:brand-品牌,store-门店,system-系统
}
