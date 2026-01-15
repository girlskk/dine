package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
)

type Auth struct{}

func (u *Auth) Name() string {
	return "Auth"
}

func NewAuth(handlers []ugin.Handler) *Auth {
	return &Auth{}
}

func (u *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		merchantIDStr := c.GetHeader("X-Merchant-ID")
		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			err = fmt.Errorf("failed to get header Merchant-ID: %w", err)
			c.Error(err)
			c.Abort()
			return
		}
		ctx = domain.NewFrontendContext(ctx, &domain.FrontendContext{
			MerchantID: merchantID,
		})
		c.Request = c.Request.Clone(ctx)
		c.Next()
	}
}

// 检查 full path 是否在允许 store_id 为空的白名单中
func isStoreIDOptional(fullPath string) bool {
	optionalPaths := map[string]struct{}{
		"/store/list": {}, // 获取门店列表接口路由
	}

	_, exists := optionalPaths[fullPath]
	return exists
}
