package domain

import (
	"context"
	"errors"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// CategoryRepository 仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/category_repository.go -package=mock . CategoryRepository
type CategoryRepository interface {
	FindByID(ctx context.Context, id int) (*Category, error)
	Exists(ctx context.Context, params CategoryExistsParams) (bool, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params CategorySearchParams) (*CategorySearchRes, error)
}

// CategoryInteractor 用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/category_interactor.go -package=mock . CategoryInteractor
type CategoryInteractor interface {
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params CategorySearchParams) (*CategorySearchRes, error)
}

var (
	ErrCategoryNameExists = errors.New("分类名称已存在")
	ErrCategoryNotExists  = errors.New("分类不存在")
	ErrCategoryUsing      = errors.New("分类正在使用中，无法删除")
)

// Category 商品分类实体
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // 分类名称
	StoreID   int       `json:"store_id"`   // 所属门店ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// Categories 分类集合
type Categories []*Category

// CategorySearchParams 查询参数
type CategorySearchParams struct {
	StoreID int
	Name    string
	Status  int
}

type CategorySearchRes struct {
	*upagination.Pagination
	Items Categories `json:"items"`
}

// CategoryExistsParams 存在性检查参数
type CategoryExistsParams struct {
	StoreID int
	Name    string
}
