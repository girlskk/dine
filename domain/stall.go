package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrStallNotExists  = errors.New("出品部门不存在")
	ErrStallNameExists = errors.New("出品部门名称已存在")
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/stall_repository.go -package=mock . StallRepository
type StallRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (stall *Stall, err error)
	Create(ctx context.Context, stall *Stall) (err error)
	Update(ctx context.Context, stall *Stall) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetStalls(ctx context.Context, pager *upagination.Pagination, filter *StallListFilter, orderBys ...StallOrderBy) (stalls []*Stall, total int, err error)
	Exists(ctx context.Context, params StallExistsParams) (exists bool, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/stall_interactor.go -package=mock . StallInteractor
type StallInteractor interface {
	Create(ctx context.Context, stall *Stall) (err error)
	Update(ctx context.Context, stall *Stall) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetStall(ctx context.Context, id uuid.UUID) (*Stall, error)
	GetStalls(ctx context.Context, pager *upagination.Pagination, filter *StallListFilter, orderBys ...StallOrderBy) (stalls []*Stall, total int, err error)
	StallSimpleUpdate(ctx context.Context, updateField StallSimpleUpdateType, stall *Stall) (err error)
}

type StallSimpleUpdateType string

const (
	StallSimpleUpdateTypeEnabled StallSimpleUpdateType = "enabled"
)

type StallOrderByType int

const (
	_ StallOrderByType = iota
	StallOrderByID
	StallOrderByCreatedAt
	StallOrderBySortOrder
)

type StallOrderBy struct {
	OrderBy StallOrderByType
	Desc    bool
}

func NewStallOrderByID(desc bool) StallOrderBy {
	return StallOrderBy{
		OrderBy: StallOrderByID,
		Desc:    desc,
	}
}

func NewStallOrderByCreatedAt(desc bool) StallOrderBy {
	return StallOrderBy{
		OrderBy: StallOrderByCreatedAt,
		Desc:    desc,
	}
}

func NewStallOrderBySortOrder(desc bool) StallOrderBy {
	return StallOrderBy{
		OrderBy: StallOrderBySortOrder,
		Desc:    desc,
	}
}

type StallPrintType string

const (
	StallPrintTypeReceipt StallPrintType = "receipt" // 小票/收据
	StallPrintTypeLabel   StallPrintType = "label"   // 标签
)

func (StallPrintType) Values() []string {
	return []string{string(StallPrintTypeReceipt), string(StallPrintTypeLabel)}
}

type StallType string

const (
	StallTypeSystem StallType = "system" // 系统出品部门
	StallTypeBrand  StallType = "brand"  // 品牌出品部门
	StallTypeStore  StallType = "store"  // 门店出品部门
)

func (StallType) Values() []string {
	return []string{string(StallTypeSystem), string(StallTypeBrand), string(StallTypeStore)}
}

type Stall struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`        // 出品部门名称
	StallType  StallType      `json:"stall_type"`  // system/brand/store
	PrintType  StallPrintType `json:"print_type"`  // 打印类型 ：receipt/label
	Enabled    bool           `json:"enabled"`     // 使用状态
	SortOrder  int            `json:"sort_order"`  // 排序
	MerchantID uuid.UUID      `json:"merchant_id"` // 商户 ID
	StoreID    uuid.UUID      `json:"store_id"`    // 门店 ID
	CreatedAt  time.Time      `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time      `json:"updated_at"`  // 更新时间
}

// StallExistsParams 存在性检查参数
type StallExistsParams struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	Name       string
	ExcludeID  uuid.UUID
}

// StallListFilter 查询过滤参数
type StallListFilter struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	StallType  StallType
	PrintType  StallPrintType
	Enabled    *bool
	Name       string
}
