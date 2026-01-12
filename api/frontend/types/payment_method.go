package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// PaymentMethodListReq 结算方式列表请求
type PaymentMethodListReq struct {
	upagination.RequestPagination
	StoreID string `json:"store_id,omitempty" form:"store_id" ` // 门店ID（可选）
}
