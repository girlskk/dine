package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_user_repository.go -package=mock . StoreUserRepository
type StoreUserRepository interface {
	Create(ctx context.Context, user *StoreUser) error
	FindByUsername(ctx context.Context, username string) (*StoreUser, error)
	Find(ctx context.Context, id uuid.UUID) (*StoreUser, error)
	Update(ctx context.Context, user *StoreUser) error
	Exists(ctx context.Context, params StoreUserExistsParams) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetUsers(ctx context.Context, pager *upagination.Pagination, filter *StoreUserListFilter, orderBys ...StoreUserOrderBy) (users []*StoreUser, total int, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_user_interactor.go -package=mock . StoreUserInteractor
type StoreUserInteractor interface {
	Login(ctx context.Context, username string, password string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *StoreUser, err error)

	Create(ctx context.Context, user *StoreUser) error
	Update(ctx context.Context, user *StoreUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUser(ctx context.Context, id uuid.UUID) (*StoreUser, error)
	GetUsers(ctx context.Context, pager *upagination.Pagination, filter *StoreUserListFilter, orderBys ...StoreUserOrderBy) (users []*StoreUser, total int, err error)
}

type StoreUser struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	Nickname       string    `json:"nickname"`
	MerchantID     uuid.UUID `json:"merchant_id"`
	StoreID        uuid.UUID `json:"store_id"`      // 门店ID
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
	Store      *Store      `json:"store,omitempty"`      // 所属门店
}

// GetUserID 实现 User 接口
func (u *StoreUser) GetUserID() uuid.UUID {
	return u.ID
}

// GetMerchantID 实现 User 接口
func (u *StoreUser) GetMerchantID() uuid.UUID {
	return u.MerchantID
}

// GetStoreID 实现 User 接口（门店用户的 StoreID 为 uuid.Nil）
func (u *StoreUser) GetStoreID() uuid.UUID {
	return u.StoreID
}

// GetUserType 实现 User 接口 (门店用户的UserType为UserTypeStore)
func (u *StoreUser) GetUserType() UserType {
	return UserTypeStore
}

func (u *StoreUser) SetPassword(password string) error {
	hashed, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	u.HashedPassword = hashed
	return nil
}

func (u *StoreUser) CheckPassword(password string) error {
	if err := util.CheckPassword(password, u.HashedPassword); err != nil {
		return ErrMismatchedHashAndPassword
	}
	return nil
}

type (
	storeUserKey struct{}
)

func NewStoreUserContext(ctx context.Context, u *StoreUser) context.Context {
	return context.WithValue(ctx, storeUserKey{}, u)
}

func FromStoreUserContext(ctx context.Context) *StoreUser {
	if v, ok := ctx.Value(storeUserKey{}).(*StoreUser); ok {
		return v
	}
	return nil
}

type StoreUserOrderByType int

const (
	_ StoreUserOrderByType = iota
	StoreUserOrderByID
	StoreUserOrderByCreatedAt
)

type StoreUserOrderBy struct {
	OrderBy StoreUserOrderByType
	Desc    bool
}

func NewStoreUserOrderByID(desc bool) StoreUserOrderBy {
	return StoreUserOrderBy{OrderBy: StoreUserOrderByID, Desc: desc}
}

func NewStoreUserOrderByCreatedAt(desc bool) StoreUserOrderBy {
	return StoreUserOrderBy{OrderBy: StoreUserOrderByCreatedAt, Desc: desc}
}

type StoreUserListFilter struct {
	UserIDs     []uuid.UUID `json:"user_ids"`     // 用户ID列表
	Code        string      `json:"code"`         // 编号
	RealName    string      `json:"real_name"`    // 真实姓名
	Gender      Gender      `json:"gender"`       // 性别
	Email       string      `json:"email"`        // 电子邮箱
	PhoneNumber string      `json:"phone_number"` // 手机号
	Enabled     *bool       `json:"enabled"`      // 是否启用
	MerchantID  uuid.UUID   `json:"merchant_id"`  // 品牌商ID
	StoreID     uuid.UUID   `json:"store_id"`     // 门店ID
}
type StoreUserExistsParams struct {
	Username  string
	ExcludeID uuid.UUID
}
