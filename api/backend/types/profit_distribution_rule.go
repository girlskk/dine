package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ProfitDistributionRuleCreateReq 创建分账方案请求
type ProfitDistributionRuleCreateReq struct {
	Name          string                                    `json:"name" binding:"required,max=255"`                      // 分账方案名称（必选）
	StoreIDs      []uuid.UUID                               `json:"store_ids" binding:"required,min=1,dive,uuid"`         // 门店ID列表（必选，多选）
	SplitRatio    decimal.Decimal                           `json:"split_ratio" binding:"required"`                       // 分账比例（必选，0-1之间）
	BillingCycle  domain.ProfitDistributionRuleBillingCycle `json:"billing_cycle" binding:"required,oneof=daily monthly"` // 账单生成周期（必选）
	EffectiveDate time.Time                                 `json:"effective_date" binding:"required"`                    // 方案生效日期（必选）
	ExpiryDate    time.Time                                 `json:"expiry_date" binding:"required"`                       // 方案失效日期（必选）
}

// ProfitDistributionRuleUpdateReq 更新分账方案请求
type ProfitDistributionRuleUpdateReq ProfitDistributionRuleCreateReq

// ProfitDistributionRuleListReq 分账方案列表请求
type ProfitDistributionRuleListReq struct {
	upagination.RequestPagination
	Name   string                              `json:"name" form:"name"`     // 分账方案名称（模糊匹配）
	Status domain.ProfitDistributionRuleStatus `json:"status" form:"status"` // 状态筛选（可选）
}
