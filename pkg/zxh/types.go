package zxh

const (
	StatusPending Status = 1
	StatusSuccess Status = 2
	StatusFail    Status = 3
	StatusWaiting Status = 4
)

type Status int

type PointWithAmount struct {
	Code   string `json:"code"`
	Amount string `json:"amount"`
}

type PaiPointWithAmount struct {
	GoodsCode string `json:"goodsCode"`
	Amount    string `json:"amount"`
}

type ZhixinhuaPointPayCallBack struct {
	OutOrderID string               `json:"outOrderId"`
	MchId      string               `json:"mchId"`
	NonceStr   string               `json:"nonceStr"`
	Timestamp  string               `json:"timestamp"`
	Sign       string               `json:"sign"`
	Status     int                  `json:"status"`
	ErrMsg     string               `json:"errMsg,optional"`
	Points     []PointWithAmount    `json:"points,optional"`
	PaiPoints  []PaiPointWithAmount `json:"paiPoints,optional"`
}

func (c *ZhixinhuaPointPayCallBack) Verify(accessKey string) bool {
	points := make([]any, len(c.Points))
	for i, p := range c.Points {
		points[i] = p
	}

	paiPoints := make([]any, len(c.PaiPoints))
	for i, p := range c.PaiPoints {
		paiPoints[i] = p
	}

	data := map[string]any{
		"outOrderId": c.OutOrderID,
		"mchId":      c.MchId,
		"nonceStr":   c.NonceStr,
		"timestamp":  c.Timestamp,
		"status":     c.Status,
		"errMsg":     c.ErrMsg,
		"points":     points,
		"paiPoints":  paiPoints,
	}

	return c.Sign == sign(data, accessKey)
}
