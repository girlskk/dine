package domain

import (
	"context"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"time"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/backend_user_interactor.go -package=mock . BackendUserInteractor
type BackendUserInteractor interface {
	Login(ctx context.Context, username string, password string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *BackendUser, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/backend_user_repository.go -package=mock . BackendUserRepository
type BackendUserRepository interface {
	Create(ctx context.Context, user *BackendUser) error
	FindByUsername(ctx context.Context, username string) (*BackendUser, error)
	FindByStoreID(ctx context.Context, storeID int) (*BackendUser, error)
	Exists(ctx context.Context, username string) (bool, error)
	Find(ctx context.Context, id int) (*BackendUser, error)
	Update(ctx context.Context, user *BackendUser) error
}

type BackendUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"`
	Nickname       string `json:"nickname"`
	StoreID        int    `json:"store_id"`

	Store *Store `json:"store,omitempty"` // 所属门店
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
