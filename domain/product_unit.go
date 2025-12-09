package domain

import (
	"context"
	"errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"time"
)

var (
	ErrUnitNameExists = errors.New("商品单位已存在")
	ErrUnitNotExists  = errors.New("商品单位不存在")
	ErrUnitUsing      = errors.New("商品单位正在使用中，无法删除")
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_unit_repository.go -package=mock . ProductUnitRepository
type ProductUnitRepository interface {
	FindByID(ctx context.Context, id int) (*ProductUnit, error)
	Exists(ctx context.Context, params UnitExistsParams) (bool, error)
	Create(ctx context.Context, unit *ProductUnit) error
	Update(ctx context.Context, unit *ProductUnit) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params UnitSearchParams) (*UnitSearchRes, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_unit_interactor.go -package=mock . ProductUnitInteractor
type ProductUnitInteractor interface {
	Create(ctx context.Context, unit *ProductUnit) error
	Update(ctx context.Context, unit *ProductUnit) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params UnitSearchParams) (*UnitSearchRes, error)
}

type UnitSearchParams struct {
	StoreID int
	Name    string
}

type UnitSearchRes struct {
	*upagination.Pagination
	Items []*ProductUnit `json:"items"`
}

type UnitExistsParams struct {
	StoreID int
	Name    string
}

type ProductUnit struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // 单位名称（如：份、杯、碗）
	StoreID   int       `json:"store_id"`   // 所属门店ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

type ProductUnits []*ProductUnit
