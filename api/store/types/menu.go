package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// MenuListReq 菜单列表请求
type MenuListReq struct {
	upagination.RequestPagination
	Name string `json:"name" form:"name"` // 菜单名称（模糊匹配）
}
