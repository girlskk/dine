package domain

import "context"

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/merchant_business_type_repository.go -package=mock . MerchantBusinessTypeRepository
type MerchantBusinessTypeRepository interface {
	FindById(ctx context.Context, id int) (businessType *MerchantBusinessType, err error)
	GetAll(ctx context.Context) (ts []*MerchantBusinessType, err error)
}

type MerchantBusinessType struct {
	ID       int    `json:"id"`
	TypeCode string `json:"type_code"` // 业态类型编码（保留字段）
	TypeName string `json:"type_name"` // 业态类型名称
}
