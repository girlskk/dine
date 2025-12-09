package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.FrontendUserInteractor = (*FrontendUserInteractor)(nil)

type FrontendUserInteractor struct {
	AuthConfig domain.AuthConfig
	DataStore  domain.DataStore
}

func NewFrontendUserInteractor(authConfig domain.AuthConfig, dataStore domain.DataStore) *FrontendUserInteractor {
	return &FrontendUserInteractor{
		AuthConfig: authConfig,
		DataStore:  dataStore,
	}
}

type AuthToken struct {
	ID int
	jwt.RegisteredClaims
}

func (interactor *FrontendUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserInteractor.Login")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user, err := interactor.DataStore.GetFrontendUserRepository().FindByUsername(ctx, username)
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

func (interactor *FrontendUserInteractor) Logout(ctx context.Context) error { return nil }

func (interactor *FrontendUserInteractor) Authenticate(ctx context.Context, token string) (user *domain.FrontendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserInteractor.Authenticate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("FrontendUserInteractor.Authenticate")

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

	user, err = interactor.DataStore.GetFrontendUserRepository().Find(ctx, claims.ID)
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

func (interactor *FrontendUserInteractor) Create(ctx context.Context, user *domain.FrontendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := interactor.DataStore.GetFrontendUserRepository().Exists(ctx, user.Username)
		if err != nil {
			return err
		}
		if exists {
			return domain.ParamsError(domain.ErrUserExists)
		}
		return interactor.DataStore.GetFrontendUserRepository().Create(ctx, user)
	})
}

func (interactor *FrontendUserInteractor) Update(ctx context.Context, user *domain.FrontendUser) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if err = interactor.DataStore.GetFrontendUserRepository().Update(ctx, user); err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrUserNotFound)
		}

		if domain.IsConflict(err) {
			err = domain.ParamsError(domain.ErrUserExists)
		}

		err = fmt.Errorf("failed to update user: %w", err)
		return
	}

	return
}

func (interactor *FrontendUserInteractor) Find(ctx context.Context, id int) (user *domain.FrontendUser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserInteractor.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if user, err = interactor.DataStore.GetFrontendUserRepository().Find(ctx, id); err != nil {
		if domain.IsNotFound(err) {
			err = domain.ParamsError(domain.ErrUserNotFound)
		}
		return
	}

	return
}

func (interactor *FrontendUserInteractor) List(ctx context.Context, pager *upagination.Pagination, filter *domain.FrontendUserListFilter) (users []*domain.FrontendUser, total int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FrontendUserInteractor.List")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DataStore.GetFrontendUserRepository().List(ctx, pager, filter)
}
