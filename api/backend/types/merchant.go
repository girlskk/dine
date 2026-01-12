package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type Address struct {
	CountryID    uuid.UUID `json:"country_id" binding:"omitempty"`      // 国家/地区 ID
	ProvinceID   uuid.UUID `json:"province_id" binding:"omitempty"`     // 省份 ID
	CityID       uuid.UUID `json:"city_id" binding:"omitempty"`         // 城市 ID
	DistrictID   uuid.UUID `json:"district_id" binding:"omitempty"`     // 区县 ID
	CountryName  string    `json:"country_name" binding:"omitempty"`    // 国家/地区 名称
	ProvinceName string    `json:"province_name" binding:"omitempty"`   // 省份名称
	CityName     string    `json:"city_name" binding:"omitempty"`       // 城市名称
	DistrictName string    `json:"district_name" binding:"omitempty"`   // 区县名称
	Address      string    `json:"address" binding:"omitempty,max=100"` // 详细地址
	Lng          string    `json:"lng" binding:"omitempty,max=50"`      // 经度
	Lat          string    `json:"lat" binding:"omitempty,max=50"`      // 纬度
}

type UpdateStoreMerchantReq struct {
	Merchant UpdateMerchantReq
	Store    UpdateMStoreReq
}

type UpdateMerchantReq struct {
	MerchantCode      string                `json:"merchant_code" binding:"omitempty,max=50"`       // 商户编号(保留字段)
	MerchantName      string                `json:"merchant_name" binding:"required,max=50"`        // 商户名称,最长不得超过50个字
	MerchantShortName string                `json:"merchant_short_name" binding:"omitempty,max=50"` // 商户简称
	BrandName         string                `json:"brand_name" binding:"omitempty,max=50"`          // 品牌名称
	AdminPhoneNumber  string                `json:"admin_phone_number" binding:"required,max=20"`   // 管理员手机号
	BusinessTypeCode  string                `json:"business_type_code" binding:"required"`          // 业务类型
	MerchantLogo      string                `json:"merchant_logo" binding:"omitempty,max=500"`      // logo 图片地址
	Description       string                `json:"description" binding:"omitempty,max=255"`        // 商户描述(保留字段)
	Status            domain.MerchantStatus `json:"status" binding:"omitempty"`                     // 状态: 正常,停用,过期
	Address           Address               `json:"address" binding:"omitempty"`                    // 地址
}

// UpdateMStoreReq 更新门店商户中门店请求参数
type UpdateMStoreReq struct {
	AdminPhoneNumber        string                 `json:"admin_phone_number" binding:"required,max=20"`           // 管理员手机号
	StoreName               string                 `json:"store_name" binding:"required,max=30"`                   // 门店名称,长度不超过30个字
	StoreShortName          string                 `json:"store_short_name" binding:"omitempty,max=30"`            // 门店简称
	StoreCode               string                 `json:"store_code" binding:"omitempty,max=50"`                  // 门店编码(保留字段)
	Status                  domain.StoreStatus     `json:"status" binding:"required"`                              // 营业/停业
	LocationNumber          string                 `json:"location_number" binding:"required,max=255"`             // 门店位置编号
	ContactName             string                 `json:"contact_name" binding:"omitempty,max=255"`               // 联系人
	ContactPhone            string                 `json:"contact_phone" binding:"omitempty,max=255"`              // 联系电话
	UnifiedSocialCreditCode string                 `json:"unified_social_credit_code" binding:"omitempty,max=255"` // 统一社会信用代码
	StoreLogo               string                 `json:"store_logo" binding:"omitempty,max=500"`                 // logo 图片地址
	BusinessLicenseURL      string                 `json:"business_license_url" binding:"omitempty,max=500"`       // 营业执照图片地址
	StorefrontURL           string                 `json:"storefront_url" binding:"omitempty,max=500"`             // 门店门头照片地址
	CashierDeskURL          string                 `json:"cashier_desk_url" binding:"omitempty,max=500"`           // 收银台照片地址
	DiningEnvironmentURL    string                 `json:"dining_environment_url" binding:"omitempty,max=500"`     // 就餐环境照片地址
	FoodOperationLicenseURL string                 `json:"food_operation_license_url" binding:"omitempty,max=500"` // 食品经营许可证照片地址
	BusinessHours           []domain.BusinessHours `json:"business_hours" binding:"required"`                      // 营业时间段
	DiningPeriods           []domain.DiningPeriod  `json:"dining_periods" binding:"required"`                      // 就餐时段
	ShiftTimes              []domain.ShiftTime     `json:"shift_times" binding:"required"`                         // 班次时间
}
type MerchantInfoResp struct {
	Merchant *domain.Merchant `json:"merchant"` // 商户信息
	Store    *domain.Store    `json:"store"`    // 门店信息
}

type MerchantRenewalReq struct {
	MerchantID           uuid.UUID                   `json:"merchant_id" binding:"required"`            // 商户 ID
	PurchaseDuration     int                         `json:"purchase_duration" binding:"required"`      // 购买时长
	PurchaseDurationUnit domain.PurchaseDurationUnit `json:"purchase_duration_unit" binding:"required"` // 购买时长单位
}

type MerchantBusinessTypeListResp struct {
	BusinessTypes []*domain.BusinessType `json:"business_types"` // 业务类型列表
}
