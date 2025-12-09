package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type SystemHandler struct {
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Routes(r gin.IRouter) {
	r = r.Group("/system")
	r.POST("/now", h.Now())
}

func (h *SystemHandler) NoAuths() []string {
	return []string{
		"/system/now",
	}
}

// Now
//
//	@Tags		系统管理
//	@Security	BearerAuth
//	@Summary	获取当前时间
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	types.SystemNowResp	"成功"
//	@Router		/system/now [post]
func (h *SystemHandler) Now() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "SystemHandler.Now")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("SystemHandler.Now")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		response.Ok(c, types.SystemNowResp{
			Now: time.Now(),
		})
	}
}
