package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRemarkCategoryNotExists  = errors.New("备注分类不存在")
	ErrRemarkCategoryNameExists = errors.New("备注分类名称已存在")
)

// RemarkCategories 聚合类型
type RemarkCategories []*RemarkCategory

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/remark_category_repository.go -package=mock . RemarkCategoryRepository
type RemarkCategoryRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (remarkCategory *RemarkCategory, err error)
	Create(ctx context.Context, remarkCategory *RemarkCategory) (err error)
	Update(ctx context.Context, remarkCategory *RemarkCategory) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetRemarkCategories(ctx context.Context, filter *RemarkCategoryListFilter) (remarkCategories RemarkCategories, err error)
	Exists(ctx context.Context, params RemarkCategoryExistsParams) (exists bool, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/remark_category_interactor.go -package=mock . RemarkCategoryInteractor
type RemarkCategoryInteractor interface {
	Create(ctx context.Context, remarkCategory *RemarkCategory) (err error)
	Update(ctx context.Context, remarkCategory *RemarkCategory) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetRemarkCategories(ctx context.Context, filter *RemarkCategoryListFilter) (remarkCategories RemarkCategories, err error)
	GetRemarkGroup(ctx context.Context, params RemarkGroupListFilter) (remarkGroups []RemarkGroup, err error)
}

type RemarkCategory struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`         // 分类名称
	RemarkScene RemarkScene `json:"remark_scene"` // 使用场景：整单备注/单品备注/退菜原因等
	MerchantID  uuid.UUID   `json:"merchant_id"`  // 品牌商ID，可为空表示系统级分类
	Description string      `json:"description"`  // 分类描述
	SortOrder   int         `json:"sort_order"`   // 排序，值越小越靠前
	RemarkCount int         `json:"remark_count"` // 该分类下备注数量
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type RemarkCategoryExistsParams struct {
	MerchantID uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 更新时排除自身
}

type RemarkCategoryListFilter struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	CountScene RemarkType // 统计备注数量场景
}
