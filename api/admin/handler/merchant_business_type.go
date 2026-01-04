package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type MerchantBusinessTypeHandler struct {
	BusinessTypeInteractor domain.MerchantBusinessTypeInteractor
}

func NewMerchantBusinessTypeHandler(businessTypeInteractor domain.MerchantBusinessTypeInteractor) *MerchantBusinessTypeHandler {
	return &MerchantBusinessTypeHandler{
		BusinessTypeInteractor: businessTypeInteractor,
	}
}

func (h *MerchantBusinessTypeHandler) Routes(r gin.IRouter) {
	r = r.Group("merchant/business_type")
	r.GET("/list", h.GetAll())
}

// GetAll 业态列表
//
//	@Summary		业态列表
//	@Description	业态列表
//	@Tags			商户管理-业态列表
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=types.MerchantBusinessTypeListResp}
//	@Router			/merchant/business_type/list [get]
func (h *MerchantBusinessTypeHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantBusinessTypeHandler.GetAll")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		businessTypes, err := h.BusinessTypeInteractor.GetAll(ctx)
		if err != nil {
			err = fmt.Errorf("failed to list businessTypes: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, &types.MerchantBusinessTypeListResp{
			BusinessTypes: businessTypes,
		})
	}
}
