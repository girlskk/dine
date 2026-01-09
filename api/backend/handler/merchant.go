package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type MerchantHandler struct {
	MerchantInteractor domain.MerchantInteractor
	StoreInteractor    domain.StoreInteractor
}

func NewMerchantHandler(merchantInteractor domain.MerchantInteractor, storeInteractor domain.StoreInteractor) *MerchantHandler {
	return &MerchantHandler{
		MerchantInteractor: merchantInteractor,
		StoreInteractor:    storeInteractor,
	}
}

func (h *MerchantHandler) Routes(r gin.IRouter) {
	r = r.Group("merchant")
	r.PUT("/brand", h.UpdateBrandMerchant())
	r.PUT("/store", h.UpdateStoreMerchant())
	r.GET("", h.GetMerchant())
	r.POST("/renewal", h.MerchantRenewal())
	r.PUT("/enable", h.Enable())
	r.PUT("/disable", h.Disable())
}

// UpdateBrandMerchant 更新品牌商户
//
//	@Summary		更新品牌商户
//	@Description	更新品牌商户
//	@Tags			商户管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"商户ID"
//	@Param			data	body	types.UpdateMerchantReq	true	"更新品牌商户请求"
//	@Success		200		"No Content"
//	@Router			/merchant/brand [put]
func (h *MerchantHandler) UpdateBrandMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.UpdateBrandMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var err error

		user := domain.FromBackendUserContext(ctx)

		var req types.UpdateMerchantReq
		if err = c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		updateBrandMerchant := &domain.UpdateMerchantParams{
			ID:                user.MerchantID,
			MerchantCode:      req.MerchantCode,
			MerchantName:      req.MerchantName,
			MerchantShortName: req.MerchantShortName,
			BrandName:         req.BrandName,
			AdminPhoneNumber:  req.AdminPhoneNumber,
			BusinessTypeCode:  req.BusinessTypeCode,
			MerchantLogo:      req.MerchantLogo,
			Description:       req.Description,
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

// UpdateStoreMerchant 更新门店商户
//
//	@Summary		更新门店商户
//	@Description	更新门店商户（商户 + 门店）
//	@Tags			商户管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string							true	"商户ID"
//	@Param			data	body	types.UpdateStoreMerchantReq	true	"更新门店商户请求"
//	@Success		200		"No Content"
//	@Router			/merchant/store [put]
func (h *MerchantHandler) UpdateStoreMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.UpdateStoreMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.UpdateStoreMerchantReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		address := &domain.Address{
			CountryID:  req.Merchant.Address.CountryID,
			ProvinceID: req.Merchant.Address.ProvinceID,
			CityID:     req.Merchant.Address.CityID,
			DistrictID: req.Merchant.Address.DistrictID,
			Address:    req.Merchant.Address.Address,
			Lng:        req.Merchant.Address.Lng,
			Lat:        req.Merchant.Address.Lat,
		}
		updateMerchant := &domain.UpdateMerchantParams{
			ID:                user.MerchantID,
			MerchantCode:      req.Merchant.MerchantCode,
			MerchantName:      req.Merchant.MerchantName,
			MerchantShortName: req.Merchant.MerchantShortName,
			BrandName:         req.Merchant.BrandName,
			AdminPhoneNumber:  req.Merchant.AdminPhoneNumber,
			BusinessTypeCode:  req.Merchant.BusinessTypeCode,
			MerchantLogo:      req.Merchant.MerchantLogo,
			Description:       req.Merchant.Description,
			Address:           address, // 门店商户的地址使用门店的地址
		}

		storeMerchant, err := h.StoreInteractor.GetStoreByMerchantID(ctx, user.MerchantID)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get store by merchant id: %w", err)
			c.Error(err)
			return
		}
		updateStore := &domain.UpdateStoreParams{
			ID:                      storeMerchant.ID,
			AdminPhoneNumber:        req.Store.AdminPhoneNumber,
			StoreName:               req.Store.StoreName,
			StoreShortName:          req.Store.StoreShortName,
			StoreCode:               req.Store.StoreCode,
			Status:                  req.Store.Status,
			BusinessModel:           domain.BusinessModelDirect,
			BusinessTypeCode:        req.Merchant.BusinessTypeCode,
			LocationNumber:          req.Store.LocationNumber,
			ContactName:             req.Store.ContactName,
			ContactPhone:            req.Store.ContactPhone,
			UnifiedSocialCreditCode: req.Store.UnifiedSocialCreditCode,
			StoreLogo:               req.Store.StoreLogo,
			BusinessLicenseURL:      req.Store.BusinessLicenseURL,
			StorefrontURL:           req.Store.StorefrontURL,
			CashierDeskURL:          req.Store.CashierDeskURL,
			DiningEnvironmentURL:    req.Store.DiningEnvironmentURL,
			FoodOperationLicenseURL: req.Store.FoodOperationLicenseURL,
			BusinessHours:           req.Store.BusinessHours,
			DiningPeriods:           req.Store.DiningPeriods,
			ShiftTimes:              req.Store.ShiftTimes,
			Address:                 address,
		}

		if err := h.MerchantInteractor.UpdateMerchantAndStore(ctx, updateMerchant, updateStore); err != nil {
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.MerchantNameExists, err))
				return
			}
			err = fmt.Errorf("failed to update store merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// GetMerchant 获取商户信息
//
//	@Summary		获取商户信息
//	@Description	根据商户ID获取商户信息
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"商户ID"
//	@Success		200	{object}	response.Response{data=types.MerchantInfoResp}
//	@Router			/merchant [get]
func (h *MerchantHandler) GetMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.GetMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)
		var err error

		user := domain.FromBackendUserContext(ctx)

		domainMerchant, err := h.MerchantInteractor.GetMerchant(ctx, user.MerchantID)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get merchant: %w", err)
			c.Error(err)
			return
		}
		merchantInfo := &types.MerchantInfoResp{
			Merchant: domainMerchant,
		}
		if domainMerchant.MerchantType == domain.MerchantTypeStore {
			store, err := h.StoreInteractor.GetStoreByMerchantID(ctx, domainMerchant.ID)
			if err != nil {
				if domain.IsNotFound(err) {
					c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
					return
				}
				err = fmt.Errorf("failed to get store by merchant id: %w", err)
				c.Error(err)
				return
			}
			merchantInfo.Store = store
		}
		response.Ok(c, merchantInfo)
	}
}

// MerchantRenewal 商户续期
//
//	@Summary		商户续期
//	@Description	商户续费续期
//	@Tags			商户管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.MerchantRenewalReq	true	"商户续期请求"
//	@Success		200		"No Content"
//	@Router			/merchant/renewal [post]
func (h *MerchantHandler) MerchantRenewal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.MerchantRenewal")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.MerchantRenewalReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromBackendUserContext(ctx)

		merchantRenewal := &domain.MerchantRenewal{
			MerchantID:           user.MerchantID,
			PurchaseDuration:     req.PurchaseDuration,
			PurchaseDurationUnit: req.PurchaseDurationUnit,
			OperatorAccount:      user.ID.String(),
			OperatorName:         user.Username,
		}

		if err := h.MerchantInteractor.MerchantRenewal(ctx, merchantRenewal); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to renew merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Enable 启用商户
//
//	@Summary		启用商户
//	@Description	将商户状态置为激活
//	@Tags			商户管理
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	"No Content"
//	@Router			/merchant/enable [put]
func (h *MerchantHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromBackendUserContext(ctx)

		updateParams := &domain.Merchant{ID: user.MerchantID, Status: domain.MerchantStatusActive}

		if err := h.MerchantInteractor.MerchantSimpleUpdate(ctx, domain.MerchantSimpleUpdateTypeStatus, updateParams); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to simple update merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用商户
//
//	@Summary		禁用商户
//	@Description	将商户状态置为禁用
//	@Tags			商户管理
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	"No Content"
//	@Router			/merchant/disable [put]
func (h *MerchantHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromBackendUserContext(ctx)

		updateParams := &domain.Merchant{ID: user.MerchantID, Status: domain.MerchantStatusDisabled}

		if err := h.MerchantInteractor.MerchantSimpleUpdate(ctx, domain.MerchantSimpleUpdateTypeStatus, updateParams); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to simple update merchant: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}
