package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type AccountListReq struct {
	Enabled *bool `form:"enabled"`
}

type AccountListResp struct {
	Users []*domain.StoreUser `json:"users"`
	Total int                 `json:"total"`
}
