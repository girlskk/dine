package types

import (
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// TaxFeeCreateReq 创建税费
type TaxFeeCreateReq struct {
	Name    string          `json:"name" binding:"required,max=50"` // 税费名称
	TaxRate decimal.Decimal `json:"tax_rate" binding:"required"`    // 税率 tax_rate 示例：6% -> 0.06
}

// TaxFeeUpdateReq 更新税费（仅可修改部分字段）
type TaxFeeUpdateReq struct {
	Name    string          `json:"name" binding:"required,max=50"` // 税费名称
	TaxRate decimal.Decimal `json:"tax_rate" binding:"required"`    // 税率 tax_rate 示例：6% -> 0.06
}

// TaxFeeListReq 列表查询
type TaxFeeListReq struct {
	upagination.RequestPagination
	Name string `form:"name"`
}

type TaxFeeListResp struct {
	TaxFees []*domain.TaxFee `json:"tax_fees"` // 税费列表
	Total   int              `json:"total"`    // 总数
}
