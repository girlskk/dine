package types

import "gitlab.jiguang.dev/pos-dine/dine/domain"

type PayHuifuCallback struct {
	Sign     string `form:"sign" binding:"required"`
	RespData string `form:"resp_data" binding:"required"`
}

// PayPollingReq 支付轮询请求
type PayPollingReq struct {
	SeqNo string `json:"seq_no" binding:"required"`
}

type PayPollingResp struct {
	State      domain.PayState `json:"state"`       // 支付状态, P 处理中, S 成功, F 失败
	FailReason string          `json:"fail_reason"` // 失败原因
}
