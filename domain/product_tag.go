package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrProductTagNotExists         = errors.New("商品标签不存在")
	ErrProductTagNameExists        = errors.New("商品标签名称已存在")
	ErrProductTagDeleteHasProducts = errors.New("商品标签下有商品，不能删除")
)

// ProductTagRepository 商品标签仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_tag_repository.go -package=mock . ProductTagRepository
type ProductTagRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*ProductTag, error)
	Create(ctx context.Context, tag *ProductTag) error
	Update(ctx context.Context, tag *ProductTag) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params ProductTagExistsParams) (bool, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductTagSearchParams) (*ProductTagSearchRes, error)
}

// ProductTagInteractor 商品标签用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_tag_interactor.go -package=mock . ProductTagInteractor
type ProductTagInteractor interface {
	Create(ctx context.Context, tag *ProductTag) error
	Update(ctx context.Context, tag *ProductTag) error
	Delete(ctx context.Context, id uuid.UUID) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params ProductTagSearchParams) (*ProductTagSearchRes, error)
}

// ProductTag 商品标签实体
type ProductTag struct {
	ID           uuid.UUID `json:"id"`            // 标签ID
	Name         string    `json:"name"`          // 标签名称
	MerchantID   uuid.UUID `json:"merchant_id"`   // 品牌商ID
	StoreID      uuid.UUID `json:"store_id"`      // 门店ID
	ProductCount int       `json:"product_count"` // 关联的商品数量
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
}

// ProductTags 商品标签集合
type ProductTags []*ProductTag

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// ProductTagExistsParams 存在性检查参数
type ProductTagExistsParams struct {
	MerchantID uuid.UUID
	Name       string
	ExcludeID  uuid.UUID // 排除的ID（用于更新时检查名称唯一性）
}

// ProductTagSearchParams 查询参数
type ProductTagSearchParams struct {
	MerchantID uuid.UUID
	Name       string
}

type ProductTagSearchRes struct {
	*upagination.Pagination
	Items ProductTags `json:"items"`
}
