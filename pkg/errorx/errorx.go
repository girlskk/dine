package errorx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
)

// 错误类型定义
type Error struct {
	Code       errcode.ErrCode `json:"code"`    // 业务错误代码
	HTTPStatus int             `json:"-"`       // HTTP 状态码
	Message    string          `json:"message"` // 错误消息
	// debug 模式显示的信息
	Err      error  `json:"err,omitempty"`      // 底层错误
	Func     string `json:"func,omitempty"`     // 函数名
	Position string `json:"position,omitempty"` // 错误发生的位置

	debug bool `json:"-"` // 是否为 debug 模式
}

// Error 实现 Error 接口
func (e *Error) Error() string {
	return e.Message
}

// Unwrap 实现 errors.Unwrap，用于错误链
func (e *Error) Unwrap() error {
	return e.Err
}

// HTTPStatusCode 返回对应的 HTTP 状态码
func (e *Error) HTTPStatusCode() int {
	return e.HTTPStatus
}

// New 创建新的错误
// httpStatus: HTTP 状态码
// code: 业务错误码（英文短语）
// err: 底层错误（可选，nil 时使用 code 作为默认消息）
func New(httpStatus int, code errcode.ErrCode, err error) *Error {
	message := string(code)

	e := &Error{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
		Err:        err,
	}

	// 如果是系统级别错误（5xx），记录堆栈信息
	if httpStatus >= http.StatusInternalServerError {
		e.caller(2)
	}

	return e
}

// WithMessage 链式设置自定义错误消息
func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

// WithDebug 设置 debug 模式，debug 模式下会在 JSON 中返回底层错误和堆栈信息
func (e *Error) WithDebug(debug bool) *Error {
	e.debug = debug
	return e
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

func (e *Error) IsMessageEmpty() bool {
	return e.Message == e.Code.String()
}

// MarshalJSON 自定义 JSON 序列化，根据 debug 模式决定是否返回底层错误
func (e *Error) MarshalJSON() ([]byte, error) {
	type Alias Error
	aux := &struct {
		*Alias
		Err      *string `json:"err,omitempty"`
		Func     *string `json:"func,omitempty"`
		Position *string `json:"position,omitempty"`
	}{
		Alias: (*Alias)(e),
	}

	// 只有在 debug 模式下才序列化底层错误和堆栈信息
	if e.debug {
		if e.Err != nil {
			errMsg := e.Err.Error()
			aux.Err = &errMsg
		}
		if e.Func != "" {
			aux.Func = &e.Func
		}
		if e.Position != "" {
			aux.Position = &e.Position
		}
	}

	return json.Marshal(aux)
}
