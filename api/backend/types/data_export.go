package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// DataExportListReq 数据导出列表请求
type DataExportListReq struct {
	upagination.RequestPagination
	Type           domain.DataExportType   `json:"type" binding:"omitempty,oneof=order_list"`               // 导出类型: order_list 订单列表
	Status         domain.DataExportStatus `json:"status" binding:"omitempty,oneof=pending success failed"` // 导出状态
	CreatedAtStart util.RequestDate        `json:"created_at_start"`                                        // 创建时间开始日期（格式：YYYY-MM-DD）
	CreatedAtEnd   util.RequestDate        `json:"created_at_end"`                                          // 创建时间结束日期（格式：YYYY-MM-DD）
}

// DataExportListResp 数据导出列表响应
type DataExportListResp struct {
	DataExports []*domain.DataExport `json:"data_exports"` // 数据导出列表
	Total       int                  `json:"total"`        // 总数
}

// DataExportRetryReq 重新导出请求
type DataExportRetryReq struct {
	ID int `json:"id" binding:"required"` // 数据导出ID
}
