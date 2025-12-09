package domain

import (
	"context"
	"errors"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// TableAreaRepository 区域仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/table_area_repository.go -package=mock . TableAreaRepository
type TableAreaRepository interface {
	FindByID(ctx context.Context, id int) (*TableArea, error)
	Exists(ctx context.Context, params AreaExistsParams) (bool, error)
	Create(ctx context.Context, area *TableArea) error
	Update(ctx context.Context, area *TableArea) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params AreaSearchParams) (*AreaSearchRes, error)
	IncreaseTableCount(ctx context.Context, id int) error
	DecreaseTableCount(ctx context.Context, id int) error
}

// TableAreaInteractor 区域用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/table_area_interactor.go -package=mock . TableAreaInteractor
type TableAreaInteractor interface {
	Create(ctx context.Context, area *TableArea) error
	Update(ctx context.Context, area *TableArea) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params AreaSearchParams) (*AreaSearchRes, error)
}

var (
	ErrAreaNameExists = errors.New("区域名称已存在")
	ErrAreaNotExists  = errors.New("区域不存在")
	ErrAreaHasTables  = errors.New("区域下存在台桌，无法删除")
)

// AreaSearchParams 区域查询参数
type AreaSearchParams struct {
	StoreID int
	Name    string
}

type AreaSearchRes struct {
	*upagination.Pagination
	Items TableAreas `json:"items"`
}

// AreaExistsParams 区域存在性检查参数
type AreaExistsParams struct {
	StoreID int
	Name    string
}

// TableArea 区域实体
type TableArea struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`        // 区域名称（如：一楼大厅、二楼包间）
	StoreID    int       `json:"store_id"`    // 所属门店ID
	TableCount int       `json:"table_count"` // 台桌个数
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}

// TableAreas 区域集合
type TableAreas []*TableArea
