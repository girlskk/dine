package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTimeLimiter(timeout int, skippers ...SkipperFunc) *TimeLimiter {
	return &TimeLimiter{timeout: timeout, skippers: skippers}
}

type TimeLimiter struct {
	timeout  int
	skippers []SkipperFunc
}

func (t *TimeLimiter) Name() string {
	return "TimeLimiter"
}

func (t *TimeLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, t.skippers...) {
			c.Next()
			return
		}

		r := c.Request
		ctx, cancel := context.WithTimeout(
			r.Context(),
			time.Second*time.Duration(t.timeout),
		)
		defer cancel()
		c.Request = r.Clone(ctx)

		c.Next()
	}
}
