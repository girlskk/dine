package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrStoreNameExists               = errors.New("门店名称已存在")
	ErrStoreNotExists                = errors.New("门店不存在")
	ErrStoreBusinessHoursConflict    = errors.New("门店营业时间段时间重复")
	ErrStoreBusinessHoursTimeInvalid = errors.New("门店营业时间段开始时间需早于结束时间")
	ErrStoreDiningPeriodConflict     = errors.New("门店就餐时段时间重复")
	ErrStoreDiningPeriodTimeInvalid  = errors.New("门店就餐时段开始时间需早于结束时间")
	ErrStoreDiningPeriodNameExists   = errors.New("门店就餐时段名称已存在")
	ErrStoreShiftTimeConflict        = errors.New("门店班次时间重复")
	ErrStoreShiftTimeTimeInvalid     = errors.New("门店班次开始时间需早于结束时间")
	ErrStoreShiftTimeNameExists      = errors.New("门店班次名称已存在")
)

// StoreRepository 门店仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_repository.go -package=mock . StoreRepository
type StoreRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (domainStore *Store, err error)
	FindStoreMerchant(ctx context.Context, merchantID uuid.UUID) (domainStore *Store, err error)
	Create(ctx context.Context, domainStore *Store) (err error)
	Update(ctx context.Context, domainStore *Store) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetStores(ctx context.Context, pager *upagination.Pagination, filter *StoreListFilter, orderBys ...StoreListOrderBy) (domainStores []*Store, total int, err error)
	ExistsStore(ctx context.Context, existsStoreParams *ExistsStoreParams) (exists bool, err error)
	CountStoresByMerchantID(ctx context.Context, merchantIDs []uuid.UUID) (storeCounts []*MerchantStoreCount, err error)
}

// StoreInteractor 门店用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_interactor.go -package=mock . StoreInteractor
type StoreInteractor interface {
	CreateStore(ctx context.Context, domainCStoreParams *CreateStoreParams) (err error)
	UpdateStore(ctx context.Context, domainUStoreParams *UpdateStoreParams) (err error)
	DeleteStore(ctx context.Context, id uuid.UUID) (err error)
	GetStore(ctx context.Context, id uuid.UUID) (domainStore *Store, err error)
	GetStores(ctx context.Context, pager *upagination.Pagination, filter *StoreListFilter, orderBys ...StoreListOrderBy) (domainStores []*Store, total int, err error)
	GetStoreByMerchantID(ctx context.Context, merchantID uuid.UUID) (domainStore *Store, err error)
	StoreSimpleUpdate(ctx context.Context, updateField StoreSimpleUpdateType, domainUStoreParams *UpdateStoreParams) (err error)
}

type StoreListOrderByType int

const (
	_ StoreListOrderByType = iota
	StoreListOrderByID
	StoreListOrderByCreatedAt
)

type StoreListOrderBy struct {
	OrderBy StoreListOrderByType
	Desc    bool
}

func NewStoreListOrderByID(desc bool) StoreListOrderBy {
	return StoreListOrderBy{
		OrderBy: StoreListOrderByID,
		Desc:    desc,
	}
}

func NewStoreListOrderByCreatedAt(desc bool) StoreListOrderBy {
	return StoreListOrderBy{
		OrderBy: StoreListOrderByCreatedAt,
		Desc:    desc,
	}
}

type StoreStatus string

const (
	StoreStatusOpen   StoreStatus = "open"   // 营业
	StoreStatusClosed StoreStatus = "closed" // 停业
)

func (StoreStatus) Values() []string {
	return []string{
		string(StoreStatusOpen),
		string(StoreStatusClosed),
	}
}

func (s StoreStatus) ToString() string {
	switch s {
	case StoreStatusOpen:
		return "营业"
	case StoreStatusClosed:
		return "停业"
	default:
		return ""
	}
}

type BusinessModel string

const (
	BusinessModelDirect     BusinessModel = "direct"     // 直营
	BusinessModelFranchisee BusinessModel = "franchisee" // 加盟
)

func (BusinessModel) Values() []string {
	return []string{
		string(BusinessModelDirect),
		string(BusinessModelFranchisee),
	}
}

func (b BusinessModel) ToString() string {
	switch b {
	case BusinessModelDirect:
		return "直营"
	case BusinessModelFranchisee:
		return "加盟"
	default:
		return ""
	}
}

type StoreSimpleUpdateType string

const (
	StoreSimpleUpdateTypeStatus StoreSimpleUpdateType = "status" // 状态更新
)

type Store struct {
	ID                      uuid.UUID       `json:"id"`
	MerchantID              uuid.UUID       `json:"merchant_id"`                // 商户 ID
	MerchantName            string          `json:"merchant_name"`              // 商户名称
	AdminPhoneNumber        string          `json:"admin_phone_number"`         // 管理员手机号
	StoreName               string          `json:"store_name"`                 // 门店名称,长度不超过30个字
	StoreShortName          string          `json:"store_short_name"`           // 门店简称
	StoreCode               string          `json:"store_code"`                 // 门店编码(保留字段)
	Status                  StoreStatus     `json:"status"`                     // 状态: 营业 停业
	BusinessModel           BusinessModel   `json:"business_model"`             // 经营模式：直营 加盟
	BusinessTypeID          uuid.UUID       `json:"business_type_id"`           // 业态类型
	BusinessTypeName        string          `json:"business_type_name"`         // 业务类型名称
	LocationNumber          string          `json:"location_number"`            // 门店位置编号
	ContactName             string          `json:"contact_name"`               // 联系人
	ContactPhone            string          `json:"contact_phone"`              // 联系电话
	UnifiedSocialCreditCode string          `json:"unified_social_credit_code"` // 统一社会信用代码
	StoreLogo               string          `json:"store_logo"`                 // logo 图片地址
	BusinessLicenseURL      string          `json:"business_license_url"`       // 营业执照图片地址
	StorefrontURL           string          `json:"storefront_url"`             // 门店门头照片地址
	CashierDeskURL          string          `json:"cashier_desk_url"`           // 收银台照片地址
	DiningEnvironmentURL    string          `json:"dining_environment_url"`     // 就餐环境照片地址
	FoodOperationLicenseURL string          `json:"food_operation_license_url"` // 食品经营许可证照片地址
	LoginAccount            string          `json:"login_account"`              // 登录账号
	LoginPassword           string          `json:"login_password"`             // 登录密码(加密存储)
	BusinessHours           []BusinessHours `json:"business_hours"`             // 营业时间段
	DiningPeriods           []DiningPeriod  `json:"dining_periods"`             // 就餐时段
	ShiftTimes              []ShiftTime     `json:"shift_times"`                // 班次时间
	Address                 *Address        `json:"address"`                    // 地址
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

type StoreSimple struct {
	ID        uuid.UUID `json:"id"`
	StoreName string    `json:"store_name"` // 门店名称
}

type StoreListFilter struct {
	StoreName        string        `json:"store_name"`         // 门店名称
	MerchantID       uuid.UUID     `json:"merchant_id"`        // 商户 ID
	BusinessTypeID   uuid.UUID     `json:"business_type_id"`   // 业态类型
	AdminPhoneNumber string        `json:"admin_phone_number"` // 管理员手机号
	Status           StoreStatus   `json:"status"`             // 状态: 营业 停业
	BusinessModel    BusinessModel `json:"business_model"`     // 经营模式：直营 加盟
	CreatedAtGte     *time.Time    `json:"created_at_gte"`     // 创建时间 大于等于
	CreatedAtLte     *time.Time    `json:"created_at_lte"`     // 创建时间 小于等于
	ProvinceID       uuid.UUID     `json:"province_id"`        // 省份 ID
}

type CreateStoreParams struct {
	MerchantID              uuid.UUID       `json:"merchant_id"`                // 商户 ID
	AdminPhoneNumber        string          `json:"admin_phone_number"`         // 管理员手机号
	StoreName               string          `json:"store_name"`                 // 门店名称,长度不超过30个字
	StoreShortName          string          `json:"store_short_name"`           // 门店简称
	StoreCode               string          `json:"store_code"`                 // 门店编码(保留字段)
	Status                  StoreStatus     `json:"status"`                     // 状态: 营业 停业
	BusinessModel           BusinessModel   `json:"business_model"`             // 经营模式：直营 加盟
	BusinessTypeID          uuid.UUID       `json:"business_type_id"`           // 业态类型
	LocationNumber          string          `json:"location_number"`            // 门店位置编号
	ContactName             string          `json:"contact_name"`               // 联系人
	ContactPhone            string          `json:"contact_phone"`              // 联系电话
	UnifiedSocialCreditCode string          `json:"unified_social_credit_code"` // 统一社会信用代码
	StoreLogo               string          `json:"store_logo"`                 // logo 图片地址
	BusinessLicenseURL      string          `json:"business_license_url"`       // 营业执照图片地址
	StorefrontURL           string          `json:"storefront_url"`             // 门店门头照片地址
	CashierDeskURL          string          `json:"cashier_desk_url"`           // 收银台照片地址
	DiningEnvironmentURL    string          `json:"dining_environment_url"`     // 就餐环境照片地址
	FoodOperationLicenseURL string          `json:"food_operation_license_url"` // 食品经营许可证照片地址
	LoginAccount            string          `json:"login_account"`              // 登录账号
	LoginPassword           string          `json:"login_password"`             // 登录密码(加密存储)
	BusinessHours           []BusinessHours `json:"business_hours"`             // 营业时间段
	DiningPeriods           []DiningPeriod  `json:"dining_periods"`             // 就餐时段
	ShiftTimes              []ShiftTime     `json:"shift_times"`                // 班次时间
	Address                 *Address        `json:"address"`                    // 地址
}

type UpdateStoreParams struct {
	ID                      uuid.UUID       `json:"id"`
	AdminPhoneNumber        string          `json:"admin_phone_number"`         // 管理员手机号
	StoreName               string          `json:"store_name"`                 // 门店名称,长度不超过30个字
	StoreShortName          string          `json:"store_short_name"`           // 门店简称
	StoreCode               string          `json:"store_code"`                 // 门店编码(保留字段)
	Status                  StoreStatus     `json:"status"`                     // 状态: 营业 停业
	BusinessModel           BusinessModel   `json:"business_model"`             // 经营模式：直营 加盟
	BusinessTypeID          uuid.UUID       `json:"business_type_id"`           // 业态类型
	LocationNumber          string          `json:"location_number"`            // 门店位置编号
	ContactName             string          `json:"contact_name"`               // 联系人
	ContactPhone            string          `json:"contact_phone"`              // 联系电话
	UnifiedSocialCreditCode string          `json:"unified_social_credit_code"` // 统一社会信用代码
	StoreLogo               string          `json:"store_logo"`                 // logo 图片地址
	BusinessLicenseURL      string          `json:"business_license_url"`       // 营业执照图片地址
	StorefrontURL           string          `json:"storefront_url"`             // 门店门头照片地址
	CashierDeskURL          string          `json:"cashier_desk_url"`           // 收银台照片地址
	DiningEnvironmentURL    string          `json:"dining_environment_url"`     // 就餐环境照片地址
	FoodOperationLicenseURL string          `json:"food_operation_license_url"` // 食品经营许可证照片地址
	LoginPassword           string          `json:"login_password"`             // 登录密码(加密存储)
	BusinessHours           []BusinessHours `json:"business_hours"`             // 营业时间段
	DiningPeriods           []DiningPeriod  `json:"dining_periods"`             // 就餐时段
	ShiftTimes              []ShiftTime     `json:"shift_times"`                // 班次时间
	Address                 *Address        `json:"address"`                    // 地址
}
type MerchantStoreCount struct {
	MerchantID uuid.UUID `json:"merchant_id"` // 商户 ID
	StoreCount int       `json:"store_count"` // 门店数量
}
type ExistsStoreParams struct {
	MerchantID uuid.UUID `json:"merchant_id"` // 商户 ID
	StoreName  string    `json:"store_name"`  // 门店名称
	ExcludeID  uuid.UUID `json:"exclude_id"`  // 排除的门店 ID
}

// BusinessHours 表示营业时间段
type BusinessHours struct {
	Weekdays  []time.Weekday `json:"weekdays"`   // 适用的星期几，0=星期日，1=星期一，依此类推
	StartTime string         `json:"start_time"` // 开始时间，格式 HH:MM:SS
	EndTime   string         `json:"end_time"`   // 结束时间，格式 HH:MM:SS
}

// DiningPeriod 表示用餐时段
type DiningPeriod struct {
	Name      string `json:"name"`       // 名称，如“用餐时段1”
	StartTime string `json:"start_time"` // 开始时间，格式 HH:MM:SS
	EndTime   string `json:"end_time"`   // 结束时间，格式 HH:MM:SS
}

// ShiftTime 表示班次时间
type ShiftTime struct {
	Name      string `json:"name"`       // 名称，如“班次1”
	StartTime string `json:"start_time"` // 开始时间，格式 HH:MM:SS
	EndTime   string `json:"end_time"`   // 结束时间，格式 HH:MM:SS
}
