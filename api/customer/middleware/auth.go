package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/api/customer"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
)

type Auth struct {
	customerInteractor domain.CustomerInteractor
	skipper            middleware.SkipperFunc
}

func (u *Auth) Name() string {
	return "Auth"
}

func NewAuth(handlers []ugin.Handler, customerInteractor domain.CustomerInteractor) *Auth {
	var prefixes []string
	for _, h := range handlers {
		switch v := h.(type) {
		case interface{ NoAuths() []string }:
			for _, n := range v.NoAuths() {
				prefixes = append(prefixes, customer.ApiPrefixV1+n)
			}
		}
	}
	skipper := middleware.AllowPathPrefixSkipper(prefixes...)

	return &Auth{
		customerInteractor: customerInteractor,
		skipper:            skipper,
	}
}

func (u *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if middleware.SkipHandler(c, u.skipper) {
			c.Next()
			return
		}

		auths := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(auths) != 2 || !strings.EqualFold(auths[0], "Bearer") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := auths[1]
		ctx := c.Request.Context()

		u, err := u.customerInteractor.Authenticate(ctx, token)
		if err != nil {
			if errors.Is(err, domain.ErrTokenInvalid) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			err = fmt.Errorf("failed to authenticate customer: %w", err)
			c.Error(err)
			c.Abort()
			return
		}

		ctx = domain.NewCustomerContext(ctx, u)
		c.Request = c.Request.Clone(ctx)

		c.Next()
	}
}
