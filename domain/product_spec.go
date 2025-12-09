package domain

import (
	"context"
	"errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"time"
)

// ProductSpecRepository 规格仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_spec_repository.go -package=mock . ProductSpecRepository
type ProductSpecRepository interface {
	FindByID(ctx context.Context, id int) (*ProductSpec, error)
	Exists(ctx context.Context, params SpecExistsParams) (bool, error)
	Create(ctx context.Context, spec *ProductSpec) error
	Update(ctx context.Context, spec *ProductSpec) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params SpecSearchParams) (*SpecSearchRes, error)
	ListByIDs(ctx context.Context, ids []int) (ProductSpecs, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_spec_interactor.go -package=mock . ProductSpecInteractor
type ProductSpecInteractor interface {
	Create(ctx context.Context, spec *ProductSpec) error
	Update(ctx context.Context, spec *ProductSpec) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params SpecSearchParams) (*SpecSearchRes, error)
}

var (
	ErrSpecNameExists = errors.New("规格名称已存在")
	ErrSpecNotExists  = errors.New("规格不存在")
	ErrSpecUsing      = errors.New("商品规格正在使用，无法删除")
)

// ProductSpec 商品规格实体
type ProductSpec struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // 规格名称（如：大杯、中杯）
	StoreID   int       `json:"store_id"`   // 所属门店ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// ProductSpecs 规格集合
type ProductSpecs []*ProductSpec

// SpecExistsParams 存在性检查参数
type SpecExistsParams struct {
	StoreID int
	Name    string
}

type SpecSearchParams struct {
	StoreID int
	Name    string
}

type SpecSearchRes struct {
	*upagination.Pagination
	Items ProductSpecs `json:"items"`
}
