package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/e"
)

// 错误类型定义
type Error struct {
	Code       int    `json:"code"`    // 业务错误代码
	HTTPStatus int    `json:"-"`       // HTTP 状态码（不序列化到 JSON）
	Message    string `json:"message"` // 错误消息
	Func       string `json:"-"`       // 函数名（仅系统错误）
	Position   string `json:"-"`       // 错误发生的位置（仅系统错误）
}

// Error 实现 Error 接口
func (e *Error) Error() string {
	return e.Message
}

// HTTPStatusCode 返回对应的 HTTP 状态码
func (e *Error) HTTPStatusCode() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}
	// 根据业务错误码自动映射 HTTP 状态码
	return mapCodeToHTTPStatus(e.Code)
}

// Fail 错误处理函数
func Fail(code int, err error) *Error {
	// 判断 err 是否为自定义错误类型 *Error
	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr
	}

	// 如果只传入错误码，返回默认错误对象；否则返回自定义错误对象
	var message string
	if err == nil {
		message = e.GetMsg(code)
	} else {
		message = err.Error()
	}

	customErr = newError(code, message)

	// 如果是系统级别错误，需要返回调用堆栈信息
	if code >= e.InternalError {
		return customErr.caller(2)
	}
	return customErr
}

// FailWithStatus 创建错误并指定 HTTP 状态码（用于特殊场景覆盖默认映射）
func FailWithStatus(code int, httpStatus int, err error) *Error {
	var message string
	if err == nil {
		message = e.GetMsg(code)
	} else {
		message = err.Error()
	}

	customErr := &Error{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
	}

	// 如果是系统级别错误，需要返回调用堆栈信息
	if code >= e.InternalError {
		return customErr.caller(2)
	}
	return customErr
}

// Failf 使用格式化字符串创建错误
func Failf(code int, format string, args ...any) *Error {
	message := fmt.Errorf(format, args...).Error()
	customErr := newError(code, message)

	// 如果是系统级别错误，需要返回调用堆栈信息
	if code >= e.InternalError {
		return customErr.caller(2)
	}
	return customErr
}

// FailWithStatusf 使用格式化字符串创建错误并指定 HTTP 状态码
func FailWithStatusf(code int, httpStatus int, format string, args ...any) *Error {
	message := fmt.Errorf(format, args...).Error()
	customErr := &Error{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
	}

	// 如果是系统级别错误，需要返回调用堆栈信息
	if code >= e.InternalError {
		return customErr.caller(2)
	}
	return customErr
}

// IsCode 检查错误是否为指定的错误码
func IsCode(err error, code int) bool {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.Code == code
}

// 获取调用堆栈信息
func (e *Error) caller(skip ...int) *Error {
	skip = append(skip, 1)
	pc, file, line, ok := runtime.Caller(skip[0])
	if ok {
		e.Func = runtime.FuncForPC(pc).Name()
		e.Position = fmt.Sprintf("%s:%d", file, line)
	}
	return e
}

// 创建自定义错误对象
func newError(code int, message string) *Error {
	return &Error{
		Code:       code,
		HTTPStatus: mapCodeToHTTPStatus(code),
		Message:    message,
	}
}

// mapCodeToHTTPStatus 根据业务错误码映射到 HTTP 状态码
func mapCodeToHTTPStatus(code int) int {
	switch {
	case code == 0:
		return http.StatusOK
	case code >= 40000 && code < 41000:
		// 40xxx -> 400 Bad Request
		return http.StatusBadRequest
	case code >= 41000 && code < 43000:
		// 41xxx -> 401 Unauthorized
		return http.StatusUnauthorized
	case code >= 43000 && code < 44000:
		// 43xxx -> 403 Forbidden
		return http.StatusForbidden
	case code >= 44000 && code < 50000:
		// 44xxx -> 404 Not Found
		return http.StatusNotFound
	case code >= 49000 && code < 50000:
		// 49xxx -> 409 Conflict
		return http.StatusConflict
	case code >= 50000 && code < 60000:
		// 50xxx -> 500 Internal Server Error
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
