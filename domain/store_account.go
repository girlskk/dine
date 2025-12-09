package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrStoreAccountNotExists = errors.New("门店账户不存在")
)

// 门店账户
type StoreAccount struct {
	ID              int             `json:"id"`
	StoreID         int             `json:"store_id"`         // 关联门店
	Balance         decimal.Decimal `json:"balance"`          // 账户余额（可提现金额）
	PendingWithdraw decimal.Decimal `json:"pending_withdraw"` // 待提现金额
	Withdrawn       decimal.Decimal `json:"withdrawn"`        // 已提现金额
	TotalAmount     decimal.Decimal `json:"total_amount"`     // 总收益金额
	CreatedAt       time.Time       `json:"created_at"`       // 创建时间
	UpdatedAt       time.Time       `json:"updated_at"`       // 更新时间
}

type StoreAccountAdjustments struct {
	BalanceDelta         decimal.Decimal // 账户余额增减量
	PendingWithdrawDelta decimal.Decimal // 待提现金额增减量
	WithdrawnDelta       decimal.Decimal // 已提现金额增减量
	TotalDelta           decimal.Decimal // 总收益增减量
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_account_repository.go -package=mock . StoreAccountRepository
type StoreAccountRepository interface {
	Create(ctx context.Context, storeAccount *StoreAccount) error
	FindByStoreForUpdate(ctx context.Context, storeID int) (*StoreAccount, error)
	FindByStore(ctx context.Context, storeID int) (*StoreAccount, error)
	// 金额变更，增量数据
	AdjustAmount(ctx context.Context, storeID int, adjustments StoreAccountAdjustments) error
	// 记录资金流水
	RecordTransaction(ctx context.Context, tx *StoreAccountTransaction) error
	// 分页查询资金流水
	PagedListTransactions(ctx context.Context, page *upagination.Pagination, params StoreAccountTransactionSearchParams) (*StoreAccountTransactionSearchRes, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_account_interactor.go -package=mock . StoreAccountInteractor
type StoreAccountInteractor interface {
	GetDetail(ctx context.Context, storeID int) (*StoreAccount, error)
	PagedListTransactions(ctx context.Context, page *upagination.Pagination, params StoreAccountTransactionSearchParams) (*StoreAccountTransactionSearchRes, error)
}

type StoreAccountTransactionSearchParams struct {
	StartAt *time.Time
	EndAt   *time.Time
	StoreID int
	Type    TransactionType
}

type StoreAccountTransactionSearchRes struct {
	*upagination.Pagination
	Items StoreAccountTransactions `json:"items"`
}
