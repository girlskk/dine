package domain

import (
	"context"
	"errors"
	"time"

	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var (
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/frontend_user_repository.go -package=mock . FrontendUserRepository
type FrontendUserRepository interface {
	FindByUsername(ctx context.Context, username string) (*FrontendUser, error)
	Find(ctx context.Context, id int) (*FrontendUser, error)
	Exists(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, user *FrontendUser) error
	List(ctx context.Context, pager *upagination.Pagination, filter *FrontendUserListFilter) (users []*FrontendUser, total int, err error)
	Update(ctx context.Context, user *FrontendUser) error
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/frontend_user_interactor.go -package=mock . FrontendUserInteractor
type FrontendUserInteractor interface {
	Login(ctx context.Context, username string, password string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *FrontendUser, err error)
	Create(ctx context.Context, user *FrontendUser) error
	List(ctx context.Context, pager *upagination.Pagination, filter *FrontendUserListFilter) (users []*FrontendUser, total int, err error)
	Update(ctx context.Context, user *FrontendUser) error
	Find(ctx context.Context, id int) (*FrontendUser, error)
}

type FrontendUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"` // 用户名
	HashedPassword string `json:"-"`        // 密码哈希
	Nickname       string `json:"nickname"` // 昵称
	StoreID        int    `json:"store_id"` // 所属门店ID

	Store *Store `json:"store,omitempty"` // 所属门店
}

func (u *FrontendUser) GetOperatorID() int {
	return u.ID
}

func (u *FrontendUser) GetOperatorName() string {
	return u.Nickname
}

func (u *FrontendUser) GetOperatorType() OperatorType {
	return OperatorTypeFrontend
}

func (u *FrontendUser) GetOperatorStoreID() int {
	return u.StoreID
}

func (u *FrontendUser) SetPassword(password string) error {
	hashed, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	u.HashedPassword = hashed
	return nil
}

func (u *FrontendUser) CheckPassword(password string) error {
	if err := util.CheckPassword(password, u.HashedPassword); err != nil {
		return ErrMismatchedHashAndPassword
	}

	return nil
}

type (
	frontendUserKey struct{}
)

func NewFrontendUserContext(ctx context.Context, u *FrontendUser) context.Context {
	return context.WithValue(ctx, frontendUserKey{}, u)
}

func FromFrontendUserContext(ctx context.Context) *FrontendUser {
	if v, ok := ctx.Value(frontendUserKey{}).(*FrontendUser); ok {
		return v
	}
	return nil
}

type FrontendUserListFilter struct {
	StoreID int `json:"store_id"`
}
