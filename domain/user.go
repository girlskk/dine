package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotExists    = errors.New("用户不存在")
	ErrUsernameExist    = errors.New("用户名已存在")
	ErrUserRoleRequired = errors.New("管理员用户至少需要分配一个角色")
)

// User 通用用户接口，用于验证用户身份
type User interface {
	GetUserID() uuid.UUID
	GetMerchantID() uuid.UUID
	GetStoreID() uuid.UUID
	GetUserType() UserType
}

// 性别
type Gender string

const (
	GenderMale    Gender = "male"    // 男性
	GenderFemale  Gender = "female"  // 女性
	GenderOther   Gender = "other"   // 其他
	GenderUnknown Gender = "unknown" // 未知
)

func (Gender) Values() []string {
	return []string{string(GenderMale), string(GenderFemale), string(GenderOther), string(GenderUnknown)}
}

type UserType string

const (
	UserTypeAdmin   UserType = "admin"   // admin表用户
	UserTypeBackend UserType = "backend" // backend用户
	UserTypeStore   UserType = "store"   // store用户
)

func (UserType) Values() []string {
	return []string{string(UserTypeAdmin), string(UserTypeBackend), string(UserTypeStore)}
}
