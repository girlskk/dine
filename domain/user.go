package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotExists              = errors.New("用户不存在")
	ErrUsernameExist              = errors.New("用户名已存在")
	ErrSuperUserCannotDelete      = errors.New("超级管理员无法删除")
	ErrSuperUserCannotDisable     = errors.New("超级管理员不能被禁用")
	ErrSuperUserCannotUpdate      = errors.New("超级管理员不能编辑")
	ErrUserDisabled               = errors.New("用户已被禁用")
	ErrDepartmentDisabled         = errors.New("用户所属部门已被禁用")
	ErrRoleDisabled               = errors.New("用户所属角色已被禁用")
	ErrUserRoleRequired           = errors.New("用户至少需要分配一个角色")
	ErrUserDepartmentRequired     = errors.New("用户所属部门不能为空")
	ErrUserRoleTypeMismatch       = errors.New("用户角色类型不匹配")
	ErrUserDepartmentTypeMismatch = errors.New("用户部门类型不匹配")
)

// User 通用用户接口，用于验证用户身份
type User interface {
	GetUserID() uuid.UUID
	GetMerchantID() uuid.UUID
	GetStoreID() uuid.UUID
	GetUserType() UserType
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
