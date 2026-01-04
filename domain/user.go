package domain

import "github.com/google/uuid"

// User 通用用户接口，用于验证用户身份
type User interface {
	GetMerchantID() uuid.UUID
	GetStoreID() uuid.UUID
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
