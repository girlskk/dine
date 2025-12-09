package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// 门店合作类型
type StoreCooperationType string

const (
	StoreCooperationTypeJoin StoreCooperationType = "join"
)

func (StoreCooperationType) Values() []string {
	return []string{
		string(StoreCooperationTypeJoin),
	}
}

// 门店类型
type StoreType string

const (
	StoreTypeRestaurant StoreType = "restaurant" // 中餐点餐（需要台桌）
	StoreTypeCafeteria  StoreType = "cafeteria"  // 大食堂（不需要台桌）
)

func (StoreType) Values() []string {
	return []string{
		string(StoreTypeRestaurant),
		string(StoreTypeCafeteria),
	}
}

type Store struct {
	ID                  int                  `json:"id"`
	Name                string               `json:"name"`                  // 名称
	Type                StoreType            `json:"type"`                  // 类型
	CooperationType     StoreCooperationType `json:"cooperation_type"`      // 合作类型
	NeedAudit           bool                 `json:"need_audit"`            // 是否需要审核
	Enabled             bool                 `json:"enabled"`               // 是否启用
	PointSettlementRate decimal.Decimal      `json:"point_settlement_rate"` // 积分结算费率（单位：百分比，例如 0.1234 表示 12.34%）
	PointWithdrawalRate decimal.Decimal      `json:"point_withdrawal_rate"` // 积分提现费率
	HuifuID             string               `json:"huifu_id"`              // 汇付ID
	ZxhID               string               `json:"zxh_id"`                // 知心话ID
	ZxhSecret           string               `json:"zxh_secret"`            // 知心话密钥
	CreatedAt           time.Time            `json:"created_at"`            // 创建时间
	UpdatedAt           time.Time            `json:"updated_at"`            // 更新时间

	Info        *StoreInfo    `json:"info,omitempty"`         // 门店资料
	Finance     *StoreFinance `json:"finance,omitempty"`      // 财务信息
	BackendUser *BackendUser  `json:"backend_user,omitempty"` // 后台用户
}

type Stores []*Store // 门店列表

type StoreInfoImages struct {
	Logo        string `json:"logo"`         // 店标图片
	Front       string `json:"front"`        // 门店正面图
	Dine        string `json:"dine"`         // 就餐环境图
	License     string `json:"license"`      // 营业执照图
	FoodLicense string `json:"food_license"` // 食品经营许可证
}

// 门店资料
type StoreInfo struct {
	City         string          `json:"city"`          // 省市地区
	Address      string          `json:"address"`       // 详细地址
	ContactName  string          `json:"contact_name"`  // 门店联系人
	ContactPhone string          `json:"contact_phone"` // 联系人电话
	Images       StoreInfoImages `json:"images"`        // 门店图片
}

// 门店财务信息
type StoreFinance struct {
	StoreID          int    `json:"store_id"`
	BankAccount      string `json:"bank_account"`       // 银行账号
	BankCardName     string `json:"bank_card_name"`     // 银行账户名称
	BankName         string `json:"bank_name"`          // 银行名称
	BranchName       string `json:"branch_name"`        // 开户支行
	PublicAccount    string `json:"public_account"`     // 对公账号
	CompanyName      string `json:"company_name"`       // 公司名称
	PublicBankName   string `json:"public_bank_name"`   // 对公银行名称
	PublicBranchName string `json:"public_branch_name"` // 对公开户支行
	CreditCode       string `json:"credit_code"`        // 统一社会信用代码
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_interactor.go -package=mock . StoreInteractor
type StoreInteractor interface {
	Create(ctx context.Context, store *Store, user *BackendUser) error
	Update(ctx context.Context, store *Store, user *BackendUser) error
	UpdateByStore(ctx context.Context, store *Store, user *BackendUser) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params StoreSearchParams) (*StoreSearchRes, error)
	GetDetail(ctx context.Context, id int) (*Store, error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_repository.go -package=mock . StoreRepository
type StoreRepository interface {
	Find(ctx context.Context, id int) (*Store, error)
	Create(ctx context.Context, store *Store) error
	Exists(ctx context.Context, params StoreExistsParams) (bool, error)
	Update(ctx context.Context, store *Store) error
	PagedListBySearch(ctx context.Context, page *upagination.Pagination, params StoreSearchParams) (*StoreSearchRes, error)
	GetDetail(ctx context.Context, id int) (*Store, error)
	ListAll(ctx context.Context) (Stores, error)
}

var (
	ErrStoreNameExists = errors.New("门店名称已存在")
	ErrStoreNotExists  = errors.New("门店不存在")
)

type (
	storeKey struct{}
)

func NewStoreContext(ctx context.Context, s *Store) context.Context {
	return context.WithValue(ctx, storeKey{}, s)
}

func FromStoreContext(ctx context.Context) *Store {
	if v, ok := ctx.Value(storeKey{}).(*Store); ok {
		return v
	}
	return nil
}

type StoreExistsParams struct {
	Name string
}

type StoreSearchParams struct {
	City string
	Name string
}

type StoreSearchRes struct {
	*upagination.Pagination
	Items Stores `json:"items"`
}
