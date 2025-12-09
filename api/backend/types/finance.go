package types

import (
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type ReconciliationListReq struct {
	Page    int                     `json:"page"`
	Size    int                     `json:"size"`
	StartAt util.RequestDate        `json:"start_at"`                                                                // 开始日期
	EndAt   util.RequestDate        `json:"end_at"`                                                                  // 截止日期
	Channel domain.OrderPaidChannel `json:"channel" binding:"omitempty,oneof=cash wechat alipay point point_wallet"` // 支付渠道
}

type ReconciliationDetailReq struct {
	ID int `json:"id" binding:"required"` // 财务对账单ID
}

type ReconciliationSummaryReq struct {
	StartAt util.RequestDate        `json:"start_at"`                                                                // 开始日期
	EndAt   util.RequestDate        `json:"end_at"`                                                                  // 截止日期
	Channel domain.OrderPaidChannel `json:"channel" binding:"omitempty,oneof=cash wechat alipay point point_wallet"` // 支付渠道
}

type PointSettlementListReq struct {
	Page    int              `json:"page"`
	Size    int              `json:"size"`
	StartAt util.RequestDate `json:"start_at"` // 开始日期
	EndAt   util.RequestDate `json:"end_at"`   // 截止日期
}

type PointSettlementListExportReq struct {
	StartAt util.RequestDate `json:"start_at"` // 开始日期
	EndAt   util.RequestDate `json:"end_at"`   // 截止日期
}

type PointSettlementDetailReq struct {
	ID int `json:"id" binding:"required"` // 积分结算单ID
}

type StoreAccountTransactionListReq struct {
	Page    int                    `json:"page"`
	Size    int                    `json:"size"`
	StartAt util.RequestDate       `json:"start_at"`                                 // 开始日期
	EndAt   util.RequestDate       `json:"end_at"`                                   // 截止日期
	Type    domain.TransactionType `json:"type" binding:"omitempty,oneof=1 2 3 4 5"` // 流水类型：1-销售进账 2-销售撤回 3-申请提现 4-提现通过 5-提现驳回
}

type StoreWithdrawApplyReq struct {
	Amount        decimal.Decimal    `json:"amount"`                                                // 提现金额
	AccountType   domain.AccountType `json:"account_type" binding:"omitempty,oneof=public private"` // 账户类型：public-对公 private-对私
	BankAccount   string             `json:"bank_account"`                                          // 银行账号
	BankCardName  string             `json:"bank_card_name"`                                        // 银行卡名称
	BankName      string             `json:"bank_name"`                                             // 银行名称
	BankBranch    string             `json:"bank_branch"`                                           // 开户支行
	InvoiceAmount decimal.Decimal    `json:"invoice_amount"`                                        // 开票金额
}

type StoreWithdrawUpdateReq struct {
	ID            int                `json:"id"`                                                    // 提现单ID
	Amount        decimal.Decimal    `json:"amount"`                                                // 提现金额
	AccountType   domain.AccountType `json:"account_type" binding:"omitempty,oneof=public private"` // 账户类型：public-对公 private-对私
	BankAccount   string             `json:"bank_account"`                                          // 银行账号
	BankCardName  string             `json:"bank_card_name"`                                        // 银行卡名称
	BankName      string             `json:"bank_name"`                                             // 银行名称
	BankBranch    string             `json:"bank_branch"`                                           // 开户支行
	InvoiceAmount decimal.Decimal    `json:"invoice_amount"`                                        // 开票金额
}

type StoreWithdrawListReq struct {
	Page    int                        `json:"page"`
	Size    int                        `json:"size"`
	StartAt util.RequestDate           `json:"start_at"`                                 // 开始日期
	EndAt   util.RequestDate           `json:"end_at"`                                   // 截止日期
	Status  domain.StoreWithdrawStatus `json:"status" binding:"omitempty,oneof=1 2 3 4"` // 提现状态：1-待审核 2-已审核 3-已驳回 4-待提交
}

type StoreWithdrawIDReq struct {
	ID int `json:"id" binding:"required"` // 提现单ID
}

type ReconciliationListExportReq struct {
	StartAt util.RequestDate        `json:"start_at"`                                                                // 开始日期
	EndAt   util.RequestDate        `json:"end_at"`                                                                  // 截止日期
	Channel domain.OrderPaidChannel `json:"channel" binding:"omitempty,oneof=cash wechat alipay point point_wallet"` // 支付渠道
}
