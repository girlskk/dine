package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}
