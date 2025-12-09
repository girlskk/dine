package domain

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_spec_rel_repository.go -package=mock . ProductSpecRelRepository
type ProductSpecRelRepository interface {
	ListByIDs(ctx context.Context, ids []int) (ProductSpecRels, error)
	Exists(ctx context.Context, params ProductSpecRelExistsParams) (bool, error)
	UpdateSaleStatusByIDs(ctx context.Context, ids []int, status ProductSaleStatus) error
	FindByID(ctx context.Context, id int) (*ProductSpecRel, error)
}

// 商品的商品规格
type ProductSpecRel struct {
	ID         int               `json:"id"`          // 商品-规格ID
	SpecID     int               `json:"spec_id"`     // 规格ID
	SpecName   string            `json:"spec_name"`   // 规格名称
	ProductID  int               `json:"product_id"`  // 商品ID
	Price      decimal.Decimal   `json:"price"`       // 规格价格
	SaleStatus ProductSaleStatus `json:"sale_status"` // 销售状态
}

type ProductSpecRels []*ProductSpecRel

type ProductSpecRelExistsParams struct {
	SpecID    int
	ProductID int
}
