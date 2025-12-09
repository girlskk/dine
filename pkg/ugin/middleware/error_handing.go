package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
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

		var apiErr *uerr.BizError
		ok := errors.As(err, &apiErr)
		if !ok {
			apiErr = uerr.NewBizError(uerr.CodeUnknownError, "服务异常")
		}

		if apiErr.Code == uerr.CodeUnknownError {
			ctx := c.Request.Context()
			logger := logging.FromContext(ctx).Named("middleware.ErrorHandling")
			logger.Errorw("http handle internal error", "error", err)

			go m.alert.Notify(ctx, err.Error())

			c.AbortWithStatusJSON(http.StatusInternalServerError, apiErr)
			return
		}

		c.JSON(http.StatusOK, apiErr)
	}
}
