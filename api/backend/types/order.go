package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// SalesReportReq 销售报表请求
type SalesReportReq struct {
	MerchantID        string `form:"merchant_id" binding:"required"`         // 品牌商ID
	StoreIDs          string `form:"store_ids"`                              // 门店ID列表（逗号分隔）
	BusinessDateStart string `form:"business_date_start" binding:"required"` // 营业日开始
	BusinessDateEnd   string `form:"business_date_end" binding:"required"`   // 营业日结束

	upagination.RequestPagination
}

// SalesReportResp 销售报表响应
type SalesReportResp struct {
	Items      []*domain.OrderSalesReportItem `json:"items"`      // 报表数据
	Pagination *upagination.Pagination        `json:"pagination"` // 分页信息
}
