package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var (
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/admin_user_repository.go -package=mock . AdminUserRepository
type AdminUserRepository interface {
	FindByUsername(ctx context.Context, username string) (*AdminUser, error)
	Find(ctx context.Context, id uuid.UUID) (*AdminUser, error)
	Create(ctx context.Context, user *AdminUser) error
	Update(ctx context.Context, user *AdminUser) error
	Delete(ctx context.Context, id uuid.UUID) error
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/admin_user_interactor.go -package=mock . AdminUserInteractor
type AdminUserInteractor interface {
	Login(ctx context.Context, username string, password string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *AdminUser, err error)
}

type AdminUser struct {
	ID             uuid.UUID            `json:"id"`
	Username       string               `json:"username"`     // 用户名
	HashedPassword string               `json:"-"`            // 密码哈希
	Nickname       string               `json:"nickname"`     // 昵称
	AccountType    AdminUserAccountType `json:"account_type"` // 账户类型 // normal: 普通管理员, super_admin: 超级管理员
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

type AdminUserAccountType string

const (
	AdminUserAccountTypeNormal     AdminUserAccountType = "normal"      // 普通管理员
	AdminUserAccountTypeSuperAdmin AdminUserAccountType = "super_admin" // 超级管理员
)
