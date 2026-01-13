package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// SalesReportReq 销售报表请求
type SalesReportReq struct {
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

// ProductSalesSummaryReq 商品销售汇总请求
type ProductSalesSummaryReq struct {
	StoreIDs          string `form:"store_ids"`                                              // 门店ID列表（逗号分隔）
	BusinessDateStart string `form:"business_date_start" binding:"required"`                 // 营业日开始
	BusinessDateEnd   string `form:"business_date_end" binding:"required"`                   // 营业日结束
	OrderChannel      string `form:"order_channel"`                                          // 订单来源
	CategoryID        string `form:"category_id"`                                            // 商品分类ID
	ProductName       string `form:"product_name"`                                           // 商品名称（模糊搜索）
	ProductType       string `form:"product_type" binding:"omitempty,oneof=normal set_meal"` // 商品类型

	upagination.RequestPagination
}

// ProductSalesSummaryResp 商品销售汇总响应
type ProductSalesSummaryResp struct {
	Items      []*domain.ProductSalesSummaryItem `json:"items"`      // 汇总数据
	Pagination *upagination.Pagination           `json:"pagination"` // 分页信息
}
