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

type BackendUserCreateReq struct {
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

type BackendUserUpdateReq struct {
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

type AccountListReq struct {
	upagination.RequestPagination
	Code        string        `form:"code"`
	RealName    string        `form:"real_name"`
	Gender      domain.Gender `form:"gender"`
	Email       string        `form:"email"`
	PhoneNumber string        `form:"phone_number"`
	Enabled     *bool         `form:"enable"`
}

type AccountListResp struct {
	Users []*domain.BackendUser `json:"users"`
	Total int                   `json:"total"`
}

type ResetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required"`
}
