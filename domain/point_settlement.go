package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrPointSettlementNotExists     = errors.New("积分结算单不存在")
	ErrPointSettlementStatusInvalid = errors.New("积分结算单状态错误")
)

const (
	PointSettlementListExportSingleMaxSize   = 3000
	PointSettlementDetailExportSingleMaxSize = 30
)

// 积分结算单
type PointSettlement struct {
	ID                  int                   `json:"id"`
	No                  string                `json:"no"`                    // 单号
	StoreID             int                   `json:"store_id"`              // 门店ID
	StoreName           string                `json:"store_name"`            // 门店名称
	OrderCount          int                   `json:"order_count"`           // 入账笔数
	Amount              decimal.Decimal       `json:"amount"`                // 入账金额（扣除费率）
	TotalPoints         decimal.Decimal       `json:"total_points"`          // 积分总额
	Date                time.Time             `json:"date"`                  // 账单日期
	Status              PointSettlementStatus `json:"status"`                // 账单状态：1-待审核 2-已审核
	PointSettlementRate decimal.Decimal       `json:"point_settlement_rate"` // 积分结算费率（单位：百分比，例如 0.1234 表示 12.34%）
	ApprovedAt          *time.Time            `json:"approved_at"`           // 审批时间
	ApproverID          int                   `json:"approver_id"`           // 审批者ID
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

type PointSettlementStatus int

const (
	PointSettlementStatusPending  PointSettlementStatus = iota + 1 // 待审核
	PointSettlementStatusApproved                                  // 已审核
)

type PointSettlements []*PointSettlement

// 积分结算单单明细
type PointSettlementDetail struct {
	OrderID         int             `json:"order_id"`          // 订单ID
	OrderNo         string          `json:"order_no"`          // 订单号
	OrderFinishedAt time.Time       `json:"order_finished_at"` // 订单完成时间
	OrderAmount     decimal.Decimal `json:"order_amount"`      // 订单金额
	Amount          decimal.Decimal `json:"amount"`            // 入账金额
	ProductInfo     []string        `json:"product_info"`      // 订单商品信息
	MemberName      string          `json:"member_name"`       // 会员姓名
}

type PointSettlementDetails []*PointSettlementDetail

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/point_settlement_repository.go -package=mock . PointSettlementRepository
type PointSettlementRepository interface {
	BatchCreate(ctx context.Context, pointSettlements PointSettlements) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params PointSettlementSearchParams) (*PointSettlementSearchRes, error)
	FindByID(ctx context.Context, id int) (*PointSettlement, error)
	FindByIDForUpdate(ctx context.Context, id int) (*PointSettlement, error)
	UpdateStatus(ctx context.Context, id int, status PointSettlementStatus, approverID *int) error
	GetPointSettlementRange(ctx context.Context, params PointSettlementSearchParams) (PointSettlementRange, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/point_settlement_interactor.go -package=mock .	 PointSettlementInteractor
type PointSettlementInteractor interface {
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params PointSettlementSearchParams) (*PointSettlementSearchRes, error)
	ListDetails(ctx context.Context, id, storeID int) (PointSettlementDetails, error)
	Approve(ctx context.Context, id int) error
	UnApprove(ctx context.Context, id int) error
	GetPointSettlementRange(ctx context.Context, params PointSettlementSearchParams) (PointSettlementRange, error)
}

type PointSettlementSearchParams struct {
	StartAt *time.Time `json:"start_at"` // 开始日期
	EndAt   *time.Time `json:"end_at"`   // 截止日期
	StoreID int        `json:"store_id"` // 门店ID
	IDGte   int        `json:"id_gte"`
	IDLte   int        `json:"id_lte"`
}

type PointSettlementSearchRes struct {
	*upagination.Pagination
	Items PointSettlements `json:"items"`
}

type PointSettlementRange struct {
	MinID int
	MaxID int
	Count int
}

type PointSettlementListExportParams struct {
	Filter PointSettlementSearchParams `json:"filter"`
	Pager  upagination.Pagination      `json:"pager"`
}

// 积分结算单导出明细列
type PointSettlementDetailExportColumn struct {
	RecordNo        string          `json:"record_no"`         // 单据编号
	RecordDate      time.Time       `json:"record_date"`       // 账单日期
	StoreName       string          `json:"store_name"`        // 门店名称
	OrderNo         string          `json:"order_no"`          // 订单号
	OrderFinishedAt time.Time       `json:"order_finished_at"` // 订单完成时间
	OrderAmount     decimal.Decimal `json:"order_amount"`      // 订单金额
}
