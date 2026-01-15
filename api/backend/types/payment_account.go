package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// PaymentAccountCreateReq 创建收款账户请求
type PaymentAccountCreateReq struct {
	Channel        domain.PaymentChannel `json:"channel" binding:"required,oneof=rm"`        // 支付渠道（必选）
	MerchantNumber string                `json:"merchant_number" binding:"required,max=255"` // 支付商户号（必选）
	MerchantName   string                `json:"merchant_name" binding:"required,max=255"`   // 支付商户名称（必选）
}

// PaymentAccountUpdateReq 更新收款账户请求
type PaymentAccountUpdateReq struct {
	Channel        domain.PaymentChannel `json:"channel" binding:"required,oneof=rm"`        // 支付渠道（必选）
	MerchantNumber string                `json:"merchant_number" binding:"required,max=255"` // 支付商户号（必选）
	MerchantName   string                `json:"merchant_name" binding:"required,max=255"`   // 支付商户名称（必选）
}

// PaymentAccountListReq 收款账户列表请求
type PaymentAccountListReq struct {
	Page           int                   `form:"page" binding:"omitempty,min=1"`       // 页码
	Size           int                   `form:"size" binding:"omitempty,min=1"`       // 每页数量
	Channel        domain.PaymentChannel `form:"channel" binding:"omitempty,oneof=rm"` // 支付渠道（可选）
	MerchantName   string                `form:"merchant_name" binding:"omitempty"`    // 支付商户名称（可选，模糊匹配）
	CreatedAtStart string                `form:"created_at_start" binding:"omitempty"` // 创建时间开始（可选）
	CreatedAtEnd   string                `form:"created_at_end" binding:"omitempty"`   // 创建时间结束（可选）
}
