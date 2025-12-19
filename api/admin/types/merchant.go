package types

import (
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
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
type CreateStoreMerchantReq struct {
	Merchant CreateMerchantReq
	Store    CreateStoreReq
}

type CreateMerchantReq struct {
	MerchantCode         string                      `json:"merchant_code" binding:"omitempty,max=50"`       // 商户编号(保留字段)
	MerchantName         string                      `json:"merchant_name" binding:"required,max=50"`        // 商户名称,最长不得超过50个字
	MerchantShortName    string                      `json:"merchant_short_name" binding:"omitempty,max=50"` // 商户简称
	MerchantType         domain.MerchantType         `json:"merchant_type" binding:"required"`               // 商户类型: 品牌商户,门店商户
	BrandName            string                      `json:"brand_name" binding:"omitempty,max=50"`          // 品牌名称
	AdminPhoneNumber     string                      `json:"admin_phone_number" binding:"required,max=20"`   // 管理员手机号
	PurchaseDuration     int                         `json:"purchase_duration" binding:"required"`           // 购买时长
	PurchaseDurationUnit domain.PurchaseDurationUnit `json:"purchase_duration_unit" binding:"required"`      // 购买时长单位
	BusinessTypeID       uuid.UUID                   `json:"business_type_id" binding:"required"`            // 业务类型
	MerchantLogo         string                      `json:"merchant_logo" binding:"omitempty,max=500"`      // logo 图片地址
	Description          string                      `json:"description" binding:"omitempty,max=255"`        // 商户描述(保留字段)
	Status               domain.MerchantStatus       `json:"status" binding:"omitempty"`                     // 状态: 正常,停用,过期
	LoginAccount         string                      `json:"login_account" binding:"required"`               // 登录账号
	LoginPassword        string                      `json:"login_password" binding:"required"`              // 登录密码(加密存储)

	Address Address `json:"address" binding:"omitempty"` // 地址
}

type UpdateStoreMerchantReq struct {
	Merchant UpdateMerchantReq
	Store    UpdateStoreReq
}

type UpdateMerchantReq struct {
	MerchantCode      string                `json:"merchant_code" binding:"omitempty,max=50"`       // 商户编号(保留字段)
	MerchantName      string                `json:"merchant_name" binding:"required,max=50"`        // 商户名称,最长不得超过50个字
	MerchantShortName string                `json:"merchant_short_name" binding:"omitempty,max=50"` // 商户简称
	BrandName         string                `json:"brand_name" binding:"omitempty,max=50"`          // 品牌名称
	AdminPhoneNumber  string                `json:"admin_phone_number" binding:"required,max=20"`   // 管理员手机号
	BusinessTypeID    uuid.UUID             `json:"business_type_id" binding:"required"`            // 业务类型
	MerchantLogo      string                `json:"merchant_logo" binding:"omitempty,max=500"`      // logo 图片地址
	Description       string                `json:"description" binding:"omitempty,max=255"`        // 商户描述(保留字段)
	Status            domain.MerchantStatus `json:"status" binding:"omitempty"`                     // 状态: 正常,停用,过期
	LoginAccount      string                `json:"login_account" binding:"required"`               // 登录账号
	LoginPassword     string                `json:"login_password" binding:"required"`              // 登录密码(加密存储)
	Address           Address               `json:"address" binding:"omitempty"`                    // 地址
}

type MerchantInfoResp struct {
	Merchant *domain.Merchant `json:"merchant"` // 商户信息
}

type MerchantListReq struct {
	upagination.RequestPagination
	MerchantName     string                `json:"merchant_name" binding:"omitempty"`                        // 商户名称
	AdminPhoneNumber string                `json:"admin_phone_number" binding:"omitempty"`                   // 管理员手机号
	MerchantType     domain.MerchantType   `json:"merchant_type" binding:"omitempty,oneof=brand store"`      // 商户类型: 品牌商户,门店商户
	Status           domain.MerchantStatus `json:"status" binding:"omitempty,oneof=active expired disabled"` // 状态: 正常,停用,过期
	ProvinceID       uuid.UUID             `json:"province_id" binding:"omitempty"`                          // 省份 ID
	CreatedAtGte     time.Time             `json:"created_at_gte" binding:"omitempty"`
	CreatedAtLte     time.Time             `json:"created_at_lte" binding:"omitempty"`
}

type MerchantListResp struct {
	Merchants []*domain.Merchant `json:"merchants"` // 商户列表
	Total     int                `json:"total"`     // 总数
}
