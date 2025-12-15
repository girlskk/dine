package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/e"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
)

type ErrorHandling struct {
	alert alert.Alert
}

func NewErrorHandling(alert alert.Alert) *ErrorHandling {
	return &ErrorHandling{alert: alert}
}

func (m *ErrorHandling) Name() string {
	return "ErrorHandling"
}

func (m *ErrorHandling) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		ginErr := c.Errors.Last()
		if ginErr == nil {
			return
		}

		err := ginErr.Err
		// 判断是否为 errorx.Error 类型
		var apiErr *errorx.Error
		ok := errors.As(err, &apiErr)
		if !ok {
			// 如果不是 errorx.Error，转换为 UnknownError
			apiErr = errorx.Fail(e.UnknownError, err)
		}

		// 获取对应的 HTTP 状态码
		statusCode := apiErr.HTTPStatusCode()

		// 如果是系统级别错误（5xx），记录日志并发送告警
		if statusCode >= http.StatusInternalServerError {
			ctx := c.Request.Context()
			logger := logging.FromContext(ctx).Named("middleware.ErrorHandling")
			logger.Errorw("http handle internal error", "error", err, "code", apiErr.Code)

			go m.alert.Notify(ctx, err.Error())
		}
		// 使用对应的 HTTP 状态码返回错误
		c.AbortWithStatusJSON(statusCode, apiErr)
	}
}
