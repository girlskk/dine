package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/utracing"
)

const HeaderXRequestID = "X-Request-Id"

type PopulateRequestID struct{}

func NewPopulateRequestID() *PopulateRequestID {
	return &PopulateRequestID{}
}

func (p *PopulateRequestID) Name() string {
	return "PopulateRequestID"
}

func (p *PopulateRequestID) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		ctx := r.Context()

		rid := r.Header.Get(HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}

		ctx = utracing.NewRequestIDContext(ctx, rid)
		c.Request = r.Clone(ctx)
		c.Header("X-Request-Id", rid)

		c.Next()
	}
}
