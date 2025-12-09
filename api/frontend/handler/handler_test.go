package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/alert"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/middleware"
	uvalidator "gitlab.jiguang.dev/pos-dine/dine/pkg/validator"
)

type HandlerTestSuite struct {
	suite.Suite
	r *gin.Engine
}

func (suite *HandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := uvalidator.RegisterDecimalValidator(v); err != nil {
			panic(err)
		}
	}
	alert := &alert.AlertNoop{}
	suite.r = gin.New()
	suite.r.Use(middleware.NewRecovery(alert).Middleware())
	suite.r.Use(middleware.NewErrorHandling(alert).Middleware())
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
