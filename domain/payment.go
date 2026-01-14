package domain

// PaymentMethodType 支付方式类型
type PaymentMethodType string

const (
	PaymentMethodTypeCash PaymentMethodType = "CASH" // 现金
)

func (PaymentMethodType) Values() []string {
	return []string{
		string(PaymentMethodTypeCash),
	}
}
