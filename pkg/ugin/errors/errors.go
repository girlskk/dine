package errors

import (
	"fmt"
	"net/http"
)

const (
	CodeUnknownError = 500
)

type BizError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *BizError) Error() string {
	return e.Msg
}

func NewBizError(code int, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

func NotFound(msg string) *BizError {
	return NewBizError(http.StatusNotFound, msg)
}

func NotFoundf(format string, args ...any) *BizError {
	return NewBizError(http.StatusNotFound, fmt.Sprintf(format, args...))
}

func BadRequest(msg string) *BizError {
	return NewBizError(http.StatusBadRequest, msg)
}

func BadRequestf(format string, args ...any) *BizError {
	return NewBizError(http.StatusBadRequest, fmt.Sprintf(format, args...))
}

func Forbidden(msg string) *BizError {
	return NewBizError(http.StatusForbidden, msg)
}

func Forbiddenf(format string, args ...any) *BizError {
	return NewBizError(http.StatusForbidden, fmt.Sprintf(format, args...))
}

func Unauthorized(msg string) *BizError {
	return NewBizError(http.StatusUnauthorized, msg)
}

func Unauthorizedf(format string, args ...any) *BizError {
	return NewBizError(http.StatusUnauthorized, fmt.Sprintf(format, args...))
}
