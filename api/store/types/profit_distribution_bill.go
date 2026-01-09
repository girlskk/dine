package types

// ProfitDistributionBillListReq 分账账单列表请求
type ProfitDistributionBillListReq struct {
	Page          int    `form:"page" binding:"omitempty,min=1"`                // 页码
	Size          int    `form:"size" binding:"omitempty,min=1"`                // 每页数量
	BillStartDate string `form:"bill_start_date" binding:"omitempty"`           // 账单开始日期
	BillEndDate   string `form:"bill_end_date" binding:"omitempty"`             // 账单结束日期
	Status        string `form:"status" binding:"omitempty,oneof=pending paid"` // 分账状态（可选）
}
