package e

// 自定义错误码，通常错误由错误码和错误信息两部分组成，便于跟踪和维护错误信息
//
// ----------------------------------
// 错误码规则：基于 HTTP 状态码分类设计
//
// 0     - 成功
// 40xxx - 400 Bad Request（请求参数错误、格式错误等）
// 41xxx - 401 Unauthorized（认证失败、token无效等）
// 43xxx - 403 Forbidden（权限不足等）
// 44xxx - 404 Not Found（资源不存在等）
// 50xxx - 500 Internal Server Error（系统内部错误）

const (
	// Success 表示成功
	Success = 0
)

// 40xxx - 400 Bad Request 相关错误
const (
	InvalidParams = iota + 40000 // 40000 - 请求参数错误
	BadRequest                   // 40001 - 业务请求错误
)

// 41xxx - 401 Unauthorized 相关错误
const (
	Unauthorized = iota + 41000 // 41000 - 未授权（认证失败、token无效等）
)

// 43xxx - 403 Forbidden 相关错误
const (
	Forbidden = iota + 43000 // 43000 - 禁止访问
)

// 44xxx - 404 Not Found 相关错误
const (
	NotFound = iota + 44000 // 44000 - 资源不存在

)

// 49xxx - 409 Conflict 相关错误
const (
	Conflict = iota + 49000 // 49000 - 资源冲突（如唯一约束冲突）
)

// 50xxx - 500 Internal Server Error 相关错误
const (
	InternalError   = iota + 50000 // 50000 - 系统内部错误
	DBError                        // 50001 - 数据库错误
	ThirdPartyError                // 50002 - 第三方服务错误
	UnknownError                   // 50003 - 未知错误
)
