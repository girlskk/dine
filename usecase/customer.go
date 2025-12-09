package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/silenceper/wechat/v2/miniprogram"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CustomerInteractor = (*CustomerInteractor)(nil)

type CustomerInteractor struct {
	AuthConfig  domain.AuthConfig
	DataStore   domain.DataStore
	MiniProgram *miniprogram.MiniProgram
}

func NewCustomerInteractor(authConfig domain.AuthConfig, dataStore domain.DataStore, miniProgram *miniprogram.MiniProgram) *CustomerInteractor {
	return &CustomerInteractor{
		AuthConfig:  authConfig,
		DataStore:   dataStore,
		MiniProgram: miniProgram,
	}
}

func (c *CustomerInteractor) WXLogin(ctx context.Context, code string) (token string, expAt time.Time, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerInteractor.WXLogin")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	//logger := logging.FromContext(ctx).Named("CustomerInteractor.WXLogin")

	//phoneInfo, err := c.MiniProgram.GetBusiness().GetPhoneNumberWithContext(ctx, &business.GetPhoneNumberRequest{
	//	Code: code,
	//})
	if err != nil {
		err = fmt.Errorf("failed to get phone number: %w", err)
		return
	}
	//logger.Debug("phoneInfo", phoneInfo)

	id, err := c.DataStore.CustomerRepo().FindOrCreate(ctx, &domain.Customer{
		Phone:    "13800138000",
		Nickname: "快来起个名字吧",
		Gender:   domain.GenderUnknown,
	})
	if err != nil {
		err = fmt.Errorf("failed to find or create customer: %w", err)
		return
	}

	expAt = time.Now().Add(time.Duration(c.AuthConfig.Expire) * time.Second)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthToken{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
	})

	token, err = claims.SignedString([]byte(c.AuthConfig.Secret))
	if err != nil {
		err = fmt.Errorf("failed to sign token: %w", err)
		return
	}

	return
}

func (c *CustomerInteractor) Logout(ctx context.Context) error { return nil }

func (c *CustomerInteractor) Authenticate(ctx context.Context, token string) (user *domain.Customer, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerInteractor.Authenticate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("CustomerInteractor.Authenticate")

	claims := &AuthToken{}
	tokenInfo, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(c.AuthConfig.Secret), nil
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

	user, err = c.DataStore.CustomerRepo().Find(ctx, claims.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ErrTokenInvalid
			return
		}
		err = fmt.Errorf("failed to find customer: %w", err)
		return
	}

	return
}
