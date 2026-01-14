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

// ListOrderReq 订单列表请求
type ListOrderReq struct {
	BusinessDate  string `form:"business_date"`                                                        // 营业日
	OrderNo       string `form:"order_no"`                                                             // 订单号
	OrderType     string `form:"order_type" binding:"omitempty,oneof=SALE REFUND PARTIAL_REFUND"`      // 订单类型
	OrderStatus   string `form:"order_status" binding:"omitempty,oneof=PLACED COMPLETED CANCELLED"`    // 订单状态
	PaymentStatus string `form:"payment_status" binding:"omitempty,oneof=UNPAID PAYING PAID REFUNDED"` // 支付状态

	upagination.RequestPagination
}

// ListOrderResp 订单列表响应
type ListOrderResp struct {
	Items      []*domain.Order         `json:"items"`      // 订单列表
	Pagination *upagination.Pagination `json:"pagination"` // 分页信息
}