package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type BusinessConfigHandler struct {
	BusinessConfigInteractor domain.BusinessConfigInteractor
}

func NewBusinessConfigHandler(BusinessConfigInteractor domain.BusinessConfigInteractor) *BusinessConfigHandler {
	return &BusinessConfigHandler{
		BusinessConfigInteractor: BusinessConfigInteractor,
	}
}

func (h *BusinessConfigHandler) Routes(r gin.IRouter) {
	r = r.Group("business/config")
	r.GET("", h.List())
}

func (h *BusinessConfigHandler) NoAuths() []string {
	return []string{}
}

// List
//
//	@Tags		经营管理
//	@Security	BearerAuth
//	@Summary	经营设置列表
//	@Param		name	query		string							false	"设置名称（模糊匹配）"
//	@Success	200		{object}	domain.BusinessConfigSearchRes	"成功"
//	@Router		/business/config [get]
func (h *BusinessConfigHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("BusinessConfigHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.BusinessConfigListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		params := domain.BusinessConfigSearchParams{
			MerchantID: user.MerchantID,
			Name:       req.Name,
		}
		res, err := h.BusinessConfigInteractor.ListBySearch(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list businessConfigs: %w", err)
				c.Error(err)
			}
			return
		}
		response.Ok(c, res)
	}
}
