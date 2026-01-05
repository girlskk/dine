package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRequestDate(t *testing.T) {
	rqe := require.New(t)
	var rt RequestDate
	err := rt.UnmarshalJSON([]byte(`"2025-03-24"`))
	rqe.NoError(err)
	t.Log(time.Time(rt))

	type TestReq struct {
		Date RequestDate `json:"date"`
	}

	gin.SetMode(gin.TestMode)
	_, ok := binding.Validator.Engine().(*validator.Validate)
	rqe.True(ok)

	t.Run("正常解析", func(t *testing.T) {
		r := gin.New()
		var testReq TestReq
		r.POST("/", func(c *gin.Context) {
			err := c.ShouldBind(&testReq)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, testReq)
		})

		rqe := require.New(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"date": "2025-03-24"}`))
		rqe.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		rqe.Equal(http.StatusOK, w.Code)
		b := w.Body.Bytes()
		t.Log(string(b))
		err = json.Unmarshal(b, &testReq)
		rqe.NoError(err)
		rqe.Equal(time.Date(2025, 3, 24, 0, 0, 0, 0, time.Local), time.Time(testReq.Date))
	})

	t.Run("格式错误", func(t *testing.T) {
		r := gin.New()
		var testReq TestReq
		r.POST("/", func(c *gin.Context) {
			err := c.ShouldBind(&testReq)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, testReq)
		})

		rqe := require.New(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"date": "2025 03 24"}`))
		rqe.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		rqe.Equal(http.StatusBadRequest, w.Code)
		t.Log(w.Body.String())
	})

	t.Run("不传参数", func(t *testing.T) {
		r := gin.New()
		var testReq TestReq
		r.POST("/", func(c *gin.Context) {
			err := c.ShouldBind(&testReq)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, testReq)
		})

		rqe := require.New(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{}`))
		rqe.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		rqe.Equal(http.StatusOK, w.Code)
		b := w.Body.Bytes()
		t.Log(string(b))
		rqe.True(time.Time(testReq.Date).IsZero())
	})

	t.Run("传空字符串", func(t *testing.T) {
		r := gin.New()
		var testReq TestReq
		r.POST("/", func(c *gin.Context) {
			err := c.ShouldBind(&testReq)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, testReq)
		})

		rqe := require.New(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"date": ""}`))
		rqe.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		rqe.Equal(http.StatusOK, w.Code)
		b := w.Body.Bytes()
		t.Log(string(b))
		rqe.True(time.Time(testReq.Date).IsZero())
	})
}

func TestLastMonthStart(t *testing.T) {
	lastMonthStart := LastMonthStart(time.Now())
	t.Log(lastMonthStart)
}

func TestLastMonthEnd(t *testing.T) {
	lastMonthEnd := LastMonthEnd(time.Now())
	t.Log(lastMonthEnd)
}
