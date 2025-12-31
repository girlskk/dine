package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ------------------------------------------------------------
// 错误定义
// ------------------------------------------------------------

var ()

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// ProfitDistributionBillStatus 分账账单状态
type ProfitDistributionBillStatus string

const (
	ProfitDistributionBillStatusUnpaid ProfitDistributionBillStatus = "unpaid" // 未打款
	ProfitDistributionBillStatusPaid   ProfitDistributionBillStatus = "paid"   // 已打款
)

func (ProfitDistributionBillStatus) Values() []string {
	return []string{
		string(ProfitDistributionBillStatusUnpaid),
		string(ProfitDistributionBillStatusPaid),
	}
}

// ProfitDistributionConfig 分账任务配置
type ProfitDistributionConfig struct {
	TaskHour   int // 定时任务执行小时
	TaskMinute int // 定时任务执行分钟
}

// ProfitDistributionBillInteractor 分账账单用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/profit_distribution_bill_interactor.go -package=mock . ProfitDistributionBillInteractor
type ProfitDistributionBillInteractor interface {
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProfitDistributionBillSearchParams) (*ProfitDistributionBillSearchRes, error)
	GenerateProfitDistributionBills(ctx context.Context) error
}

// ------------------------------------------------------------
// 实体定义
// ------------------------------------------------------------

// ProfitDistributionBill 分账账单实体
type ProfitDistributionBill struct {
	ID               uuid.UUID                       `json:"id"`                // 分账账单ID
	No               string                          `json:"no"`                // 分账账单编号
	MerchantID       uuid.UUID                       `json:"merchant_id"`       // 品牌商ID
	StoreID          uuid.UUID                       `json:"store_id"`          // 门店ID
	RevenueID        uuid.UUID                       `json:"revenue_id"`        // 门店营业额ID
	ReceivableAmount decimal.Decimal                 `json:"receivable_amount"` // 应收金额（令吉）
	PaymentAmount    decimal.Decimal                 `json:"payment_amount"`    // 打款金额（令吉）
	Status           ProfitDistributionBillStatus    `json:"status"`            // 分账状态
	BillDate         time.Time                       `json:"bill_date"`         // 账单日期
	StartDate        time.Time                       `json:"start_date"`        // 账单周期：开始日期
	EndDate          time.Time                       `json:"end_date"`          // 账单周期：结束日期
	RuleSnapshot     *ProfitDistributionRuleSnapshot `json:"rule_snapshot"`     // 分账方案快照
	CreatedAt        time.Time                       `json:"created_at"`        // 创建时间
	UpdatedAt        time.Time                       `json:"updated_at"`        // 更新时间
}

// ProfitDistributionRuleSnapshot 分账方案快照（用于账单历史追溯）
type ProfitDistributionRuleSnapshot struct {
	RuleID     uuid.UUID       `json:"rule_id"`     // 分账方案ID
	RuleName   string          `json:"rule_name"`   // 分账方案名称
	SplitRatio decimal.Decimal `json:"split_ratio"` // 分账比例
}

// ProfitDistributionBills 分账账单集合
type ProfitDistributionBills []*ProfitDistributionBill

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// ProfitDistributionBillSearchParams 查询参数
type ProfitDistributionBillSearchParams struct {
	MerchantID    uuid.UUID
	StoreIDs      []uuid.UUID
	BillStartDate *time.Time
	BillEndDate   *time.Time
	Status        ProfitDistributionBillStatus
}

// ProfitDistributionBillSearchRes 查询结果
type ProfitDistributionBillSearchRes struct {
	*upagination.Pagination
	Items ProfitDistributionBills `json:"items"`
}
