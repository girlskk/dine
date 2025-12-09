package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrReconciliationRecordNotExists = errors.New("财务对账单不存在")
)

const (
	ReconciliationRecordListExportSingleMaxSize   = 3000
	ReconciliationRecordDetailExportSingleMaxSize = 30
)

// 财务对账单
type ReconciliationRecord struct {
	ID         int              `json:"id"`
	No         string           `json:"no"`          // 单号
	StoreID    int              `json:"store_id"`    // 门店ID
	StoreName  string           `json:"store_name"`  // 门店名称
	OrderCount int              `json:"order_count"` // 入账笔数
	Amount     decimal.Decimal  `json:"amount"`      // 入账金额
	Channel    OrderPaidChannel `json:"channel"`     // 订单支付渠道
	Date       time.Time        `json:"date"`        // 账单日期
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type ReconciliationRecords []*ReconciliationRecord

// 财务对账单明细
type ReconciliationDetail struct {
	OrderID         int             `json:"order_id"`          // 订单ID
	OrderNo         string          `json:"order_no"`          // 订单号
	OrderFinishedAt time.Time       `json:"order_finished_at"` // 订单完成时间
	OrderAmount     decimal.Decimal `json:"order_amount"`      // 订单金额
	Amount          decimal.Decimal `json:"amount"`            // 入账金额
	ProductInfo     []string        `json:"product_info"`      // 订单商品信息
}

type ReconciliationDetails []*ReconciliationDetail

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/reconciliation_record_repository.go -package=mock . ReconciliationRecordRepository
type ReconciliationRecordRepository interface {
	BatchCreate(ctx context.Context, records ReconciliationRecords) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ReconciliationSearchParams) (*ReconciliationSearchRes, error)
	FindByID(ctx context.Context, id int) (*ReconciliationRecord, error)
	GetReconciliationRange(ctx context.Context, params ReconciliationSearchParams) (ReconciliationRange, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/reconciliation_record_interactor.go -package=mock . ReconciliationRecordInteractor
type ReconciliationRecordInteractor interface {
	GenerateDailyRecords(ctx context.Context) (err error)                             // 生成每天的财务对账单
	GenerateDailyRecordsByDate(ctx context.Context, recordDate time.Time) (err error) // 生成指定日期的财务对账单
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ReconciliationSearchParams) (*ReconciliationSearchRes, error)
	ListDetails(ctx context.Context, id, storeID int) (ReconciliationDetails, error)
	Summary(ctx context.Context, params ReconciliationSearchParams) (*ReconciliationSummaryRes, error)
	GetReconciliationRange(ctx context.Context, params ReconciliationSearchParams) (ReconciliationRange, error)
}

type ReconciliationSearchParams struct {
	StartAt *time.Time       `json:"start_at"` // 开始日期
	EndAt   *time.Time       `json:"end_at"`   // 截止日期
	StoreID int              `json:"store_id"` // 门店ID
	Channel OrderPaidChannel `json:"channel"`  // 支付渠道
	IDGte   int              `json:"id_gte"`
	IDLte   int              `json:"id_lte"`
}

type ReconciliationSearchRes struct {
	*upagination.Pagination
	Items ReconciliationRecords `json:"items"`
}

type ReconciliationSummaryRes struct {
	TotalCount        int             `json:"total_count"`         // 总交易笔数
	TotalAmount       decimal.Decimal `json:"total_amount"`        // 总金额
	WechatAmount      decimal.Decimal `json:"wechat_amount"`       // 微信支付金额
	AlipayAmount      decimal.Decimal `json:"alipay_amount"`       // 支付宝支付金额
	CashAmount        decimal.Decimal `json:"cash_amount"`         // 现金支付金额
	PointAmount       decimal.Decimal `json:"point_amount"`        // 积分支付金额
	PointWalletAmount decimal.Decimal `json:"point_wallet_amount"` // 知心话钱包支付金额
}

type ReconciliationRange struct {
	MinID int
	MaxID int
	Count int
}

type ReconciliationRecordListExportParams struct {
	Filter ReconciliationSearchParams `json:"filter"`
	Pager  upagination.Pagination     `json:"pager"`
}

// 财务对账单导出明细列
type ReconciliationDetailExportColumn struct {
	RecordNo        string           `json:"record_no"`         // 单据编号
	RecordDate      time.Time        `json:"record_date"`       // 账单日期
	StoreName       string           `json:"store_name"`        // 门店名称
	Channel         OrderPaidChannel `json:"channel"`           // 支付渠道
	OrderNo         string           `json:"order_no"`          // 订单号
	OrderFinishedAt time.Time        `json:"order_finished_at"` // 订单完成时间
	OrderAmount     decimal.Decimal  `json:"order_amount"`      // 订单金额
}
