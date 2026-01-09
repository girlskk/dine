package frontend

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.jiguang.dev/pos-dine/dine/api/frontend/docs"
	"gitlab.jiguang.dev/pos-dine/dine/buildinfo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin"
	uvalidator "gitlab.jiguang.dev/pos-dine/dine/pkg/validator"
	"go.uber.org/fx"
)

const ApiPrefixV1 = "/api/v1"

var middlewares = []string{
	"Recovery",
	"TimeLimiter",
	"PopulateRequestID",
	"PopulateLogger",
	"Observability",
	"Logger",
	"ErrorHandling",
	"Tenant",
}

type Params struct {
	fx.In

	AppConfig domain.AppConfig

	Middlewares []ugin.Middleware `group:"middlewares"`
	Handlers    []ugin.Handler    `group:"handlers"`
}

//	@title			前台 API
//	@version		1.0
//	@description	供收银机等客户端调用.

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and the JWT token.

// @BasePath	/api/v1
func New(p Params) (*gin.Engine, error) {
	e := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := uvalidator.RegisterDecimalValidator(v); err != nil {
			return nil, err
		}
	}

	if p.AppConfig.RunMode == domain.RunModeDev {
		gin.SetMode(gin.DebugMode)
		e.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	e.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"run_mode":      p.AppConfig.RunMode,
			"build_version": buildinfo.Version,
			"build_at":      buildinfo.BuildAt,
		})
	})

	r := e.Group(ApiPrefixV1)

	midMap := lo.KeyBy(p.Middlewares, func(m ugin.Middleware) string { return m.Name() })

	for _, name := range middlewares {
		mid, ok := midMap[name]
		if !ok {
			return nil, fmt.Errorf("middleware %s not found", name)
		}
		r.Use(mid.Middleware())
	}

	for _, h := range p.Handlers {
		h.Routes(r)
	}

	return e, nil
}
