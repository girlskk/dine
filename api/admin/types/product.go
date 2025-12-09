package types

import "gitlab.jiguang.dev/pos-dine/dine/domain"

type ProductListReq struct {
	Page    int                  `json:"page"`
	Size    int                  `json:"size"`
	Name    string               `json:"name"`                                 // 商品名称或编号
	StoreID int                  `json:"store_id"`                             // 门店ID
	Status  domain.ProductStatus `json:"status" binding:"omitempty,oneof=1 2"` // 商品状态：1-待审核 2-审核通过
}

type ProductDetailReq struct {
	ID int `json:"id" binding:"required"`
}

type ProductApproveReq struct {
	IDs           []int `json:"ids" binding:"required"` // 商品ID列表
	AllowPointPay *bool `json:"allow_point_pay"`        // 是否允许积分支付
}

type ProductUnApproveReq struct {
	IDs []int `json:"ids" binding:"required"` // 商品ID列表
}
