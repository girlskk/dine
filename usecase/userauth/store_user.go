package userauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreUserInteractor = (*StoreUserInteractor)(nil)

type StoreUserInteractor struct {
	AuthConfig domain.AuthConfig
	DataStore  domain.DataStore
}

func NewStoreUserInteractor(authConfig domain.AuthConfig, dataStore domain.DataStore) *StoreUserInteractor {
	return &StoreUserInteractor{
		AuthConfig: authConfig,
		DataStore:  dataStore,
	}
}

func (interactor *StoreUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.Login")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user, err := interactor.DataStore.StoreUserRepo().FindByUsername(ctx, username)
	if err != nil {
		return
	}
	if err = user.CheckPassword(password); err != nil {
		return
	}

	expAt = time.Now().Add(time.Duration(interactor.AuthConfig.Expire) * time.Second)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthToken{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
	})

	token, err = claims.SignedString([]byte(interactor.AuthConfig.Secret))
	if err != nil {
		err = fmt.Errorf("failed to sign token: %w", err)
		return
	}

	return
}

func (interactor *StoreUserInteractor) Logout(ctx context.Context) error { return nil }

func (interactor *StoreUserInteractor) Authenticate(ctx context.Context, token string) (user *domain.StoreUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.Authenticate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("StoreUserInteractor.Authenticate")

	claims := &AuthToken{}
	tokenInfo, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(interactor.AuthConfig.Secret), nil
	})
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenNotValidYet) && !errors.Is(err, jwt.ErrTokenExpired) {
			logger.Errorw("failed to parse token", "error", err)
		}
		err = domain.ErrTokenInvalid
		return
	}

	if !tokenInfo.Valid {
		err = domain.ErrTokenInvalid
		return
	}

	user, err = interactor.DataStore.StoreUserRepo().Find(ctx, claims.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.ErrTokenInvalid
			return
		}
		err = fmt.Errorf("failed to find user[%d]: %w", claims.ID, err)
		return
	}

	return
}
