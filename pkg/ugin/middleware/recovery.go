package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/metrics"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/uruntime"
)

type Recovery struct {
	alert alert.Alert
}

func NewRecovery(alert alert.Alert) *Recovery {
	return &Recovery{alert: alert}
}

func (r *Recovery) Name() string {
	return "Recovery"
}

func (r *Recovery) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if p := recover(); p != nil {
				ctx := c.Request.Context()
				logger := logging.FromContext(ctx).Named("middleware.Recovery")

				stack := uruntime.Stack(3)
				logger.Errorw(
					"http handle panic",
					"panic", p,
					"stack", string(stack),
				)
				metrics.RecoverCounter.Inc()

				go r.alert.Notify(ctx, fmt.Sprintf("panic: %v\n%s", p, stack))

				c.Status(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
