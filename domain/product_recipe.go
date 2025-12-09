package domain

import (
	"context"
	"errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"time"
)

// ProductRecipeRepository 做法仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_recipe_repository.go -package=mock . ProductRecipeRepository
type ProductRecipeRepository interface {
	FindByID(ctx context.Context, id int) (*ProductRecipe, error)
	Exists(ctx context.Context, params RecipeExistsParams) (bool, error)
	Create(ctx context.Context, recipe *ProductRecipe) error
	Update(ctx context.Context, recipe *ProductRecipe) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params RecipeSearchParams) (*RecipeSearchRes, error)
	ListByIDs(ctx context.Context, ids []int) (ProductRecipes, error)
	IsUsedByProduct(ctx context.Context, id int) (bool, error)
}

// ProductRecipeInteractor 用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_recipe_interactor.go -package=mock . ProductRecipeInteractor
type ProductRecipeInteractor interface {
	Create(ctx context.Context, recipe *ProductRecipe) error
	Update(ctx context.Context, recipe *ProductRecipe) error
	Delete(ctx context.Context, id int) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params RecipeSearchParams) (*RecipeSearchRes, error)
}

var (
	ErrRecipeNameExists = errors.New("做法名称已存在")
	ErrRecipeNotExists  = errors.New("做法不存在")
	ErrRecipeUsing      = errors.New("商品做法正在使用，无法删除")
)

// ProductRecipes 商品做法实体
type ProductRecipe struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // 做法名称（如：加冰、去糖）
	StoreID   int       `json:"store_id"`   // 所属门店ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// ProductRecipes 做法集合
type ProductRecipes []*ProductRecipe

// RecipeExistsParams 存在性检查参数
type RecipeExistsParams struct {
	StoreID int
	Name    string
}

type RecipeSearchParams struct {
	StoreID int
	Name    string
}

type RecipeSearchRes struct {
	*upagination.Pagination
	Items ProductRecipes `json:"items"`
}
