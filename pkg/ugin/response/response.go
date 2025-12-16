package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
)

type Response struct {
	Code errcode.ErrCode `json:"code"`
	Data any             `json:"data,omitempty"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &Response{
		Code: errcode.Success,
		Data: data,
	})
}
