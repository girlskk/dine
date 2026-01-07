package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var (
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
	ErrAdminUserNotExists        = errors.New("管理员用户不存在")
	ErrAdminUserUsernameExist    = errors.New("用户名已存在")
	ErrBackendUserRoleRequired   = errors.New("管理员用户至少需要分配一个角色")
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/admin_user_repository.go -package=mock . AdminUserRepository
type AdminUserRepository interface {
	FindByUsername(ctx context.Context, username string) (*AdminUser, error)
	Find(ctx context.Context, id uuid.UUID) (*AdminUser, error)
	Create(ctx context.Context, user *AdminUser) error
	Update(ctx context.Context, user *AdminUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, params AdminUserExistsParams) (bool, error)
	GetUsers(ctx context.Context, pager *upagination.Pagination, filter *AdminUserListFilter, orderBys ...AdminUserOrderBy) (users []*AdminUser, total int, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/admin_user_interactor.go -package=mock . AdminUserInteractor
type AdminUserInteractor interface {
	Login(ctx context.Context, username string, password string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *AdminUser, err error)

	Create(ctx context.Context, user *AdminUser) error
	Update(ctx context.Context, user *AdminUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUser(ctx context.Context, id uuid.UUID) (*AdminUser, error)
	GetUsers(ctx context.Context, pager *upagination.Pagination, filter *AdminUserListFilter, orderBys ...AdminUserOrderBy) (users []*AdminUser, total int, err error)
}

type AdminUser struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`      // 用户名
	HashedPassword string    `json:"-"`             // 密码哈希
	Nickname       string    `json:"nickname"`      // 昵称
	DepartmentID   uuid.UUID `json:"department_id"` // 所属部门ID
	Code           string    `json:"code"`          // 编号
	RealName       string    `json:"real_name"`     // 真实姓名
	Gender         Gender    `json:"gender"`        // 性别
	Email          string    `json:"email"`         // 电子邮箱
	PhoneNumber    string    `json:"phone_number"`  // 手机号
	Enabled        bool      `json:"enabled"`       // 是否启用
	IsSuperAdmin   bool      `json:"is_superadmin"` // 是否为超级管理员
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	RoleIDs    []uuid.UUID `json:"role_ids,omitempty"`   // 角色ID列表
	RoleList   []*Role     `json:"role_list,omitempty"`  // 角色列表
	Department *Department `json:"department,omitempty"` // 所属部门
}

// GetUserID 实现 User 接口
func (u *AdminUser) GetUserID() uuid.UUID {
	return u.ID
}

// GetMerchantID 实现 User 接口 (管理员用户的 MerchantID 为 uuid.Nil)
func (u *AdminUser) GetMerchantID() uuid.UUID {
	return uuid.Nil
}

// GetStoreID 实现 User 接口 (管理员用户的 StoreID 为 uuid.Nil)
func (u *AdminUser) GetStoreID() uuid.UUID {
	return uuid.Nil
}

// GetUserType 实现 User 接口 (管理员用户的UserType为UserTypeAdmin)
func (u *AdminUser) GetUserType() UserType {
	return UserTypeAdmin
}

func (u *AdminUser) SetPassword(password string) error {
	hashed, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	u.HashedPassword = hashed
	return nil
}

func (u *AdminUser) CheckPassword(password string) error {
	if err := util.CheckPassword(password, u.HashedPassword); err != nil {
		return ErrMismatchedHashAndPassword
	}

	return nil
}

type (
	adminUserKey struct{}
)

func NewAdminUserContext(ctx context.Context, u *AdminUser) context.Context {
	return context.WithValue(ctx, adminUserKey{}, u)
}

func FromAdminUserContext(ctx context.Context) *AdminUser {
	if v, ok := ctx.Value(adminUserKey{}).(*AdminUser); ok {
		return v
	}
	return nil
}

type AdminUserListFilter struct {
	UserIDs     []uuid.UUID `json:"user_ids"`
	Code        string      `json:"code"`
	RealName    string      `json:"real_name"`
	Gender      Gender      `json:"gender"`
	Email       string      `json:"email"`
	PhoneNumber string      `json:"phone_number"`
	Enabled     *bool       `json:"enabled"`
}

type AdminUserOrderByType int

const (
	_ AdminUserOrderByType = iota
	AdminUserOrderByID
	AdminUserOrderByCreatedAt
)

type AdminUserOrderBy struct {
	OrderBy AdminUserOrderByType
	Desc    bool
}

func NewAdminUserOrderByID(desc bool) AdminUserOrderBy {
	return AdminUserOrderBy{
		OrderBy: AdminUserOrderByID,
		Desc:    desc,
	}
}

func NewAdminUserOrderByCreatedAt(desc bool) AdminUserOrderBy {
	return AdminUserOrderBy{
		OrderBy: AdminUserOrderByCreatedAt,
		Desc:    desc,
	}
}

type AdminUserExistsParams struct {
	Username  string    `json:"username"`
	Code      string    `json:"code"`
	ExcludeID uuid.UUID `json:"exclude_id"`
}
