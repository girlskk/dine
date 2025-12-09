package domain

import "context"

const (
	DailySequencePrefixOrderNo         = "seq:order_no"
	DailySequencePrefixPayNo           = "seq:payment_no"
	DailySequencePrefixStoreWithdrawNo = "seq:store_withdraw_no"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/daily_sequence.go -package=mock . DailySequence
type DailySequence interface {
	Next(ctx context.Context, prefix string) (int64, error)
	Current(ctx context.Context, prefix string) (int64, error)
}
