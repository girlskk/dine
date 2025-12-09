package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrStoreWithdrawAmountNotEnough = errors.New("可提现金额不足")
	ErrStoreWithdrawNotExists       = errors.New("提现单不存在")
	ErrStoreWithdrawStatusInvalid   = errors.New("提现单状态错误")
)

// 门店提现单
type StoreWithdraw struct {
	ID                  int                 `json:"id"`                    // 提现单ID
	StoreID             int                 `json:"store_id"`              // 门店ID
	StoreName           string              `json:"store_name"`            // 门店名称
	No                  string              `json:"no"`                    // 单据编号
	Amount              decimal.Decimal     `json:"amount"`                // 提现金额（原始申请金额）
	PointWithdrawalRate decimal.Decimal     `json:"point_withdrawal_rate"` // 积分提现费率
	ActualAmount        decimal.Decimal     `json:"actual_amount"`         // 实际到账金额（扣除平台佣金）
	AccountType         AccountType         `json:"account_type"`          // 账户类型：public-对公 private-对私
	BankAccount         string              `json:"bank_account"`          // 银行账号
	BankCardName        string              `json:"bank_card_name"`        // 银行卡名称（对公时为公司名称）
	BankName            string              `json:"bank_name"`             // 银行名称
	BankBranch          string              `json:"bank_branch"`           // 开户支行
	InvoiceAmount       decimal.Decimal     `json:"invoice_amount"`        // 开票金额
	Status              StoreWithdrawStatus `json:"status"`                // 提现状态：1-待审核 2-已审核 3-已驳回 4-待提交
	CreatedAt           time.Time           `json:"created_at"`            // 创建时间
	UpdatedAt           time.Time           `json:"updated_at"`            // 更新时间
}

type StoreWithdraws []*StoreWithdraw

type AccountType string

const (
	AccountTypePublic  AccountType = "public"  // 对公账户
	AccountTypePrivate AccountType = "private" // 对私账户
)

func (AccountType) Values() []string {
	return []string{
		string(AccountTypePublic),
		string(AccountTypePrivate),
	}
}

type StoreWithdrawStatus int

const (
	_                              StoreWithdrawStatus = iota
	StoreWithdrawStatusPending                         // 待审核
	StoreWithdrawStatusApproved                        // 已审核
	StoreWithdrawStatusRejected                        // 已驳回
	StoreWithdrawStatusUncommitted                     // 待提交
)

// 提现单编号前缀
const (
	StoreWithdrawNoPrefix string = "TX"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_withdraw_interactor.go -package=mock . StoreWithdrawInteractor
type StoreWithdrawInteractor interface {
	Apply(ctx context.Context, withdraw *StoreWithdraw) error
	Update(ctx context.Context, withdraw *StoreWithdraw) error
	Commit(ctx context.Context, id int, storeID int) error
	Cancel(ctx context.Context, id int, storeID int) error
	Delete(ctx context.Context, id int, storeID int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params StoreWithdrawSearchParams) (*StoreWithdrawSearchRes, error)
	Approve(ctx context.Context, id int) error
	Reject(ctx context.Context, id int) error
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_withdraw_repository.go -package=mock . StoreWithdrawRepository
type StoreWithdrawRepository interface {
	Create(ctx context.Context, withdraw *StoreWithdraw) error
	Update(ctx context.Context, withdraw *StoreWithdraw) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params StoreWithdrawSearchParams) (*StoreWithdrawSearchRes, error)
	FindByIDForUpdate(ctx context.Context, id int) (*StoreWithdraw, error)
	FindByID(ctx context.Context, id int) (*StoreWithdraw, error)
	Delete(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, status StoreWithdrawStatus) error
}

type StoreWithdrawSearchParams struct {
	StartAt *time.Time          `json:"start_at"` // 开始日期
	EndAt   *time.Time          `json:"end_at"`   // 截止日期
	StoreID int                 `json:"store_id"` // 门店ID
	Status  StoreWithdrawStatus `json:"status"`   // 提现状态
}

type StoreWithdrawSearchRes struct {
	*upagination.Pagination
	Items StoreWithdraws `json:"items"`
}
