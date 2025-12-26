package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/product_attr_rel_repository.go -package=mock . ProductAttrRelRepository
type ProductAttrRelRepository interface {
	CreateBulk(ctx context.Context, relations ProductAttrRelations) error
}

// ProductAttrRelation 商品口味做法关联实体
type ProductAttrRelation struct {
	ID         uuid.UUID `json:"id"`           // 关联ID
	ProductID  uuid.UUID `json:"product_id"`   // 商品ID
	AttrID     uuid.UUID `json:"attr_id"`      // 口味做法ID
	AttrItemID uuid.UUID `json:"attr_item_id"` // 口味做法项ID
	IsDefault  bool      `json:"is_default"`   // 是否默认项
	CreatedAt  time.Time `json:"created_at"`   // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`   // 更新时间

	// 关联信息
	Attr     *ProductAttr     `json:"attr"`      // 口味做法
	AttrItem *ProductAttrItem `json:"attr_item"` // 口味做法项
}

type ProductAttrRelations []*ProductAttrRelation
