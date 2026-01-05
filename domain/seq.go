package domain

import "context"

const (
	DailySequencePrefixOrderNo                  = "seq:order_no"
	DailySequencePrefixPayNo                    = "seq:payment_no"
	DailySequencePrefixStoreWithdrawNo          = "seq:store_withdraw_no"
	DailySequencePrefixProfitDistributionBillNo = "seq:profit_distribution_bill_no" // 分账账单编号前缀
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/daily_sequence.go -package=mock . DailySequence
type DailySequence interface {
	Next(ctx context.Context, prefix string) (int64, error)
	Current(ctx context.Context, prefix string) (int64, error)
}
