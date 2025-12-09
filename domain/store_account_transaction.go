package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// 门店账户流水
type StoreAccountTransaction struct {
	ID        int             `json:"id"`             // 流水ID
	StoreID   int             `json:"store_id"`       // 门店ID
	No        string          `json:"no"`             // 单据编号
	Amount    decimal.Decimal `json:"amount"`         // 变动金额
	After     decimal.Decimal `json:"after"`          // 变动后金额
	Type      TransactionType `db:"type" json:"type"` // 变动类型：1-销售进账 2-进账撤回 3-申请提现 4-提现通过 5-提现驳回 6-提现撤回
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type StoreAccountTransactions []*StoreAccountTransaction

type TransactionType int // 流水类型

const (
	_                              TransactionType = iota
	TransactionTypeSaleIncome                      // 销售进账
	TransactionTypeSaleRevert                      // 销售撤回
	TransactionTypeWithdrawApply                   // 申请提现
	TransactionTypeWithdrawApprove                 // 提现通过
	TransactionTypeWithdrawReject                  // 提现驳回
	TransactionTypeWithdrawCancel                  // 提现撤回
)
