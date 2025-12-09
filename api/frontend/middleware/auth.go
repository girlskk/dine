package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
)

type Auth struct {
	userInteractor domain.FrontendUserInteractor
	skipper        middleware.SkipperFunc
}

func (u *Auth) Name() string {
	return "Auth"
}

func NewAuth(handlers []ugin.Handler, userInteractor domain.FrontendUserInteractor) *Auth {
	var prefixes []string
	for _, h := range handlers {
		switch v := h.(type) {
		case interface{ NoAuths() []string }:
			for _, n := range v.NoAuths() {
				prefixes = append(prefixes, frontend.ApiPrefixV1+n)
			}
		}
	}
	skipper := middleware.AllowPathPrefixSkipper(prefixes...)

	return &Auth{
		userInteractor: userInteractor,
		skipper:        skipper,
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

		u, err := u.userInteractor.Authenticate(ctx, token)
		if err != nil {
			if errors.Is(err, domain.ErrTokenInvalid) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			err = fmt.Errorf("failed to authenticate user: %w", err)
			c.Error(err)
			c.Abort()
			return
		}

		ctx = domain.NewFrontendUserContext(ctx, u)
		ctx = domain.NewStoreContext(ctx, u.Store)
		c.Request = c.Request.Clone(ctx)

		c.Next()
	}
}
