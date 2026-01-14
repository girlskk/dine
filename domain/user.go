package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// 用户相关
	ErrUserNotExists             = errors.New("用户不存在")
	ErrUsernameExist             = errors.New("用户账号已存在")
	ErrSuperUserCannotDelete     = errors.New("超级管理员无法删除")
	ErrSuperUserCannotDisable    = errors.New("超级管理员不能被禁用")
	ErrSuperUserCannotUpdate     = errors.New("超级管理员不能编辑")
	ErrUserDisabled              = errors.New("用户已被禁用")
	ErrPasswordCannotBeEmpty     = errors.New("密码不能为空")
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")

	// 部门相关
	ErrDepartmentDisabled = errors.New("用户所属部门已被禁用")

	// 角色相关
	ErrRoleDisabled   = errors.New("用户所属角色已被禁用")
	ErrRoleNotExists  = errors.New("角色不存在")
	ErrRoleNameExists = errors.New("角色名称已存在")
	ErrRoleCodeExists = errors.New("角色编码已存在")

	// 用户和部门关系
	ErrUserDepartmentRequired     = errors.New("用户所属部门不能为空")
	ErrUserDepartmentTypeMismatch = errors.New("用户部门类型不匹配")

	// 角色和权限关系
	ErrUserRoleNotExists         = errors.New("该用户未分配角色")
	ErrUserRoleRequired          = errors.New("用户至少需要分配一个角色")
	ErrUserRoleTypeMismatch      = errors.New("用户角色类型不匹配")
	ErrRoleAssignedCannotDisable = errors.New("角色已分配用户，无法禁用")
	ErrRoleAssignedCannotDelete  = errors.New("角色已分配用户，无法删除")

	// 前端路由菜单
	ErrRouterMenuNotExists        = errors.New("菜单不存在")
	ErrRouterMenuNameExists       = errors.New("同级菜单名称已存在")
	ErrRouterMenuForbidenAddChild = errors.New("禁止添加子菜单")

	// 角色和菜单关系
	ErrRoleMenuNotExists = errors.New("角色菜单关系不存在")
	ErrRoleMenuExists    = errors.New("角色菜单关系已存在")
)

// User 通用用户接口，用于验证用户身份
type User interface {
	GetUserID() uuid.UUID
	GetMerchantID() uuid.UUID
	GetStoreID() uuid.UUID
	GetUserType() UserType
}

func VerifyOwnerMerchant(user User, merchantID uuid.UUID) bool {
	if user.GetMerchantID() != merchantID {
		return false
	}
	return true
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
