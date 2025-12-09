package huifu

import "gitlab.jiguang.dev/pos-dine/dine/pkg/huifu"

func New(cfg huifu.MerchSysConfig) *huifu.BsPay {
	return &huifu.BsPay{IsProdMode: true, Msc: &cfg}
}
