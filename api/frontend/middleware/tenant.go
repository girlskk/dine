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
			c.AbortWithStatus(http.StatusBadRequest)
			return
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
