package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_user_interactor.go -package=mock . StoreUserInteractor
type StoreUserInteractor interface {
	Login(ctx context.Context, username string, password string) (token string, expAt time.Time, err error)
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, token string) (user *StoreUser, err error)
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/store_user_repository.go -package=mock . StoreUserRepository
type StoreUserRepository interface {
	Create(ctx context.Context, user *StoreUser) error
	FindByUsername(ctx context.Context, username string) (*StoreUser, error)
	Exists(ctx context.Context, username string) (bool, error)
	Find(ctx context.Context, id uuid.UUID) (*StoreUser, error)
	Update(ctx context.Context, user *StoreUser) error
}

type StoreUser struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	Nickname       string    `json:"nickname"`
	MerchantID     uuid.UUID `json:"merchant_id"`
	StoreID        uuid.UUID `json:"store_id"` // 门店ID

	// 关联数据
	Store *Store `json:"store,omitempty"` // 所属门店
}

// 实现 User 接口
// GetMerchantID 实现 User 接口
func (u *StoreUser) GetMerchantID() uuid.UUID {
	return u.MerchantID
}

// GetStoreID 实现 User 接口（品牌商用户的 StoreID 为 uuid.Nil）
func (u *StoreUser) GetStoreID() uuid.UUID {
	return u.StoreID
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
