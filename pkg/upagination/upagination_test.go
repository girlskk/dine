package upagination

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRequestPagination(t *testing.T) {
	rqe := require.New(t)
	var rp RequestPagination
	err := json.Unmarshal([]byte(`{"page": 1, "size": 10}`), &rp)
	rqe.NoError(err)

	t.Log(rp)

	p := rp.ToPagination()
	t.Log(p)

	gin.SetMode(gin.TestMode)
	_, ok := binding.Validator.Engine().(*validator.Validate)
	rqe.True(ok)

	type TestReq struct {
		RequestPagination
		Name string `json:"name"`
	}

	r := gin.New()
	r.POST("/", func(c *gin.Context) {
		var req TestReq
		err := c.ShouldBind(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	t.Run("正常解析", func(t *testing.T) {
		rqe := require.New(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"page": 1, "size": 10, "name": "test"}`))
		rqe.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		rqe.Equal(http.StatusOK, w.Code)

		var resp TestReq
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		rqe.NoError(err)

		rqe.Equal(1, resp.Page)
		rqe.Equal(10, resp.Size)
		rqe.Equal("test", resp.Name)
	})

	t.Run("分页参数校验", func(t *testing.T) {
		rqe := require.New(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"page": -1, "size": 10, "name": "test"}`))
		rqe.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		rqe.Equal(http.StatusBadRequest, w.Code)

		t.Log(w.Body.String())
	})

	t.Run("不传分页参数", func(t *testing.T) {
		rqe := require.New(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"name": "test"}`))
		rqe.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		rqe.Equal(http.StatusOK, w.Code)

		var resp TestReq
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		rqe.NoError(err)

		rqe.Equal(0, resp.Page)
		rqe.Equal(0, resp.Size)
		rqe.Equal("test", resp.Name)
	})
}
