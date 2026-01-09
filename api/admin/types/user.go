package types

import (
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

// AdminUserCreateReq 创建管理员用户请求
type AdminUserCreateReq struct {
	Username     string        `json:"username" binding:"required"`       // 用户登陆账号
	Password     string        `json:"password" binding:"required"`       // 用户登陆密码
	Nickname     string        `json:"nickname" binding:"omitempty"`      // 昵称
	DepartmentID uuid.UUID     `json:"department_id" binding:"required"`  // 所属部门ID
	RealName     string        `json:"real_name" binding:"required"`      // 用户姓名
	Gender       domain.Gender `json:"gender" binding:"omitempty"`        // 性别
	Email        string        `json:"email" binding:"omitempty,email"`   // 电子邮箱
	PhoneNumber  string        `json:"phone_number" binding:"omitempty"`  // 手机号
	Enabled      bool          `json:"enabled"`                           // 是否启用
	RoleIDs      []uuid.UUID   `json:"role_ids" binding:"required,min=1"` // 角色ID列表
}

// AdminUserUpdateReq 更新管理员用户请求
// Password is optional; if empty, the password remains unchanged.
type AdminUserUpdateReq struct {
	Username     string        `json:"username" binding:"required"`       // 用户登陆账号
	Password     string        `json:"password" binding:"required"`       // 用户登陆密码
	Nickname     string        `json:"nickname" binding:"omitempty"`      // 昵称
	DepartmentID uuid.UUID     `json:"department_id" binding:"required"`  // 所属部门ID
	RealName     string        `json:"real_name" binding:"required"`      // 用户姓名
	Gender       domain.Gender `json:"gender" binding:"omitempty"`        // 性别
	Email        string        `json:"email" binding:"omitempty,email"`   // 电子邮箱
	PhoneNumber  string        `json:"phone_number" binding:"omitempty"`  // 手机号
	Enabled      bool          `json:"enabled"`                           // 是否启用
	RoleIDs      []uuid.UUID   `json:"role_ids" binding:"required,min=1"` // 角色ID列表
}

// AdminUserListReq 管理员用户列表查询请求
type AdminUserListReq struct {
	upagination.RequestPagination
	Code        string        `form:"code"`         // 编号
	RealName    string        `form:"real_name"`    // 真实姓名
	Gender      domain.Gender `form:"gender"`       // 性别
	Email       string        `form:"email"`        // 电子邮箱
	PhoneNumber string        `form:"phone_number"` // 手机号
	Enabled     *bool         `form:"enabled"`      // 是否启用
	RoleID      string        `form:"role_id"`      // 角色ID
}

// AdminUserListResp 管理员用户列表响应
type AdminUserListResp struct {
	Users []*domain.AdminUser `json:"users"`
	Total int                 `json:"total"`
}

type ResetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required"`
}
