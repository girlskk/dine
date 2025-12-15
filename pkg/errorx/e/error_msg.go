package e

var MsgFlags = map[int]string{
	Success:         "ok",
	InvalidParams:   "请求参数错误",
	BadRequest:      "请求错误",
	Unauthorized:    "未授权",
	Forbidden:       "禁止访问",
	NotFound:        "资源不存在",
	Conflict:        "资源冲突",
	InternalError:   "系统内部错误",
	DBError:         "数据库错误",
	ThirdPartyError: "第三方服务错误",
	UnknownError:    "未知错误",
}

// 从错误码中获取错误提示
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[UnknownError]
}
