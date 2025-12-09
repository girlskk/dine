package domain

import (
	"context"
	"errors"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// TableRepository 台桌仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/table_repository.go -package=mock . TableRepository
type TableRepository interface {
	FindByID(ctx context.Context, id int) (*Table, error)
	Exists(ctx context.Context, params TableExistsParams) (bool, error)
	Create(ctx context.Context, table *Table) error
	Update(ctx context.Context, table *Table) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params TableSearchParams) (*TableSearchRes, error)
	// ListByAreaID(ctx context.Context, areaID int) (Tables, error)
	UpdateStatus(ctx context.Context, id int, status TableStatus) (bool, error)
	UpdateOrderIDAndStatusFrom(ctx context.Context, id int, orderID int, from, to TableStatus) (bool, error)
	FindWithOrder(ctx context.Context, id int) (*Table, error)
}

// TableInteractor 台桌用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/table_interactor.go -package=mock . TableInteractor
type TableInteractor interface {
	Get(ctx context.Context, id int) (*Table, error)
	Create(ctx context.Context, table *Table) error
	Update(ctx context.Context, table *Table) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params TableSearchParams) (*TableSearchRes, error)
	GetWithOrder(ctx context.Context, id int) (*Table, error)
}

var (
	ErrTableNameExists = errors.New("台桌名称已存在")
	ErrTableNotExists  = errors.New("台桌不存在")
	ErrTableIsOccupied = errors.New("台桌已被占用，无法删除")
)

type TableStatus int // 台桌状态：1-空闲 2-就餐中

const (
	_                   TableStatus = iota
	TableStatusFree                 // 空闲
	TableStatusOccupied             // 占用中
)

func (ts TableStatus) String() string {
	switch ts {
	case TableStatusFree:
		return "空闲"
	case TableStatusOccupied:
		return "就餐中"
	default:
		return "未知状态"
	}
}

// Table 台桌实体
type Table struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`       // 台桌名称（如：A1, B2）
	SeatCount int         `json:"seat_count"` // 座位数
	Status    TableStatus `json:"status"`     // 状态：1-空闲 2-就餐中
	StoreID   int         `json:"store_id"`   // 所属门店ID
	AreaID    int         `json:"area_id"`    // 台桌区域ID
	OrderID   int         `json:"order_id"`   // 订单ID
	CreatedAt time.Time   `json:"created_at"` // 创建时间
	UpdatedAt time.Time   `json:"updated_at"` // 更新时间

	Area  *TableArea `json:"area"`  // 台桌区域
	Order *Order     `json:"order"` // 台桌当前订单
}

// Tables 台桌集合
type Tables []*Table

// TableSearchParams 台桌查询参数
type TableSearchParams struct {
	StoreID int
	AreaID  int
	Name    string
	Status  TableStatus
}

type TableSearchRes struct {
	*upagination.Pagination
	Items Tables `json:"items"`
}

// TableExistsParams 台桌存在性检查参数
type TableExistsParams struct {
	StoreID int
	Name    string
}
