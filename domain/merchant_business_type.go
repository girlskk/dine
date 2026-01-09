package domain

import (
	"context"

	"github.com/google/uuid"
)

// MerchantBusinessTypeRepository 业态类型仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/merchant_business_type_repository.go -package=mock . MerchantBusinessTypeRepository
type MerchantBusinessTypeRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (businessType *MerchantBusinessType, err error)
	GetAll(ctx context.Context) (list []*MerchantBusinessType, err error)
	FindByCode(ctx context.Context, typeCode string) (businessType *MerchantBusinessType, err error)
}

// MerchantBusinessTypeInteractor 业态类型用例接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/merchant_business_type_interactor.go -package=mock . MerchantBusinessTypeInteractor
type MerchantBusinessTypeInteractor interface {
	GetAll(ctx context.Context) (list []*MerchantBusinessType, err error)
}
type MerchantBusinessType struct {
	ID       uuid.UUID `json:"id"`
	TypeCode string    `json:"type_code"` // 业态类型编码（保留字段）
	TypeName string    `json:"type_name"` // 业态类型名称
}
