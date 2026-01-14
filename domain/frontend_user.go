package domain

import (
	"context"

	"github.com/google/uuid"
)

// FrontendUser 前台用户（收银机等客户端）
type FrontendUser struct {
	MerchantID uuid.UUID `json:"merchant_id"` // 品牌商ID
	StoreID    uuid.UUID `json:"store_id"`    // 门店ID
}

// GetMerchantID 实现 User 接口
func (u *FrontendUser) GetMerchantID() uuid.UUID {
	return u.MerchantID
}

// GetStoreID 实现 User 接口
func (u *FrontendUser) GetStoreID() uuid.UUID {
	return u.StoreID
}

// GetUserID 实现 User 接口
func (u *FrontendUser) GetUserID() uuid.UUID {
	return uuid.Nil
}

// GetUserType 实现 User 接口
func (u *FrontendUser) GetUserType() UserType {
	return UserTypeFrontend
}

type (
	frontendUserKey struct{}
)

func NewFrontendUserContext(ctx context.Context, u *FrontendUser) context.Context {
	return context.WithValue(ctx, frontendUserKey{}, u)
}

func FromFrontendUserContext(ctx context.Context) *FrontendUser {
	if v, ok := ctx.Value(frontendUserKey{}).(*FrontendUser); ok {
		return v
	}
	return nil
}
