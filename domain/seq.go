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

const (
	AdminTaxSequenceKey          = "tax:admin:sequence"          // 运营后台税率序列号
	BackendTaxSequenceKey        = "tax:backend:sequence"        // 品牌后台税率序列号
	StoreTaxSequenceKey          = "tax:store:sequence"          // 门店后台税率序列号
	AdminDepartmentSequenceKey   = "department:admin:sequence"   // 运营后台部门序列号
	BackendDepartmentSequenceKey = "department:backend:sequence" // 品牌后台部门序列号
	StoreDepartmentSequenceKey   = "department:store:sequence"   // 门店后台部门序列号
	AdminUserSequenceKey         = "user:admin:sequence"         // 运营后台用户序列号
	BackendUserSequenceKey       = "user:backend:sequence"       // 品牌后台用户序列号
	StoreUserSequenceKey         = "user:store:sequence"         // 门店后台用户序列号
	AdminRoleSequenceKey         = "role:admin:sequence"         // 运营后台角色序列号
	BackendRoleSequenceKey       = "role:backend:sequence"       // 品牌后台角色序列号
	StoreRoleSequenceKey         = "role:store:sequence"         // 门店后台角色序列号
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/incr_sequence.go -package=mock . IncrSequence
type IncrSequence interface {
	Next(ctx context.Context) (string, error)
	Current(ctx context.Context) (string, error)
}
