package types

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// ProductAttrItemReq 口味做法项请求（用于创建和更新）
type ProductAttrItemReq struct {
	ID        uuid.UUID       `json:"id,omitempty"`                    // 口味做法项ID（更新时传入）
	Name      string          `json:"name" binding:"required,max=255"` // 口味做法项名称
	Image     string          `json:"image,omitempty"`                 // 图片URL（可选）
	BasePrice decimal.Decimal `json:"base_price" binding:"required"`   // 基础加价（单位：分）
}

// ProductAttrCreateReq 创建商品口味做法请求
type ProductAttrCreateReq struct {
	Name     string               `json:"name" binding:"required,max=255"`                                                     // 口味做法名称
	Channels []domain.SaleChannel `json:"channels" binding:"required,min=1,dive,oneof=POS Mobile Scan SelfService ThirdParty"` // 售卖渠道列表（必选，可多选）
	Items    []ProductAttrItemReq `json:"items,omitempty"`                                                                     // 口味做法项列表（可选）
}

// ProductAttrUpdateReq 更新商品口味做法请求
type ProductAttrUpdateReq struct {
	Name     string               `json:"name" binding:"required,max=255"`                                                     // 口味做法名称
	Channels []domain.SaleChannel `json:"channels" binding:"required,min=1,dive,oneof=POS Mobile Scan SelfService ThirdParty"` // 售卖渠道列表（必选，可多选）
	Items    []ProductAttrItemReq `json:"items,omitempty"`                                                                     // 口味做法项列表（可选，用于新增、修改、删除）
}
