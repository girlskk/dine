package zxh

import "context"

type PointLimit struct {
	Code  string `json:"code"`
	Limit string `json:"limit"`
}

type PaiPointLimit struct {
	GoodsCode string `json:"goodsCode"`
	Limit     string `json:"limit"`
}

type PaymentParam struct {
	PayCode      string          `json:"payCode"`
	MerchantName string          `json:"merchantName"`
	OutOrderID   string          `json:"outOrderId"`
	Desc         string          `json:"desc"`
	NotifyURL    string          `json:"notifyUrl"`
	Amount       string          `json:"amount"`
	Points       []PointLimit    `json:"points"`
	PaiPoints    []PaiPointLimit `json:"paiPoints"`
}

func (c *Client) Payment(ctx context.Context, param *PaymentParam) error {
	var res any

	points := make([]any, len(param.Points))
	for i, p := range param.Points {
		points[i] = p
	}

	paiPoints := make([]any, len(param.PaiPoints))
	for i, p := range param.PaiPoints {
		paiPoints[i] = p
	}

	data := map[string]any{
		"payCode":      param.PayCode,
		"merchantName": param.MerchantName,
		"outOrderId":   param.OutOrderID,
		"desc":         param.Desc,
		"notifyUrl":    param.NotifyURL,
		"amount":       param.Amount,
		"points":       points,
		"paiPoints":    paiPoints,
	}
	req := newRequest(c, PaymentPath, data, res)
	return req.do(ctx)
}
