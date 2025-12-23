package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrRemarkNotExists    = errors.New("备注不存在")
	ErrRemarkNameExists   = errors.New("备注名称已存在")
	ErrRemarkDeleteSystem = errors.New("系统备注不能删除")
)

type RemarkType string

const (
	RemarkTypeSystem RemarkType = "system"
	RemarkTypeBrand  RemarkType = "brand"
)

func (RemarkType) Values() []string {
	return []string{string(RemarkTypeSystem), string(RemarkTypeBrand)}
}

type RemarkScene string

const (
	RemarkSceneWholeOrder   RemarkScene = "whole_order"   // 整单备注
	RemarkSceneItem         RemarkScene = "item"          // 单品备注
	RemarkSceneCancelReason RemarkScene = "cancel_reason" // 退菜原因
	RemarkSceneDiscount     RemarkScene = "discount"      // 优惠原因
	RemarkSceneGift         RemarkScene = "gift"          // 赠菜原��
	RemarkSceneRebill       RemarkScene = "rebill"        // 反结账原因
	RemarkSceneRefundReject RemarkScene = "refund_reject" // 拒绝退款
)

func (RemarkScene) Values() []string {
	return []string{
		string(RemarkSceneWholeOrder),
		string(RemarkSceneItem),
		string(RemarkSceneCancelReason),
		string(RemarkSceneDiscount),
		string(RemarkSceneGift),
		string(RemarkSceneRebill),
		string(RemarkSceneRefundReject),
	}
}

type RemarkSimpleUpdateType string

const (
	RemarkSimpleUpdateTypeEnabled RemarkSimpleUpdateType = "enabled"
)

type RemarkOrderByType int

const (
	_ RemarkOrderByType = iota
	RemarkOrderByID
	RemarkOrderByCreatedAt
	RemarkOrderBySortOrder
)

type RemarkOrderBy struct {
	OrderBy RemarkOrderByType
	Desc    bool
}

func NewRemarkOrderByID(desc bool) RemarkOrderBy {
	return RemarkOrderBy{
		OrderBy: RemarkOrderByID,
		Desc:    desc,
	}
}

func NewRemarkOrderByCreatedAt(desc bool) RemarkOrderBy {
	return RemarkOrderBy{
		OrderBy: RemarkOrderByCreatedAt,
		Desc:    desc,
	}
}

func NewRemarkOrderBySortOrder(desc bool) RemarkOrderBy {
	return RemarkOrderBy{
		OrderBy: RemarkOrderBySortOrder,
		Desc:    desc,
	}
}

type Remark struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`          // 备注名称
	RemarkType   RemarkType `json:"remark_type"`   // 备注类型：系统/品牌
	Enabled      bool       `json:"enabled"`       // 是否启用
	SortOrder    int        `json:"sort_order"`    // 排序，值越小越靠前
	CategoryID   uuid.UUID  `json:"category_id"`   // 分类ID
	CategoryName string     `json:"category_name"` // 分类名称
	MerchantID   uuid.UUID  `json:"merchant_id"`   // 品牌商ID，仅品牌备注需要
	StoreID      uuid.UUID  `json:"store_id"`      // 门店 ID
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Remarks []*Remark

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/remark_repository.go -package=mock . RemarkRepository
type RemarkRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Remark, error)
	Create(ctx context.Context, remark *Remark) (err error)
	Update(ctx context.Context, remark *Remark) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetRemarks(ctx context.Context, pager *upagination.Pagination, filter *RemarkListFilter, orderBys ...RemarkOrderBy) (remarks Remarks, total int, err error)
	Exists(ctx context.Context, filter RemarkExistsParams) (exists bool, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/remark_interactor.go -package=mock . RemarkInteractor
type RemarkInteractor interface {
	Create(ctx context.Context, remark *Remark) (err error)
	Update(ctx context.Context, remark *Remark) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetRemark(ctx context.Context, id uuid.UUID) (remark *Remark, err error)
	GetRemarks(ctx context.Context, pager *upagination.Pagination, filter *RemarkListFilter, orderBys ...RemarkOrderBy) (remarks Remarks, total int, err error)
	Exists(ctx context.Context, filter RemarkExistsParams) (exists bool, err error)
	RemarkSimpleUpdate(ctx context.Context, updateField RemarkSimpleUpdateType, remark *Remark) (err error)
}

type RemarkExistsParams struct {
	CategoryID uuid.UUID
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 更新时排除自身
}

type RemarkListFilter struct {
	CategoryID uuid.UUID  `json:"category_id"`
	MerchantID uuid.UUID  `json:"merchant_id"`
	StoreID    uuid.UUID  `json:"store_id"`
	Enabled    *bool      `json:"enabled"` // nil: 不筛; true/false: 按状态筛
	RemarkType RemarkType `json:"remark_type"`
}
