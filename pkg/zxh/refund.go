package zxh

import "context"

type RefundParam struct {
	OutOrderID string `json:"outOrderId"`
	Amount     string `json:"amount"`
}

type RefundResult struct {
	Points    []PointWithAmount    `json:"points"`
	PaiPoints []PaiPointWithAmount `json:"paiPoints"`
}

func (c *Client) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	res := new(RefundResult)
	data := map[string]any{
		"outOrderId": param.OutOrderID,
		"amount":     param.Amount,
	}
	req := newRequest(c, RefundPath, data, res)
	if err := req.do(ctx); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) RefundQuery(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	res := new(RefundResult)
	data := map[string]any{
		"outOrderId": param.OutOrderID,
		"amount":     param.Amount,
	}
	req := newRequest(c, RefundQueryPath, data, res)
	if err := req.do(ctx); err != nil {
		return nil, err
	}
	return res, nil
}
