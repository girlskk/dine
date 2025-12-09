package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/utracing"
	"go.uber.org/zap"
)

func NewPopulateLogger(originalLogger *zap.SugaredLogger) *PopulateLogger {
	return &PopulateLogger{originalLogger: originalLogger}
}

type PopulateLogger struct {
	originalLogger *zap.SugaredLogger
}

func (p *PopulateLogger) Name() string {
	return "PopulateLogger"
}

func (p *PopulateLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		ctx := r.Context()

		logger := p.originalLogger

		if id := utracing.RequestIDFromContext(ctx); id != "" {
			logger = logger.With("request_id", id)
		}

		ctx = logging.NewContext(ctx, logger)
		c.Request = r.Clone(ctx)

		c.Next()
	}
}
