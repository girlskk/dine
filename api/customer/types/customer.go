package types

type WXLoginReq struct {
	Code string `json:"code"`
}

type WXLoginResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}
