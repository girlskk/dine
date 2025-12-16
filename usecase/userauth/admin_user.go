package userauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.AdminUserInteractor = (*AdminUserInteractor)(nil)

type AdminUserInteractor struct {
	AuthConfig domain.AuthConfig
	DS         domain.DataStore
}

func NewAdminUserInteractor(authConfig domain.AuthConfig, ds domain.DataStore) *AdminUserInteractor {
	return &AdminUserInteractor{
		AuthConfig: authConfig,
		DS:         ds,
	}
}

func (interactor *AdminUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.Login")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user, err := interactor.DS.AdminUserRepo().FindByUsername(ctx, username)
	if err != nil {
		err = fmt.Errorf("failed to find user by username: %w", err)
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

func (interactor *AdminUserInteractor) Logout(ctx context.Context) error { return nil }

func (interactor *AdminUserInteractor) Authenticate(ctx context.Context, token string) (user *domain.AdminUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.Authenticate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("AdminUserInteractor.Authenticate")

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

	user, err = interactor.DS.AdminUserRepo().Find(ctx, claims.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ErrTokenInvalid
			return
		}
		err = fmt.Errorf("failed to find user[%d]: %w", claims.ID, err)
		return
	}

	return
}
