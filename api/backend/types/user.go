package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type AccountCreateReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nickname" binding:"required"`
}

type AccountListReq struct {
	upagination.RequestPagination
}

// AccountListResp 账号列表响应
type AccountListResp struct {
	Items []*domain.FrontendUser `json:"items"` // 账号列表
	Total int                    `json:"total"` // 总数
}

type AccountUpdateReq struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"` // 可选，不填则不更新密码
	NickName string `json:"nickname" binding:"required"`
}
