package domain

import (
	"context"

	"github.com/google/uuid"
)

type FrontendContext struct {
	MerchantID uuid.UUID `json:"merchant_id"` // 品牌商ID
}

func (u *FrontendContext) GetMerchantID() uuid.UUID {
	return u.MerchantID
}

type (
	frontendContextKey struct{}
)

func NewFrontendContext(ctx context.Context, u *FrontendContext) context.Context {
	return context.WithValue(ctx, frontendContextKey{}, u)
}

func FromFrontendContext(ctx context.Context) *FrontendContext {
	if v, ok := ctx.Value(frontendContextKey{}).(*FrontendContext); ok {
		return v
	}
	return nil
}
