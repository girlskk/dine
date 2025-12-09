package zxh

import (
	"testing"
)

const (
	secretKey = "8d68ae9285df4f759dba62aaa8611b95" // 签名的密钥
)

func TestSign(t *testing.T) {
	params := map[string]any{
		"mchId":     "100001",
		"nonceStr":  "1845030987138863104",
		"timestamp": "1728724641",
		"status":    3,
		"errMsg":    "拍拍积分不足",
		"points": []PointWithAmount{
			{
				Code:   "TotalScore",
				Amount: "10",
			},
		},
		"paiPoints":  []PaiPointWithAmount{},
		"outOrderId": "10011325902917202410121717209110",
	}

	s := sign(params, secretKey)
	t.Log(s)
}
