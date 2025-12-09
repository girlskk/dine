package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	OrderFinanceLogType string
)

const (
	OrderFinanceLogTypePaid   OrderFinanceLogType = "paid"   // 支付
	OrderFinanceLogTypeRefund OrderFinanceLogType = "refund" // 退款
)

func (OrderFinanceLogType) Values() []string {
	return []string{
		string(OrderFinanceLogTypePaid),
		string(OrderFinanceLogTypeRefund),
	}
}

type OrderFinanceLog struct {
	ID          int                 `json:"id"`
	OrderID     int                 `json:"order_id"`
	Amount      decimal.Decimal     `json:"amount"`
	Type        OrderFinanceLogType `json:"type"`
	Channel     OrderPaidChannel    `json:"channel"`
	SeqNo       string              `json:"seq_no"`
	CreatorType OperatorType        `json:"creator_type"` // 创建人类型 frontend: 前台用户, backend: 后台用户, admin: 管理员, system: 系统
	CreatorID   int                 `json:"creator_id"`   // 创建人ID
	CreatorName string              `json:"creator_name"` // 创建人姓名
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
