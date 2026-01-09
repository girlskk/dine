package types

import (
	"github.com/google/uuid"
)

// StorePaymentAccountCreateReq 创建门店收款账户请求
type StorePaymentAccountCreateReq struct {
	StoreID          uuid.UUID `json:"store_id" binding:"required,uuid"`           // 门店ID（必选）
	PaymentAccountID uuid.UUID `json:"payment_account_id" binding:"required,uuid"` // 品牌商收款账户ID（必选）
	MerchantNumber   string    `json:"merchant_number" binding:"required,max=255"` // 支付商户号（必选）
}

// StorePaymentAccountUpdateReq 更新门店收款账户请求
type StorePaymentAccountUpdateReq struct {
	MerchantNumber string `json:"merchant_number" binding:"required,max=255"` // 支付商户号（必选）
}

// StorePaymentAccountListReq 门店收款账户列表请求
type StorePaymentAccountListReq struct {
	Page           int      `form:"page" binding:"omitempty,min=1"`       // 页码
	Size           int      `form:"size" binding:"omitempty,min=1"`       // 每页数量
	StoreIDs       []string `form:"store_ids" binding:"omitempty"`        // 门店ID列表（可选，多选）
	MerchantName   string   `form:"merchant_name" binding:"omitempty"`    // 品牌商支付商户名称（可选，模糊匹配）
	CreatedAtStart string   `form:"created_at_start" binding:"omitempty"` // 创建时间开始（可选）
	CreatedAtEnd   string   `form:"created_at_end" binding:"omitempty"`   // 创建时间结束（可选）
}
