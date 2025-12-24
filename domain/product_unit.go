package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrProductUnitNotExists         = errors.New("商品单位不存在")
	ErrProductUnitNameExists        = errors.New("商品单位名称已存在")
	ErrProductUnitDeleteHasProducts = errors.New("商品单位下有商品，不能删除")
)

// ProductUnitRepository 商品单位仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_unit_repository.go -package=mock . ProductUnitRepository
type ProductUnitRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*ProductUnit, error)
	Create(ctx context.Context, unit *ProductUnit) error
	Update(ctx context.Context, unit *ProductUnit) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params ProductUnitExistsParams) (bool, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductUnitSearchParams) (*ProductUnitSearchRes, error)
}

// ProductUnitInteractor 商品单位用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_unit_interactor.go -package=mock . ProductUnitInteractor
type ProductUnitInteractor interface {
	Create(ctx context.Context, unit *ProductUnit) error
	Update(ctx context.Context, unit *ProductUnit) error
	Delete(ctx context.Context, id uuid.UUID) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductUnitSearchParams) (*ProductUnitSearchRes, error)
}

// 商品单位类型
type ProductUnitType string

const (
	ProductUnitTypeQuantity ProductUnitType = "quantity" // 数量单位
	ProductUnitTypeWeight   ProductUnitType = "weight"   // 重量单位
)

func (ProductUnitType) Values() []string {
	return []string{
		string(ProductUnitTypeQuantity),
		string(ProductUnitTypeWeight),
	}
}

// ProductUnit 商品单位实体
type ProductUnit struct {
	ID           uuid.UUID       `json:"id"`            // 单位ID
	Name         string          `json:"name"`          // 单位名称
	Type         ProductUnitType `json:"type"`          // 单位类型：quantity（数量单位）、weight（重量单位）
	MerchantID   uuid.UUID       `json:"merchant_id"`   // 品牌商ID
	StoreID      uuid.UUID       `json:"store_id"`      // 门店ID
	ProductCount int             `json:"product_count"` // 关联的商品数量
	CreatedAt    time.Time       `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time       `json:"updated_at"`    // 更新时间
}

// ProductUnits 商品单位集合
type ProductUnits []*ProductUnit

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// ProductUnitExistsParams 存在性检查参数
type ProductUnitExistsParams struct {
	MerchantID uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// ProductUnitSearchParams 查询参数
type ProductUnitSearchParams struct {
	MerchantID uuid.UUID
	Name       string
	Type       ProductUnitType
}

type ProductUnitSearchRes struct {
	*upagination.Pagination
	Items ProductUnits `json:"items"`
}
