package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ------------------------------------------------------------
// 错误定义
// ------------------------------------------------------------

var (
	ErrProfitDistributionRuleNotExists         = errors.New("分账方案不存在")
	ErrProfitDistributionRuleNameExists        = errors.New("分账方案名称已存在")
	ErrProfitDistributionRuleStoreBound        = errors.New("门店已绑定其他分账方案")
	ErrProfitDistributionRuleStatusInvalid     = errors.New("分账方案状态无效")
	ErrProfitDistributionRuleDateInvalid       = errors.New("方案生效日期必须必须晚于今天，生效日期要早于失效日期")
	ErrProfitDistributionRuleSplitRatioInvalid = errors.New("分账比例必须为0-1之间的数值")
)

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// ProfitDistributionRuleStatus 分账方案状态
type ProfitDistributionRuleStatus string

const (
	ProfitDistributionRuleStatusEnabled  ProfitDistributionRuleStatus = "enabled"  // 启用
	ProfitDistributionRuleStatusDisabled ProfitDistributionRuleStatus = "disabled" // 禁用
)

func (ProfitDistributionRuleStatus) Values() []string {
	return []string{
		string(ProfitDistributionRuleStatusEnabled),
		string(ProfitDistributionRuleStatusDisabled),
	}
}

// ProfitDistributionRuleBillingCycle 账单生成周期
type ProfitDistributionRuleBillingCycle string

const (
	ProfitDistributionRuleBillingCycleDaily   ProfitDistributionRuleBillingCycle = "daily"   // 按日
	ProfitDistributionRuleBillingCycleMonthly ProfitDistributionRuleBillingCycle = "monthly" // 按月
)

func (ProfitDistributionRuleBillingCycle) Values() []string {
	return []string{
		string(ProfitDistributionRuleBillingCycleDaily),
		string(ProfitDistributionRuleBillingCycleMonthly),
	}
}

// ------------------------------------------------------------
// 实体定义
// ------------------------------------------------------------

// ProfitDistributionRule 分账方案实体
type ProfitDistributionRule struct {
	ID                uuid.UUID                          `json:"id"`                  // 分账方案ID
	MerchantID        uuid.UUID                          `json:"merchant_id"`         // 品牌商ID
	Name              string                             `json:"name"`                // 分账方案名称
	SplitRatio        decimal.Decimal                    `json:"split_ratio"`         // 分账比例（0-1，单位：小数）
	BillingCycle      ProfitDistributionRuleBillingCycle `json:"billing_cycle"`       // 账单生成周期
	EffectiveDate     time.Time                          `json:"effective_date"`      // 方案生效日期
	ExpiryDate        time.Time                          `json:"expiry_date"`         // 方案失效日期
	BillGenerationDay int                                `json:"bill_generation_day"` // 账单生成日：1-28号
	Status            ProfitDistributionRuleStatus       `json:"status"`              // 状态
	StoreCount        int                                `json:"store_count"`         // 关联门店数量
	CreatedAt         time.Time                          `json:"created_at"`          // 创建时间
	UpdatedAt         time.Time                          `json:"updated_at"`          // 更新时间

	// 关联信息
	Stores []*StoreSimple `json:"stores,omitempty"` // 关联门店列表
}

// ProfitDistributionRules 分账方案集合
type ProfitDistributionRules []*ProfitDistributionRule

// ------------------------------------------------------------
// 仓储和用例接口
// ------------------------------------------------------------

// ProfitDistributionRuleRepository 分账方案仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/profit_distribution_rule_repository.go -package=mock . ProfitDistributionRuleRepository
type ProfitDistributionRuleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*ProfitDistributionRule, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*ProfitDistributionRule, error)
	Create(ctx context.Context, rule *ProfitDistributionRule) error
	Update(ctx context.Context, rule *ProfitDistributionRule) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params ProfitDistributionRuleExistsParams) (bool, error)
	CheckStoreBound(ctx context.Context, storeIDs []uuid.UUID, excludeRuleID uuid.UUID) (bool, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProfitDistributionRuleSearchParams) (*ProfitDistributionRuleSearchRes, error)
	ListAllEnabled(ctx context.Context) (ProfitDistributionRules, error)
}

// ProfitDistributionRuleInteractor 分账方案用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/profit_distribution_rule_interactor.go -package=mock . ProfitDistributionRuleInteractor
type ProfitDistributionRuleInteractor interface {
	Create(ctx context.Context, rule *ProfitDistributionRule, user User) error
	Update(ctx context.Context, rule *ProfitDistributionRule, user User) error
	Delete(ctx context.Context, id uuid.UUID, user User) error
	Enable(ctx context.Context, id uuid.UUID, user User) error
	Disable(ctx context.Context, id uuid.UUID, user User) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProfitDistributionRuleSearchParams) (*ProfitDistributionRuleSearchRes, error)
}

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// ProfitDistributionRuleExistsParams 存在性检查参数
type ProfitDistributionRuleExistsParams struct {
	MerchantID uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// ProfitDistributionRuleSearchParams 查询参数
type ProfitDistributionRuleSearchParams struct {
	MerchantID uuid.UUID
	Name       string                       // 分账方案名称（模糊匹配）
	Status     ProfitDistributionRuleStatus // 状态筛选（可选）
}

// ProfitDistributionRuleSearchRes 查询结果
type ProfitDistributionRuleSearchRes struct {
	*upagination.Pagination
	Items ProfitDistributionRules `json:"items"`
}
