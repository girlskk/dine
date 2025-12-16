package errcode

// 业务错误码
type ErrCode string

func (ec ErrCode) String() string {
	return string(ec)
}

const (
	Success ErrCode = "SUCCESS" // Success 表示成功
	// 请求错误
	InvalidParams ErrCode = "INVALID_PARAMS" //   参数错误
	BadRequest    ErrCode = "BAD_REQUEST"    //   请求错误
	// 认证错误
	Unauthorized ErrCode = "UNAUTHORIZED" // 未授权（认证失败、token无效等）
	// 授权错误
	Forbidden ErrCode = "FORBIDDEN" // 禁止访问
	// 资源错误
	NotFound ErrCode = "NOT_FOUND" // 资源不存在
	Conflict ErrCode = "CONFLICT"  // 资源冲突

	// 系统错误
	InternalError ErrCode = "INTERNAL_ERROR" // 系统内部错误
	UnknownError  ErrCode = "UNKNOWN_ERROR"  // 未知错误

	UserNotFound ErrCode = "USER_NOT_FOUND" // 用户不存在
)
