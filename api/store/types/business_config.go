package types

import "gitlab.jiguang.dev/pos-dine/dine/domain"

// BusinessConfigListReq 经营设置列表列表请求
type BusinessConfigListReq struct {
	Name  string                     `json:"name" form:"name"`   // 设置名称（模糊匹配）
	Group domain.BusinessConfigGroup `json:"group" form:"group"` // 配置分组
}
