package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/uhttp"
)

type JSONResponseWriterLogger struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func NewJSONResponseWriterLogger(writer gin.ResponseWriter) *JSONResponseWriterLogger {
	return &JSONResponseWriterLogger{
		ResponseWriter: writer,
		buf:            new(bytes.Buffer),
	}
}

func (w *JSONResponseWriterLogger) isJSONResponse() bool {
	mt, _, err := mime.ParseMediaType(w.Header().Get("Content-Type"))
	if err == nil && (mt == "application/json" || mt == "application/javascript") {
		return true
	}
	return false
}

func (w *JSONResponseWriterLogger) Write(b []byte) (int, error) {
	if w.isJSONResponse() {
		w.buf.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *JSONResponseWriterLogger) WriteString(s string) (int, error) {
	if w.isJSONResponse() {
		w.buf.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w *JSONResponseWriterLogger) Content() map[string]any {
	if !w.isJSONResponse() {
		return nil
	}

	var resp map[string]any
	if err := json.Unmarshal(w.buf.Bytes(), &resp); err != nil {
		return nil
	}

	return resp
}

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Name() string {
	return "Logger"
}

func (l *Logger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		ctx := r.Context()
		logger := logging.FromContext(ctx).Named("middleware.Logger")
		start := time.Now()

		mediaType, _, _ := mime.ParseMediaType(c.GetHeader("Content-Type"))
		params := make(map[string]any)
		if mediaType != uhttp.CONTENT_TYPE_FILE {
			var reqBody []byte
			if r.Body != nil { // Read
				var err error
				if reqBody, err = uhttp.BytesFromRequest(r); err != nil {
					logger.Errorw("failed to uhttp.BytesFromRequest", "error", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			var err error
			if mediaType == uhttp.CONTENT_TYPE_FORM {
				_params := make(map[string]string)
				if err = c.ShouldBind(&_params); err == nil {
					for k, v := range _params {
						params[k] = v
					}
				}
			} else {
				c.ShouldBind(&params)
			}

			if reqBody != nil {
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset
			}
		}

		// resp body log
		rwl := NewJSONResponseWriterLogger(c.Writer)
		c.Writer = rwl

		c.Next()

		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		logger = logger.With(
			"ip", c.ClientIP(),
			"method", r.Method,
			"url", r.URL.String(),
			"ua", c.GetHeader("User-Agent"),
			"rs", c.Writer.Status(),
			"tc", timeConsuming,
		)

		if len(params) > 0 {
			logger = logger.With("req", params)
		}

		if resp := rwl.Content(); resp != nil {
			logger = logger.With("resp", resp)
		}

		logger.Info("access")
	}
}
