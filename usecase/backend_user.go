package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.BackendUserInteractor = (*BackendUserInteractor)(nil)

type BackendUserInteractor struct {
	AuthConfig domain.AuthConfig
	DataStore  domain.DataStore
}

func NewBackendUserInteractor(authConfig domain.AuthConfig, dataStore domain.DataStore) *BackendUserInteractor {
	return &BackendUserInteractor{
		AuthConfig: authConfig,
		DataStore:  dataStore,
	}
}

func (interactor *BackendUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserInteractor.Authentication")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user, err := interactor.DataStore.BackendUserRepo().FindByUsername(ctx, username)
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

func (interactor *BackendUserInteractor) Logout(ctx context.Context) error { return nil }

func (interactor *BackendUserInteractor) Authenticate(ctx context.Context, token string) (user *domain.BackendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BackendUserInteractor.Authentication")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("BackendUserInteractor.Authenticate")

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

	user, err = interactor.DataStore.BackendUserRepo().Find(ctx, claims.ID)
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
