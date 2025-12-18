package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/i18n"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
)

type ErrorHandling struct {
	alert     alert.Alert
	appConfig domain.AppConfig
}

func NewErrorHandling(alert alert.Alert, appConfig domain.AppConfig) *ErrorHandling {
	return &ErrorHandling{alert: alert, appConfig: appConfig}
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
			apiErr = errorx.New(http.StatusInternalServerError, errcode.UnknownError, err)
		}

		// 如果 Message 为空，则翻译错误码
		if apiErr.IsMessageEmpty() {
			ctx := c.Request.Context()
			// 使用 i18n 翻译错误码
			translated := i18n.Translate(ctx, apiErr.Code.String(), nil)

			fmt.Println("translated", translated)
			fmt.Println("apiErr.Code.String()", apiErr.Code.String())
			if translated != apiErr.Code.String() {
				// 翻译成功，更新 Message
				apiErr.Message = translated
			}
		}

		// 获取对应的 HTTP 状态码
		statusCode := apiErr.HTTPStatusCode()

		// 如果 debug 模式，则返回底层错误和堆栈信息
		isDebug := m.appConfig.RunMode == domain.RunModeDev
		apiErr.WithDebug(isDebug)

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
