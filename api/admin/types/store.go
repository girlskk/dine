package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type CreateStoreReq struct {
	MerchantID              uuid.UUID              `json:"merchant_id" binding:"required"`                         // 所属商户 ID
	AdminPhoneNumber        string                 `json:"admin_phone_number" binding:"required,max=20"`           // 管理员手机号
	StoreName               string                 `json:"store_name" binding:"required,max=30"`                   // 门店名称,长度不超过30个字
	StoreShortName          string                 `json:"store_short_name" binding:"omitempty,max=30"`            // 门店简称
	StoreCode               string                 `json:"store_code" binding:"omitempty,max=50"`                  // 门店编码(保留字段)
	Status                  domain.StoreStatus     `json:"status" binding:"required"`                              // 营业/停业
	BusinessModel           domain.BusinessModel   `json:"business_model" binding:"required"`                      // 直营/加盟
	BusinessTypeID          uuid.UUID              `json:"business_type_id" binding:"required"`                    // 业态类型
	ContactName             string                 `json:"contact_name" binding:"omitempty,max=255"`               // 联系人
	ContactPhone            string                 `json:"contact_phone" binding:"omitempty,max=255"`              // 联系电话
	UnifiedSocialCreditCode string                 `json:"unified_social_credit_code" binding:"omitempty,max=255"` // 统一社会信用代码
	StoreLogo               string                 `json:"store_logo" binding:"omitempty,max=500"`                 // logo 图片地址
	BusinessLicenseURL      string                 `json:"business_license_url" binding:"omitempty,max=500"`       // 营业执照图片地址
	StorefrontURL           string                 `json:"storefront_url" binding:"omitempty,max=500"`             // 门店门头照片地址
	CashierDeskURL          string                 `json:"cashier_desk_url" binding:"omitempty,max=500"`           // 收银台照片地址
	DiningEnvironmentURL    string                 `json:"dining_environment_url" binding:"omitempty,max=500"`     // 就餐环境照片地址
	FoodOperationLicenseURL string                 `json:"food_operation_license_url" binding:"omitempty,max=500"` // 食品经营许可证照片地址
	Address                 Address                `json:"address" binding:"omitempty"`                            // 地址 todo 门店地址校验必填
	LoginAccount            string                 `json:"login_account" binding:"required"`                       // 登录账号
	LoginPassword           string                 `json:"login_password" binding:"required"`                      // 登录密码(加密存储)
	BusinessHours           []domain.BusinessHours `json:"business_hours" binding:"required"`                      // 营业时间段
	DiningPeriods           []domain.DiningPeriod  `json:"dining_periods" binding:"required"`                      // 就餐时段
	ShiftTimes              []domain.ShiftTime     `json:"shift_times" binding:"required"`                         // 班次时间
}

type UpdateStoreReq struct {
	AdminPhoneNumber        string                 `json:"admin_phone_number" binding:"required,max=20"`           // 管理员手机号
	StoreName               string                 `json:"store_name" binding:"required,max=30"`                   // 门店名称,长度不超过30个字
	StoreShortName          string                 `json:"store_short_name" binding:"omitempty,max=30"`            // 门店简称
	StoreCode               string                 `json:"store_code" binding:"omitempty,max=50"`                  // 门店编码(保留字段)
	Status                  domain.StoreStatus     `json:"status" binding:"required"`                              // 营业/停业
	BusinessModel           domain.BusinessModel   `json:"business_model" binding:"required"`                      // 直营/加盟
	BusinessTypeID          uuid.UUID              `json:"business_type_id" binding:"required"`                    // 业态类型
	ContactName             string                 `json:"contact_name" binding:"omitempty,max=255"`               // 联系人
	ContactPhone            string                 `json:"contact_phone" binding:"omitempty,max=255"`              // 联系电话
	UnifiedSocialCreditCode string                 `json:"unified_social_credit_code" binding:"omitempty,max=255"` // 统一社会信用代码
	StoreLogo               string                 `json:"store_logo" binding:"omitempty,max=500"`                 // logo 图片地址
	BusinessLicenseURL      string                 `json:"business_license_url" binding:"omitempty,max=500"`       // 营业执照图片地址
	StorefrontURL           string                 `json:"storefront_url" binding:"omitempty,max=500"`             // 门店门头照片地址
	CashierDeskURL          string                 `json:"cashier_desk_url" binding:"omitempty,max=500"`           // 收银台照片地址
	DiningEnvironmentURL    string                 `json:"dining_environment_url" binding:"omitempty,max=500"`     // 就餐环境照片地址
	FoodOperationLicenseURL string                 `json:"food_operation_license_url" binding:"omitempty,max=500"` // 食品经营许可证照片地址
	Address                 Address                `json:"address" binding:"omitempty"`                            // 地址 todo 门店地址校验必填
	LoginAccount            string                 `json:"login_account" binding:"required"`                       // 登录账号
	LoginPassword           string                 `json:"login_password" binding:"required"`                      // 登录密码(加密存储)
	BusinessHours           []domain.BusinessHours `json:"business_hours" binding:"required"`                      // 营业时间段
	DiningPeriods           []domain.DiningPeriod  `json:"dining_periods" binding:"required"`                      // 就餐时段
	ShiftTimes              []domain.ShiftTime     `json:"shift_times" binding:"required"`                         // 班次时间
}
