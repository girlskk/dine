package domain

import "github.com/google/uuid"

// User 通用用户接口，用于验证用户身份
type User interface {
	GetMerchantID() uuid.UUID
	GetStoreID() uuid.UUID
}
