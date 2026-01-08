package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/backend_user_repository.go -package=mock . BackendUserRepository
type BackendUserRepository interface {
	Create(ctx context.Context, user *BackendUser) error
	FindByUsername(ctx context.Context, username string) (*BackendUser, error)
	Exists(ctx context.Context, params BackendUserExistsParams) (bool, error)
	Find(ctx context.Context, id uuid.UUID) (*BackendUser, error)
	Update(ctx context.Context, user *BackendUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUsers(ctx context.Context, pager *upagination.Pagination, filter *BackendUserListFilter, orderBys ...BackendUserOrderBy) (users []*BackendUser, total int, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/backend_user_interactor.go -package=mock . BackendUserInteractor
type BackendUserInteractor interface {
	Login(ctx context.Context, username string, password string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *BackendUser, err error)

	Create(ctx context.Context, user *BackendUser) error
	Update(ctx context.Context, user *BackendUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUser(ctx context.Context, id uuid.UUID) (*BackendUser, error)
	GetUsers(ctx context.Context, pager *upagination.Pagination, filter *BackendUserListFilter, orderBys ...BackendUserOrderBy) (users []*BackendUser, total int, err error)
	SimpleUpdate(ctx context.Context, updateField BackendUserSimpleUpdateField, params BackendUserSimpleUpdateParams) error
}

type BackendUser struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	Nickname       string    `json:"nickname"`
	MerchantID     uuid.UUID `json:"merchant_id"`   // 品牌商ID
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

	// 关联数据
	RoleIDs    []uuid.UUID `json:"role_ids,omitempty"`   // 角色ID列表
	RoleList   []*Role     `json:"role_list,omitempty"`  // 角色列表
	Department *Department `json:"department,omitempty"` // 所属部门
	Merchant   *Merchant   `json:"merchant,omitempty"`   // 所属品牌商
}

// GetUserID 实现 User 接口
func (u *BackendUser) GetUserID() uuid.UUID {
	return u.ID
}

// GetMerchantID 实现 User 接口
func (u *BackendUser) GetMerchantID() uuid.UUID {
	return u.MerchantID
}

// GetStoreID 实现 User 接口（品牌商用户的 StoreID 为 uuid.Nil）
func (u *BackendUser) GetStoreID() uuid.UUID {
	return uuid.Nil
}

// GetUserType 实现 User 接口 (品牌商用户的UserType为UserTypeBackend)
func (u *BackendUser) GetUserType() UserType {
	return UserTypeBackend
}

func (u *BackendUser) SetPassword(password string) error {
	hashed, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	u.HashedPassword = hashed
	return nil
}

func (u *BackendUser) CheckPassword(password string) error {
	if err := util.CheckPassword(password, u.HashedPassword); err != nil {
		return ErrMismatchedHashAndPassword
	}
	return nil
}

type (
	backendUserKey struct{}
)

func NewBackendUserContext(ctx context.Context, u *BackendUser) context.Context {
	return context.WithValue(ctx, backendUserKey{}, u)
}

func FromBackendUserContext(ctx context.Context) *BackendUser {
	if v, ok := ctx.Value(backendUserKey{}).(*BackendUser); ok {
		return v
	}
	return nil
}

type BackendUserListFilter struct {
	UserIDs     []uuid.UUID `json:"user_ids"`     //
	Code        string      `json:"code"`         // 编号
	RealName    string      `json:"real_name"`    // 真实姓名
	Gender      Gender      `json:"gender"`       // 性别
	Email       string      `json:"email"`        // 电子邮箱
	PhoneNumber string      `json:"phone_number"` // 手机号
	Enabled     *bool       `json:"enabled"`      // 是否启用
	MerchantID  uuid.UUID   `json:"merchant_id"`  // 品牌商ID
}

type BackendUserOrderByType int

const (
	_ BackendUserOrderByType = iota
	BackendUserOrderByID
	BackendUserOrderByCreatedAt
)

type BackendUserOrderBy struct {
	OrderBy BackendUserOrderByType
	Desc    bool
}

func NewBackendUserOrderByID(desc bool) BackendUserOrderBy {
	return BackendUserOrderBy{
		OrderBy: BackendUserOrderByID,
		Desc:    desc,
	}
}

func NewBackendUserOrderByCreatedAt(desc bool) BackendUserOrderBy {
	return BackendUserOrderBy{
		OrderBy: BackendUserOrderByCreatedAt,
		Desc:    desc,
	}
}

type BackendUserExistsParams struct {
	Username  string
	ExcludeID uuid.UUID
}

// BackendUserSimpleUpdateField Simple update types for backend user
type BackendUserSimpleUpdateField string

const (
	BackendUserSimpleUpdateFieldEnable BackendUserSimpleUpdateField = "enable"
)

type BackendUserSimpleUpdateParams struct {
	ID      uuid.UUID `json:"id"`
	Enabled bool      `json:"enabled"`
}
