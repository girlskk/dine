package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// ------------------------------------------------------------
// 错误定义
// ------------------------------------------------------------

var (
	ErrPaymentAccountNotExists           = errors.New("收款账户不存在")
	ErrPaymentAccountMerchantNumberExist = errors.New("支付商户号在当前品牌商+渠道下已存在")
	ErrPaymentAccountHasStoreAccounts    = errors.New("收款账户下有绑定门店收款账户，不可删除")
)

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// PaymentChannel 支付渠道
type PaymentChannel string

const (
	PaymentChannelRM PaymentChannel = "rm" // Revenue Monster
)

func (PaymentChannel) Values() []string {
	return []string{
		string(PaymentChannelRM),
	}
}

// ------------------------------------------------------------
// 仓储接口
// ------------------------------------------------------------

// PaymentAccountRepository 收款账户仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/payment_account_repository.go -package=mock . PaymentAccountRepository
type PaymentAccountRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*PaymentAccount, error)
	Create(ctx context.Context, account *PaymentAccount) error
	Update(ctx context.Context, account *PaymentAccount) error
	UpdateAllDefaultStatus(ctx context.Context, merchantID uuid.UUID, isDefault bool) error
	FindForUpdateByMerchantID(ctx context.Context, merchantID uuid.UUID) (PaymentAccounts, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params PaymentAccountExistsParams) (bool, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params PaymentAccountSearchParams) (*PaymentAccountSearchRes, error)
	CountStoreAccounts(ctx context.Context, id uuid.UUID) (int, error)
}

// ------------------------------------------------------------
// 用例接口
// ------------------------------------------------------------

// PaymentAccountInteractor 收款账户用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/payment_account_interactor.go -package=mock . PaymentAccountInteractor
type PaymentAccountInteractor interface {
	Create(ctx context.Context, account *PaymentAccount) error
	Update(ctx context.Context, account *PaymentAccount, user User) error
	UpdateDefaultStatus(ctx context.Context, id uuid.UUID, user User) error
	Delete(ctx context.Context, id uuid.UUID, user User) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params PaymentAccountSearchParams) (*PaymentAccountSearchRes, error)
}

// ------------------------------------------------------------
// 实体定义
// ------------------------------------------------------------

// PaymentAccount 收款账户实体
type PaymentAccount struct {
	ID             uuid.UUID      `json:"id"`              // 收款账户ID
	MerchantID     uuid.UUID      `json:"merchant_id"`     // 品牌商ID
	Channel        PaymentChannel `json:"channel"`         // 支付渠道
	MerchantNumber string         `json:"merchant_number"` // 支付商户号
	MerchantName   string         `json:"merchant_name"`   // 支付商户名称
	IsDefault      bool           `json:"is_default"`      // 是否默认
	CreatedAt      time.Time      `json:"created_at"`      // 创建时间
	UpdatedAt      time.Time      `json:"updated_at"`      // 更新时间
}

// PaymentAccounts 收款账户集合
type PaymentAccounts []*PaymentAccount

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// PaymentAccountExistsParams 存在性检查参数
type PaymentAccountExistsParams struct {
	MerchantID uuid.UUID      // 品牌商ID
	Channel    PaymentChannel // 支付渠道
	ExcludeID  uuid.UUID      // 排除的ID（用于更新时检查唯一性）
}

// PaymentAccountSearchParams 查询参数
type PaymentAccountSearchParams struct {
	MerchantID     uuid.UUID      // 品牌商ID（必填）
	Channel        PaymentChannel // 支付渠道（可选）
	MerchantName   string         // 支付商户名称（可选，模糊匹配）
	CreatedAtStart *time.Time     // 创建时间开始（可选）
	CreatedAtEnd   *time.Time     // 创建时间结束（可选）
}

// PaymentAccountSearchRes 查询结果
type PaymentAccountSearchRes struct {
	*upagination.Pagination
	Items PaymentAccounts `json:"items"`
}
