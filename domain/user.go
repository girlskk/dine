package domain

import (
	"github.com/google/uuid"
)

// User 通用用户接口，用于验证用户身份
type User interface {
	GetMerchantID() uuid.UUID
	GetStoreID() uuid.UUID
}

// VerifyOwnerShip 验证资源是否属于当前用户可操作
func VerifyOwnerShip(user User, merchantID, storeID uuid.UUID) bool {
	if user.GetMerchantID() != merchantID || user.GetStoreID() != storeID {
		return false
	}
	return true
}

// 性别
type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
	GenderUnknown Gender = "unknown"
)

func (Gender) Values() []string {
	return []string{string(GenderMale), string(GenderFemale), string(GenderOther), string(GenderUnknown)}
}
