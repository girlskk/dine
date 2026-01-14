package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// Tenant 租户中间件，从 header 中获取 merchant_id 和 store_id
type Tenant struct{}

func (t *Tenant) Name() string {
	return "Tenant"
}

func NewTenant() *Tenant {
	return &Tenant{}
}

func (t *Tenant) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantIDStr := c.GetHeader("X-Merchant-ID")
		storeIDStr := c.GetHeader("X-Store-ID")

		if merchantIDStr == "" || storeIDStr == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		storeID, err := uuid.Parse(storeIDStr)
		if err != nil {
			// 解析store ID失败时，检查是否允许 store_id 为空，不影响正常请求
			if !isStoreIDOptional(c.FullPath()) {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		user := &domain.FrontendUser{
			MerchantID: merchantID,
			StoreID:    storeID,
		}

		ctx := c.Request.Context()
		ctx = domain.NewFrontendUserContext(ctx, user)
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
