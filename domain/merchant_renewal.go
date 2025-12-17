package domain

import (
	"context"
	"time"
)

type PurchaseDurationUnit string

const (
	PurchaseDurationUnitDay   PurchaseDurationUnit = "day"   // 日
	PurchaseDurationUnitMonth PurchaseDurationUnit = "month" // 月
	PurchaseDurationUnitYear  PurchaseDurationUnit = "year"  // 年
	PurchaseDurationUnitWeek  PurchaseDurationUnit = "week"  // 周
)

func (u PurchaseDurationUnit) Values() []string {
	return []string{
		string(PurchaseDurationUnitDay),
		string(PurchaseDurationUnitMonth),
		string(PurchaseDurationUnitYear),
		string(PurchaseDurationUnitWeek),
	}
}

func (u PurchaseDurationUnit) ToString() string {
	switch u {
	case PurchaseDurationUnitDay:
		return "日"
	case PurchaseDurationUnitMonth:
		return "月"
	case PurchaseDurationUnitYear:
		return "年"
	case PurchaseDurationUnitWeek:
		return "周"
	default:
		return ""
	}
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/merchant_renewal_repository.go -package=mock . MerchantRenewalRepository
type MerchantRenewalRepository interface {
	GetByMerchant(ctx context.Context, merchantId int) (renewals []*MerchantRenewal, err error)
	Create(ctx context.Context, merchantRenewal *MerchantRenewal) (err error)
}

type MerchantRenewal struct {
	ID                   int                  `json:"id"`
	MerchantID           int                  `json:"merchant_id"`            // 商户 ID
	PurchaseDuration     int                  `json:"purchase_duration"`      // 购买时长
	PurchaseDurationUnit PurchaseDurationUnit `json:"purchase_duration_unit"` // 购买时长单位
	OperatorName         string               `json:"operator_name"`          // 操作人
	OperatorAccount      string               `json:"operator_account"`       // 操作人账号
	CreatedAt            *time.Time           `json:"created_at"`             // 创建时间
}
