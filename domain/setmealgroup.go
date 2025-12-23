package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 错误定义
// ------------------------------------------------------------

var (
	ErrSetMealGroupNoDetails               = errors.New("套餐组必须至少有一个详情项")
	ErrSetMealGroupNoDefaultDetail         = errors.New("每个套餐组中必须至少有一个默认项")
	ErrSetMealGroupDetailInvalid           = errors.New("套餐组详情商品无效")
	ErrSetMealGroupOptionalProductInvalid  = errors.New("备选商品无效")
	ErrSetMealGroupOptionalProductConflict = errors.New("备选商品不能是当前套餐组详情中的商品")
)

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// SetMealGroupSelectionType 套餐组点单限制
type SetMealGroupSelectionType string

const (
	SetMealGroupSelectionTypeFixed    SetMealGroupSelectionType = "fixed"    // 固定分组
	SetMealGroupSelectionTypeOptional SetMealGroupSelectionType = "optional" // 可选套餐
)

func (SetMealGroupSelectionType) Values() []string {
	return []string{
		string(SetMealGroupSelectionTypeFixed),
		string(SetMealGroupSelectionTypeOptional),
	}
}

// ------------------------------------------------------------
// 实体定义
// ------------------------------------------------------------

// SetMealGroup 套餐组
type SetMealGroup struct {
	ID            uuid.UUID                 `json:"id"`             // 套餐组ID
	Name          string                    `json:"name"`           // 套餐组名称
	ProductID     uuid.UUID                 `json:"product_id"`     // 套餐商品ID（外键）
	SelectionType SetMealGroupSelectionType `json:"selection_type"` // 点单限制：fixed（固定分组）、optional（可选套餐）
	CreatedAt     time.Time                 `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time                 `json:"updated_at"`     // 更新时间

	// 关联信息
	Details []*SetMealDetail `json:"details,omitempty"` // 套餐组详情列表
}

// SetMealGroups 套餐组集合
type SetMealGroups []*SetMealGroup

// SetMealDetail 套餐组详情
type SetMealDetail struct {
	ID                 uuid.UUID   `json:"id"`                             // 详情ID
	GroupID            uuid.UUID   `json:"group_id"`                       // 套餐组ID（外键）
	ProductID          uuid.UUID   `json:"product_id"`                     // 商品ID（外键，引用普通商品）
	Quantity           int         `json:"quantity"`                       // 数量（必选，必须为正整数）
	IsDefault          bool        `json:"is_default"`                     // 是否默认（必选，每个套餐组中只能有一个默认项）
	OptionalProductIDs []uuid.UUID `json:"optional_product_ids,omitempty"` // 备选商品ID列表（可选，多选）
	CreatedAt          time.Time   `json:"created_at"`                     // 创建时间
	UpdatedAt          time.Time   `json:"updated_at"`                     // 更新时间

	// 关联信息
	Product          *Product `json:"product,omitempty"`           // 关联商品
	OptionalProducts Products `json:"optional_products,omitempty"` // 备选商品
}

// SetMealDetails 套餐组详情集合
type SetMealDetails []*SetMealDetail

// ------------------------------------------------------------
// 仓储和用例接口
// ------------------------------------------------------------

// SetMealGroupRepository 套餐组仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/setmealgroup_repository.go -package=mock . SetMealGroupRepository
type SetMealGroupRepository interface {
	// 套餐组相关操作
	CreateGroups(ctx context.Context, groups []*SetMealGroup) error
	DeleteByProductID(ctx context.Context, productID uuid.UUID) error
	// 套餐组详情相关操作
	CreateDetails(ctx context.Context, details []*SetMealDetail) error
}
