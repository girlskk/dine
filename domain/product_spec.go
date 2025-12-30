package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrProductSpecNotExists         = errors.New("商品规格不存在")
	ErrProductSpecNameExists        = errors.New("商品规格名称已存在")
	ErrProductSpecDeleteHasProducts = errors.New("商品规格下有商品，不能删除")
)

// ProductSpecRepository 商品规格仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_spec_repository.go -package=mock . ProductSpecRepository
type ProductSpecRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*ProductSpec, error)
	Create(ctx context.Context, spec *ProductSpec) error
	CreateBulk(ctx context.Context, specs ProductSpecs) error
	Update(ctx context.Context, spec *ProductSpec) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params ProductSpecExistsParams) (bool, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductSpecSearchParams) (*ProductSpecSearchRes, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) (ProductSpecs, error)
	FindByNamesInStore(ctx context.Context, storeID uuid.UUID, names []string) (ProductSpecs, error)
}

// ProductSpecInteractor 商品规格用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_spec_interactor.go -package=mock . ProductSpecInteractor
type ProductSpecInteractor interface {
	Create(ctx context.Context, spec *ProductSpec) error
	Update(ctx context.Context, spec *ProductSpec, user User) error
	Delete(ctx context.Context, id uuid.UUID, user User) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductSpecSearchParams) (*ProductSpecSearchRes, error)
}

// ProductSpec 商品规格实体
type ProductSpec struct {
	ID           uuid.UUID `json:"id"`            // 规格ID
	Name         string    `json:"name"`          // 规格名称
	MerchantID   uuid.UUID `json:"merchant_id"`   // 品牌商ID
	StoreID      uuid.UUID `json:"store_id"`      // 门店ID
	ProductCount int       `json:"product_count"` // 关联的商品数量
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
}

// ProductSpecs 商品规格集合
type ProductSpecs []*ProductSpec

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// ProductSpecExistsParams 存在性检查参数
type ProductSpecExistsParams struct {
	MerchantID uuid.UUID
	StoreID    uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// ProductSpecSearchParams 查询参数
type ProductSpecSearchParams struct {
	MerchantID   uuid.UUID
	StoreID      uuid.UUID
	Name         string
	OnlyMerchant bool
}

type ProductSpecSearchRes struct {
	*upagination.Pagination
	Items ProductSpecs `json:"items"`
}
