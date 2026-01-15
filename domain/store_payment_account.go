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
	ErrStorePaymentAccountNotExists                       = errors.New("门店收款账户不存在")
	ErrStorePaymentAccountPaymentAccountExist             = errors.New("该门店已存在该品牌商收款账户")
	ErrStorePaymentAccountMerchantNumberExist             = errors.New("支付商户号在当前门店+品牌商收款账户下已存在")
	ErrStorePaymentAccountStoreNotBelongMerchant          = errors.New("门店不属于当前品牌商")
	ErrStorePaymentAccountPaymentAccountNotBelongMerchant = errors.New("品牌商收款账户不属于当前品牌商")
)

// ------------------------------------------------------------
// 仓储接口
// ------------------------------------------------------------

// StorePaymentAccountRepository 门店收款账户仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_payment_account_repository.go -package=mock . StorePaymentAccountRepository
type StorePaymentAccountRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*StorePaymentAccount, error)
	Create(ctx context.Context, account *StorePaymentAccount) error
	Update(ctx context.Context, account *StorePaymentAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params StorePaymentAccountExistsParams) (bool, error)
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params StorePaymentAccountSearchParams) (*StorePaymentAccountSearchRes, error)
}

// ------------------------------------------------------------
// 用例接口
// ------------------------------------------------------------

// StorePaymentAccountInteractor 门店收款账户用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_payment_account_interactor.go -package=mock . StorePaymentAccountInteractor
type StorePaymentAccountInteractor interface {
	Create(ctx context.Context, account *StorePaymentAccount, user User) error
	Update(ctx context.Context, account *StorePaymentAccount, user User) error
	Delete(ctx context.Context, id uuid.UUID, user User) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params StorePaymentAccountSearchParams, user User) (*StorePaymentAccountSearchRes, error)
}

// ------------------------------------------------------------
// 实体定义
// ------------------------------------------------------------

// StorePaymentAccount 门店收款账户实体
type StorePaymentAccount struct {
	ID               uuid.UUID `json:"id"`                 // 门店收款账户ID
	MerchantID       uuid.UUID `json:"merchant_id"`        // 品牌商ID
	StoreID          uuid.UUID `json:"store_id"`           // 门店ID
	PaymentAccountID uuid.UUID `json:"payment_account_id"` // 品牌商收款账户ID
	MerchantNumber   string    `json:"merchant_number"`    // 支付商户号
	CreatedAt        time.Time `json:"created_at"`         // 创建时间
	UpdatedAt        time.Time `json:"updated_at"`         // 更新时间

	// 关联信息
	Store          *StoreSimple    `json:"store,omitempty"`           // 门店
	PaymentAccount *PaymentAccount `json:"payment_account,omitempty"` // 品牌商收款账户
}

// StorePaymentAccounts 门店收款账户集合
type StorePaymentAccounts []*StorePaymentAccount

// ------------------------------------------------------------
// 参数定义（DTO）
// ------------------------------------------------------------

// StorePaymentAccountExistsParams 存在性检查参数
type StorePaymentAccountExistsParams struct {
	StoreID          uuid.UUID // 门店ID
	PaymentAccountID uuid.UUID // 品牌商收款账户ID
	ExcludeID        uuid.UUID // 排除的ID（用于更新时检查唯一性）
}

// StorePaymentAccountSearchParams 查询参数
type StorePaymentAccountSearchParams struct {
	MerchantID     uuid.UUID   // 品牌商ID（必填）
	StoreIDs       []uuid.UUID // 门店ID列表（可选，多选）
	MerchantName   string      // 品牌商支付商户名称（可选，模糊匹配）
	CreatedAtStart *time.Time  // 创建时间开始（可选）
	CreatedAtEnd   *time.Time  // 创建时间结束（可选）
}

// StorePaymentAccountSearchRes 查询结果
type StorePaymentAccountSearchRes struct {
	*upagination.Pagination
	Items StorePaymentAccounts `json:"items"`
}
