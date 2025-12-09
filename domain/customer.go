// domain/customer.go
package domain

import (
	"context"
	"time"
)

type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderUnknown Gender = "unknown"
)

func (g Gender) Values() []string {
	return []string{
		string(GenderMale),
		string(GenderFemale),
		string(GenderUnknown),
	}
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/customer_repository.go -package=mock . CustomerRepository
type CustomerRepository interface {
	Find(ctx context.Context, id int) (*Customer, error)
	FindOrCreate(ctx context.Context, customer *Customer) (id int, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/customer_interactor.go -package=mock . CustomerInteractor
type CustomerInteractor interface {
	WXLogin(ctx context.Context, code string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *Customer, err error)
}

// Customer 顾客实体
type Customer struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"`   // 昵称
	Phone     string    `json:"phone"`      // 手机号，必填
	Avatar    string    `json:"avatar"`     // 头像
	Gender    Gender    `json:"gender"`     // 性别
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

func (u *Customer) GetOperatorID() int {
	return u.ID
}

// 顾客手机号作为操作者
func (u *Customer) GetOperatorName() string {
	return u.Phone
}

func (u *Customer) GetOperatorType() OperatorType {
	return OperatorTypeCustomer
}

func (u *Customer) GetOperatorStoreID() int {
	return 0
}

type (
	customerKey struct{}
)

func NewCustomerContext(ctx context.Context, u *Customer) context.Context {
	return context.WithValue(ctx, customerKey{}, u)
}

func FromCustomerContext(ctx context.Context) *Customer {
	if v, ok := ctx.Value(customerKey{}).(*Customer); ok {
		return v
	}
	return nil
}
