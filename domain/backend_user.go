package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
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
	Exists(ctx context.Context, username string) (bool, error)
	Find(ctx context.Context, id uuid.UUID) (*BackendUser, error)
	Update(ctx context.Context, user *BackendUser) error
}

type BackendUser struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	Nickname       string    `json:"nickname"`
	MerchantID     uuid.UUID `json:"merchant_id"` // 品牌商ID

	// 关联数据
	Merchant *Merchant `json:"merchant,omitempty"` // 所属品牌商
}

// 实现 User 接口
// GetMerchantID 实现 User 接口
func (u *BackendUser) GetMerchantID() uuid.UUID {
	return u.MerchantID
}

// GetStoreID 实现 User 接口（品牌商用户的 StoreID 为 uuid.Nil）
func (u *BackendUser) GetStoreID() uuid.UUID {
	return uuid.Nil
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
