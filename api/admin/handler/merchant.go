package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type MerchantHandler struct {
	MerchantInteractor domain.MerchantInteractor
}

func NewMerchantHandler(merchantInteractor domain.MerchantInteractor) *MerchantHandler {
	return &MerchantHandler{
		MerchantInteractor: merchantInteractor,
	}
}

func (h *MerchantHandler) Routes(r gin.IRouter) {
	r = r.Group("merchant/merchant")
	r.POST("/brand", h.CreateBrandMerchant())
	r.PUT("/brand/:id", h.UpdateBrandMerchant())
	r.DELETE("/:id", h.DeleteMerchant())
}

// CreateBrandMerchant 创建品牌商户
// @Summary 创建品牌商户
// @Description 创建品牌商户
// @Tags Merchant
// @Accept JSON
// @Produce JSON
// @Param data body types.CreateMerchantReq true "创建品牌商户请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /merchant/merchant/brand [post]
func (h *MerchantHandler) CreateBrandMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.CreateBrandMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var err error
		var req types.CreateMerchantReq
		if err = c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		if req.MerchantType != domain.MerchantTypeBrand {
			// todo
		}
		createBrandMerchant := &domain.CreateMerchantParams{
			MerchantCode:         req.MerchantCode,
			MerchantName:         req.MerchantName,
			MerchantShortName:    req.MerchantShortName,
			MerchantType:         req.MerchantType,
			BrandName:            req.BrandName,
			AdminPhoneNumber:     req.AdminPhoneNumber,
			PurchaseDuration:     req.PurchaseDuration,
			PurchaseDurationUnit: req.PurchaseDurationUnit,
			BusinessTypeID:       req.BusinessTypeID,
			MerchantLogo:         req.MerchantLogo,
			Description:          req.Description,
			Status:               req.Status,
			LoginAccount:         req.LoginAccount,
			LoginPassword:        req.LoginPassword,
		}
		if req.Address.CountryID != uuid.Nil {
			createBrandMerchant.Address = &domain.Address{
				CountryID:  req.Address.CountryID,
				ProvinceID: req.Address.ProvinceID,
				CityID:     req.Address.CityID,
				DistrictID: req.Address.DistrictID,
				Address:    req.Address.Address,
				Lng:        req.Address.Lng,
				Lat:        req.Address.Lat,
			}
		}
		err = h.MerchantInteractor.CreateMerchant(ctx, createBrandMerchant)
		if err != nil {
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.MerchantNameExists, err))
				return
			}
			err = fmt.Errorf("failed to create brand merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// UpdateBrandMerchant 更新品牌商户
// @Summary 更新品牌商户
// @Description 更新品牌商户
// @Tags Merchant
// @Accept JSON
// @Produce JSON
// @Param id path string true "商户ID"
// @Param data body types.UpdateMerchantReq true "更新品牌商户请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /merchant/merchant/brand/{id} [put]
func (h *MerchantHandler) UpdateBrandMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.UpdateBrandMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var err error

		// 从路径参数获取商户 ID
		merchantIDStr := c.Param("id")
		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.UpdateMerchantReq
		if err = c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		updateBrandMerchant := &domain.UpdateMerchantParams{
			ID:                merchantID,
			MerchantCode:      req.MerchantCode,
			MerchantName:      req.MerchantName,
			MerchantShortName: req.MerchantShortName,
			BrandName:         req.BrandName,
			AdminPhoneNumber:  req.AdminPhoneNumber,
			BusinessTypeID:    req.BusinessTypeID,
			MerchantLogo:      req.MerchantLogo,
			Description:       req.Description,
			Status:            req.Status,
			LoginAccount:      req.LoginAccount,
			LoginPassword:     req.LoginPassword,
		}
		if req.Address.CountryID != uuid.Nil {
			updateBrandMerchant.Address = &domain.Address{
				CountryID:  req.Address.CountryID,
				ProvinceID: req.Address.ProvinceID,
				CityID:     req.Address.CityID,
				DistrictID: req.Address.DistrictID,
				Address:    req.Address.Address,
				Lng:        req.Address.Lng,
				Lat:        req.Address.Lat,
			}
		}

		err = h.MerchantInteractor.UpdateMerchant(ctx, updateBrandMerchant)
		if err != nil {
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.MerchantNameExists, err))
				return
			}
			err = fmt.Errorf("failed to update brand merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}

}

// DeleteMerchant 删除商户
// @Summary 删除商户
// @Description 删除商户
// @Tags Merchant
// @Accept JSON
// @Produce JSON
// @Param id path string true "商户ID"
// @Success 200 {object} response.Response "删除成功"
// @Success 204 {object} response.Response "删除成功，无内容或资源不存在幂等处理"
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /merchant/merchant/{id} [delete]
func (h *MerchantHandler) DeleteMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.DeleteMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var err error

		// 从路径参数获取商户 ID
		merchantIDStr := c.Param("id")
		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		err = h.MerchantInteractor.DeleteMerchant(ctx, merchantID)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNoContent, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to delete merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// GetMerchant 获取商户信息
// @Summary 获取商户信息
// @Description 根据商户ID获取商户信息
// @Tags Merchant
// @Accept JSON
// @Produce JSON
// @Param id path string true "商户ID"
// @Success 200 {object} response.Response{data=types.MerchantInfoResp}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /merchant/merchant/{id} [get]
func (h *MerchantHandler) GetMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.DeleteMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var err error

		// 从路径参数获取商户 ID
		merchantIDStr := c.Param("id")
		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		domainMerchant, err := h.MerchantInteractor.GetMerchant(ctx, merchantID)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, &types.MerchantInfoResp{
			Merchant: domainMerchant,
		})
	}
}
