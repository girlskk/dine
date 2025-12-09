package domain

import (
	"context"
	"errors"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrAttrNameExists = errors.New("商品属性已存在")
	ErrAttrNotExists  = errors.New("商品属性不存在")
	ErrAttrUsing      = errors.New("商品属性正在使用，无法删除")
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_attr_repository.go -package=mock . ProductAttrRepository
type ProductAttrRepository interface {
	FindByID(ctx context.Context, id int) (*ProductAttr, error)
	Exists(ctx context.Context, params AttrExistsParams) (bool, error)
	Create(ctx context.Context, attr *ProductAttr) error
	Update(ctx context.Context, attr *ProductAttr) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params AttrSearchParams) (*AttrSearchRes, error)
	ListByIDs(ctx context.Context, ids []int) (ProductAttrs, error)
	IsUsedByProduct(ctx context.Context, id int) (bool, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_attr_interactor.go -package=mock . ProductAttrInteractor
type ProductAttrInteractor interface {
	Create(ctx context.Context, attr *ProductAttr) error
	Update(ctx context.Context, attr *ProductAttr) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params AttrSearchParams) (*AttrSearchRes, error)
}

type AttrSearchParams struct {
	StoreID int
	Name    string
}

type AttrSearchRes struct {
	*upagination.Pagination
	Items []*ProductAttr `json:"items"`
}

type AttrExistsParams struct {
	StoreID int
	Name    string
}

type ProductAttr struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // 属性名称（如：辣度、甜度）
	StoreID   int       `json:"store_id"`   // 所属门店ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

type ProductAttrs []*ProductAttr
