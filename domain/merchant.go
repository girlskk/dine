package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrMerchantNameExists = errors.New("商户名称已存在")
	ErrMerchantNotExists  = errors.New("商户不存在")
)

// MerchantRepository 商户仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/merchant_repository.go -package=mock . MerchantRepository
type MerchantRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (domainMerchant *Merchant, err error)
	Create(ctx context.Context, domainMerchant *Merchant) (err error)
	Update(ctx context.Context, domainMerchant *Merchant) (err error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetMerchants(ctx context.Context, pager *upagination.Pagination, filter *MerchantListFilter, orderBys ...MerchantListOrderBy) (domainMerchants []*Merchant, total int, err error)
	CountMerchant(ctx context.Context) (merchantCount *MerchantCount, err error)
	ExistMerchant(ctx context.Context, merchantExistsParams *MerchantExistsParams) (exist bool, err error)
}

// MerchantInteractor 商户用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/merchant_interactor.go -package=mock . MerchantInteractor
type MerchantInteractor interface {
	CreateMerchant(ctx context.Context, domainCMerchant *CreateMerchantParams) (err error)
	CreateMerchantAndStore(ctx context.Context, domainMerchant *CreateMerchantParams, domainCStore *CreateStoreParams) (err error)
	UpdateMerchant(ctx context.Context, domainUMerchant *UpdateMerchantParams) (err error)
	UpdateMerchantAndStore(ctx context.Context, domainMerchant *UpdateMerchantParams, domainUStore *UpdateStoreParams) (err error)
	DeleteMerchant(ctx context.Context, id uuid.UUID) (err error)
	GetMerchant(ctx context.Context, id uuid.UUID) (domainMerchant *Merchant, err error)
	GetMerchants(ctx context.Context, pager *upagination.Pagination, filter *MerchantListFilter, orderBys ...MerchantListOrderBy) (domainMerchants []*Merchant, total int, err error)
	CountMerchant(ctx context.Context) (merchantCount *MerchantCount, err error)
	MerchantRenewal(ctx context.Context, merchantRenewal *MerchantRenewal) (err error)
	MerchantSimpleUpdate(ctx context.Context, updateField MerchantSimpleUpdateType, domainMerchant *Merchant) (err error)
}
type MerchantListOrderByType int

const (
	_ MerchantListOrderByType = iota
	MerchantListOrderByID
	MerchantListOrderByCreatedAt
)

type MerchantListOrderBy struct {
	OrderBy MerchantListOrderByType
	Desc    bool
}

func NewMerchantListOrderByID(desc bool) MerchantListOrderBy {
	return MerchantListOrderBy{
		OrderBy: MerchantListOrderByID,
		Desc:    desc,
	}
}

func NewMerchantListOrderByCreatedAt(desc bool) MerchantListOrderBy {
	return MerchantListOrderBy{
		OrderBy: MerchantListOrderByCreatedAt,
		Desc:    desc,
	}
}

type MerchantType string

const (
	MerchantTypeBrand MerchantType = "brand" // 品牌商户
	MerchantTypeStore MerchantType = "store" // 门店商户
)

func (t MerchantType) Values() []string {
	return []string{
		string(MerchantTypeBrand),
		string(MerchantTypeStore),
	}
}

func (t MerchantType) ToString() string {
	switch t {
	case MerchantTypeBrand:
		return "品牌商户"
	case MerchantTypeStore:
		return "门店商户"
	default:
		return ""
	}
}

type MerchantStatus string

const (
	MerchantStatusActive   MerchantStatus = "active"   // 已激活
	MerchantStatusExpired  MerchantStatus = "expired"  // 已过期
	MerchantStatusDisabled MerchantStatus = "disabled" // 已禁用
)

func (MerchantStatus) Values() []string {
	return []string{
		string(MerchantStatusActive),
		string(MerchantStatusExpired),
		string(MerchantStatusDisabled),
	}
}

func (s MerchantStatus) ToString() string {
	switch s {
	case MerchantStatusActive:
		return "已激活"
	case MerchantStatusExpired:
		return "已过期"
	case MerchantStatusDisabled:
		return "已禁用"
	default:
		return ""
	}
}

type MerchantSimpleUpdateType string

const (
	MerchantSimpleUpdateTypeStatus MerchantSimpleUpdateType = "status" // 状态
)

type Merchant struct {
	ID                   uuid.UUID            `json:"id"`
	MerchantCode         string               `json:"merchant_code"`                             // 商户编号(保留字段)
	MerchantName         string               `json:"merchant_name"`                             // 商户名称,最长不得超过50个字
	MerchantShortName    string               `json:"merchant_short_name"`                       // 商户简称
	MerchantType         MerchantType         `json:"merchant_type"`                             // 商户类型: 品牌商户,门店商户
	BrandName            string               `json:"brand_name"`                                // 品牌名称
	AdminPhoneNumber     string               `json:"admin_phone_number"`                        // 管理员手机号
	ExpireUTC            *time.Time           `json:"expire_utc"`                                // UTC 时区的过期时间
	BusinessTypeID       uuid.UUID            `json:"business_type_id"`                          // 业务类型
	BusinessTypeName     string               `json:"business_type_name"`                        // 业务类型名称
	MerchantLogo         string               `json:"merchant_logo"`                             // logo 图片地址
	Description          string               `json:"description"`                               // 商户描述(保留字段)
	Status               MerchantStatus       `json:"status"`                                    // 状态: 正常,停用,过期
	LoginAccount         string               `json:"login_account"`                             // 登录账号
	LoginPassword        string               `json:"login_password"`                            // 登录密码(加密存储)
	Address              *Address             `json:"address"`                                   // 地址
	StoreCount           int                  `json:"store_count"`                               // 关联门店数量(仅品牌商户有效)
	PurchaseDuration     int                  `json:"purchase_duration" binding:"required"`      // 购买时长
	PurchaseDurationUnit PurchaseDurationUnit `json:"purchase_duration_unit" binding:"required"` // 购买时长单位

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MerchantSimple struct {
	ID           uuid.UUID `json:"id"`
	MerchantName string    `json:"merchant_name"` // 商户名称
}

type Address struct {
	CountryID    uuid.UUID `json:"country_id"`    // 国家/地区 ID
	ProvinceID   uuid.UUID `json:"province_id"`   // 省份 ID
	CityID       uuid.UUID `json:"city_id"`       // 城市 ID
	DistrictID   uuid.UUID `json:"district_id"`   // 区县 ID
	CountryName  string    `json:"country_name"`  // 国家/地区 名称
	ProvinceName string    `json:"province_name"` // 省份名称
	CityName     string    `json:"city_name"`     // 城市名称
	DistrictName string    `json:"district_name"` // 区县名称
	Address      string    `json:"address"`       // 详细地址
	Lng          string    `json:"lng"`           // 经度
	Lat          string    `json:"lat"`           // 纬度
}
type MerchantCount struct {
	MerchantTypeBrand int `json:"merchant_type_brand"` // 品牌商户数量
	MerchantTypeStore int `json:"merchant_type_store"` // 门店商户数量
	Expired           int `json:"expired"`             // 过期商户数量
}

type MerchantListFilter struct {
	Status           MerchantStatus `json:"status"`
	MerchantName     string         `json:"merchant_name"`      // 商户名称
	AdminPhoneNumber string         `json:"admin_phone_number"` // 管理员手机号
	MerchantType     MerchantType   `json:"merchant_type"`      // 商户类型: 品牌商户,门店商户
	CreatedAtGte     *time.Time     `json:"created_at_gte"`     // 创建时间 大于等于
	CreatedAtLte     *time.Time     `json:"created_at_lte"`     // 创建时间 小于等于
	ProvinceID       uuid.UUID      `json:"province_id"`        // 省份 ID
}

type MerchantExistsParams struct {
	MerchantName string    // 商户名称
	ExcludeID    uuid.UUID // 排除的商户 ID
}

type CreateMerchantParams struct {
	MerchantCode         string               `json:"merchant_code"`          // 商户编号(保留字段)
	MerchantName         string               `json:"merchant_name"`          // 商户名称,最长不得超过50个字
	MerchantShortName    string               `json:"merchant_short_name"`    // 商户简称
	MerchantType         MerchantType         `json:"merchant_type"`          // 商户类型: 品牌商户,门店商户
	BrandName            string               `json:"brand_name"`             // 品牌名称
	AdminPhoneNumber     string               `json:"admin_phone_number"`     // 管理员手机号
	PurchaseDuration     int                  `json:"purchase_duration"`      // 购买时长
	PurchaseDurationUnit PurchaseDurationUnit `json:"purchase_duration_unit"` // 购买时长单位
	BusinessTypeID       uuid.UUID            `json:"business_type_id"`       // 业务类型
	MerchantLogo         string               `json:"merchant_logo"`          // logo 图片地址
	Description          string               `json:"description"`            // 商户描述(保留字段)
	LoginAccount         string               `json:"login_account"`          // 登录账号
	LoginPassword        string               `json:"login_password"`         // 登录密码(加密存储)
	Address              *Address             `json:"address"`                // 地址
}

type UpdateMerchantParams struct {
	ID                uuid.UUID `json:"id"`
	MerchantCode      string    `json:"merchant_code"`       // 商户编号(保留字段)
	MerchantName      string    `json:"merchant_name"`       // 商户名称,最长不得超过50个字
	MerchantShortName string    `json:"merchant_short_name"` // 商户简称
	BrandName         string    `json:"brand_name"`          // 品牌名称
	AdminPhoneNumber  string    `json:"admin_phone_number"`  // 管理员手机号
	BusinessTypeID    uuid.UUID `json:"business_type_id"`    // 业务类型
	MerchantLogo      string    `json:"merchant_logo"`       // logo 图片地址
	Description       string    `json:"description"`         // 商户描述(保留字段)
	Address           *Address  `json:"address"`             // 地址
}
