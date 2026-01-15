package types

import (
	"github.com/shopspring/decimal"
)

// ProfitDistributionBillListReq 分账账单列表请求
type ProfitDistributionBillListReq struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`               // 页码
	Size          int      `form:"size" binding:"omitempty,min=1"`               // 每页数量
	BillStartDate string   `form:"bill_start_date" binding:"omitempty"`          // 账单开始日期
	BillEndDate   string   `form:"bill_end_date" binding:"omitempty"`            // 账单结束日期
	StoreIDs      []string `form:"store_ids" binding:"omitempty"`                // 门店ID列表（可选，多选）
	Status        string   `form:"status" binding:"omitempty,oneof=unpaid paid"` // 分账状态（可选）
}

// ProfitDistributionBillPayReq 打款分账账单请求
type ProfitDistributionBillPayReq struct {
	PaymentAmount decimal.Decimal `json:"payment_amount" binding:"required"` // 打款金额（必选，单位：令吉）
}
