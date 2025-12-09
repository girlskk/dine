package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type ReconciliationListReq struct {
	Page    int                     `json:"page"`
	Size    int                     `json:"size"`
	StartAt util.RequestDate        `json:"start_at"`                                                                // 开始日期
	EndAt   util.RequestDate        `json:"end_at"`                                                                  // 截止日期
	StoreID int                     `json:"store_id"`                                                                // 门店ID
	Channel domain.OrderPaidChannel `json:"channel" binding:"omitempty,oneof=cash wechat alipay point point_wallet"` // 支付渠道
}

type ReconciliationListExportReq struct {
	StartAt util.RequestDate        `json:"start_at"`                                                                // 开始日期
	EndAt   util.RequestDate        `json:"end_at"`                                                                  // 截止日期
	StoreID int                     `json:"store_id"`                                                                // 门店ID
	Channel domain.OrderPaidChannel `json:"channel" binding:"omitempty,oneof=cash wechat alipay point point_wallet"` // 支付渠道
}

type ReconciliationDetailReq struct {
	ID int `json:"id" binding:"required"` // 财务对账单ID
}

type ReconciliationSummaryReq struct {
	StartAt util.RequestDate        `json:"start_at"`                                                                // 开始日期
	EndAt   util.RequestDate        `json:"end_at"`                                                                  // 截止日期
	Channel domain.OrderPaidChannel `json:"channel" binding:"omitempty,oneof=cash wechat alipay point point_wallet"` // 支付渠道
	StoreID int                     `json:"store_id"`                                                                // 门店ID
}

type PointSettlementListReq struct {
	Page    int              `json:"page"`
	Size    int              `json:"size"`
	StartAt util.RequestDate `json:"start_at"` // 开始日期
	EndAt   util.RequestDate `json:"end_at"`   // 截止日期
	StoreID int              `json:"store_id"` // 门店ID
}

type PointSettlementListExportReq struct {
	StartAt util.RequestDate `json:"start_at"` // 开始日期
	EndAt   util.RequestDate `json:"end_at"`   // 截止日期
	StoreID int              `json:"store_id"` // 门店ID
}

type PointSettlementIDReq struct {
	ID int `json:"id"` // 积分结算账单ID
}

type StoreWithdrawListReq struct {
	Page    int                        `json:"page"`
	Size    int                        `json:"size"`
	StartAt util.RequestDate           `json:"start_at"`                                 // 开始日期
	EndAt   util.RequestDate           `json:"end_at"`                                   // 截止日期
	Status  domain.StoreWithdrawStatus `json:"status" binding:"omitempty,oneof=1 2 3 4"` // 提现状态：1-待审核 2-已审核 3-已驳回 4-待提交
	StoreID int                        `json:"store_id"`                                 // 门店ID
}

type StoreWithdrawIDReq struct {
	ID int `json:"id" binding:"required"` // 提现单ID
}

type ReconciliationGenerateDailyRecordsByDateReq struct {
	RecordDate util.RequestDate `json:"record_date"` // 账单日期
}
