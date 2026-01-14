package handler

import (
	"errors"
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
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

type MerchantHandler struct {
	MerchantInteractor domain.MerchantInteractor
	StoreInteractor    domain.StoreInteractor
}

func NewMerchantHandler(
	merchantInteractor domain.MerchantInteractor,
	storeInteractor domain.StoreInteractor,
) *MerchantHandler {
	return &MerchantHandler{
		MerchantInteractor: merchantInteractor,
		StoreInteractor:    storeInteractor,
	}
}

func (h *MerchantHandler) Routes(r gin.IRouter) {
	r = r.Group("merchant/merchant")
	r.POST("/brand", h.CreateBrandMerchant())
	r.POST("/store", h.CreateStoreMerchant())
	r.PUT("/brand/:id", h.UpdateBrandMerchant())
	r.PUT("/store/:id", h.UpdateStoreMerchant())
	r.DELETE("/:id", h.DeleteMerchant())
	r.GET("/:id", h.GetMerchant())
	r.GET("/list", h.GetMerchants())
	r.POST("/renewal", h.MerchantRenewal())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
	r.GET("/count", h.CountMerchant())
}

// CreateBrandMerchant 创建品牌商户
//
//	@Summary		创建品牌商户
//	@Description	创建品牌商户
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.CreateMerchantReq	true	"创建品牌商户请求"
//	@Success		200		"No Content"
//	@Router			/merchant/merchant/brand [post]
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

		hashPwd, err := util.HashPassword(req.LoginPassword)
		if err != nil {
			c.Error(errorx.New(http.StatusInternalServerError, errcode.InternalError, err))
			return
		}
		createBrandMerchant := &domain.CreateMerchantParams{
			MerchantCode:         req.MerchantCode,
			MerchantName:         req.MerchantName,
			MerchantShortName:    req.MerchantShortName,
			MerchantType:         domain.MerchantTypeBrand,
			BrandName:            req.BrandName,
			AdminPhoneNumber:     req.AdminPhoneNumber,
			PurchaseDuration:     req.PurchaseDuration,
			PurchaseDurationUnit: req.PurchaseDurationUnit,
			BusinessTypeCode:     req.BusinessTypeCode,
			MerchantLogo:         req.MerchantLogo,
			Description:          req.Description,
			LoginAccount:         req.LoginAccount,
			LoginPassword:        hashPwd,
		}
		createBrandMerchant.Address = &domain.Address{
			Country:  req.Address.Country,
			Province: req.Address.Province,
			Address:  req.Address.Address,
			Lng:      req.Address.Lng,
			Lat:      req.Address.Lat,
		}
		user := domain.FromAdminUserContext(ctx)
		err = h.MerchantInteractor.CreateMerchant(ctx, createBrandMerchant, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// CreateStoreMerchant 创建门店商户
//
//	@Summary		创建门店商户
//	@Description	创建门店商户（商户 + 门店）
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.CreateStoreMerchantReq	true	"创建门店商户请求"
//	@Success		200		"No Content"
//	@Router			/merchant/merchant/store [post]
func (h *MerchantHandler) CreateStoreMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.CreateStoreMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CreateStoreMerchantReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		address := &domain.Address{
			Country:  req.Merchant.Address.Country,
			Province: req.Merchant.Address.Province,
			Address:  req.Merchant.Address.Address,
			Lng:      req.Merchant.Address.Lng,
			Lat:      req.Merchant.Address.Lat,
		}
		hashPwd, err := util.HashPassword(req.Merchant.LoginPassword)
		if err != nil {
			c.Error(errorx.New(http.StatusInternalServerError, errcode.InternalError, err))
			return
		}
		createMerchant := &domain.CreateMerchantParams{
			MerchantCode:         req.Merchant.MerchantCode,
			MerchantName:         req.Merchant.MerchantName,
			MerchantShortName:    req.Merchant.MerchantShortName,
			MerchantType:         domain.MerchantTypeStore,
			BrandName:            req.Merchant.BrandName,
			AdminPhoneNumber:     req.Merchant.AdminPhoneNumber,
			PurchaseDuration:     req.Merchant.PurchaseDuration,
			PurchaseDurationUnit: req.Merchant.PurchaseDurationUnit,
			BusinessTypeCode:     req.Merchant.BusinessTypeCode,
			MerchantLogo:         req.Merchant.MerchantLogo,
			Description:          req.Merchant.Description,
			LoginAccount:         req.Merchant.LoginAccount,
			LoginPassword:        hashPwd,
			Address:              address,
		}

		createStore := &domain.CreateStoreParams{
			AdminPhoneNumber:        req.Store.AdminPhoneNumber,
			StoreName:               req.Store.StoreName,
			StoreShortName:          req.Store.StoreShortName,
			StoreCode:               req.Store.StoreCode,
			Status:                  req.Store.Status,
			BusinessModel:           req.Store.BusinessModel,
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
			LoginAccount:            req.Merchant.LoginAccount,
			LoginPassword:           req.Merchant.LoginPassword,
			BusinessHours:           req.Store.BusinessHours,
			DiningPeriods:           req.Store.DiningPeriods,
			ShiftTimes:              req.Store.ShiftTimes,
			Address:                 address,
		}
		user := domain.FromAdminUserContext(ctx)
		err = h.MerchantInteractor.CreateMerchantAndStore(ctx, createMerchant, createStore, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// UpdateBrandMerchant 更新品牌商户
//
//	@Summary		更新品牌商户
//	@Description	更新品牌商户
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"商户ID"
//	@Param			data	body	types.UpdateMerchantReq	true	"更新品牌商户请求"
//	@Success		200		"No Content"
//	@Router			/merchant/merchant/brand/{id} [put]
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
			BusinessTypeCode:  req.BusinessTypeCode,
			MerchantLogo:      req.MerchantLogo,
			Description:       req.Description,
		}
		updateBrandMerchant.Address = &domain.Address{
			Country:  req.Address.Country,
			Province: req.Address.Province,
			Address:  req.Address.Address,
			Lng:      req.Address.Lng,
			Lat:      req.Address.Lat,
		}
		user := domain.FromAdminUserContext(ctx)
		err = h.MerchantInteractor.UpdateMerchant(ctx, updateBrandMerchant, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}

}

// UpdateStoreMerchant 更新门店商户
//
//	@Summary		更新门店商户
//	@Description	更新门店商户（商户 + 门店）
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string							true	"商户ID"
//	@Param			data	body	types.UpdateStoreMerchantReq	true	"更新门店商户请求"
//	@Success		200		"No Content"
//	@Router			/merchant/merchant/store/{id} [put]
func (h *MerchantHandler) UpdateStoreMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.UpdateStoreMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		merchantIDStr := c.Param("id")
		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.UpdateStoreMerchantReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		address := &domain.Address{
			Country:  req.Merchant.Address.Country,
			Province: req.Merchant.Address.Province,
			Address:  req.Merchant.Address.Address,
			Lng:      req.Merchant.Address.Lng,
			Lat:      req.Merchant.Address.Lat,
		}
		updateMerchant := &domain.UpdateMerchantParams{
			ID:                merchantID,
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

		storeMerchant, err := h.StoreInteractor.GetStoreByMerchantID(ctx, merchantID)
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
		user := domain.FromAdminUserContext(ctx)
		err = h.MerchantInteractor.UpdateMerchantAndStore(ctx, updateMerchant, updateStore, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// DeleteMerchant 删除商户
//
//	@Summary		删除商户
//	@Description	删除商户
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"商户ID"
//	@Success		200	"No Content"
//	@Router			/merchant/merchant/{id} [delete]
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

		user := domain.FromAdminUserContext(ctx)
		err = h.MerchantInteractor.DeleteMerchant(ctx, merchantID, user)
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
//
//	@Summary		获取商户信息
//	@Description	根据商户ID获取商户信息
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"商户ID"
//	@Success		200	{object}	response.Response{data=types.MerchantInfoResp}
//	@Router			/merchant/merchant/{id} [get]
func (h *MerchantHandler) GetMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.GetMerchant")
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

		user := domain.FromAdminUserContext(ctx)
		domainMerchant, err := h.MerchantInteractor.GetMerchant(ctx, merchantID, user)
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

// GetMerchants 商户列表
//
//	@Summary		商户列表
//	@Description	分页查询商户列表
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Produce		json
//	@Param			data	query		types.MerchantListReq	true	"商户列表查询参数"
//	@Success		200		{object}	response.Response{data=types.MerchantListResp}
//	@Router			/merchant/merchant/list [get]
func (h *MerchantHandler) GetMerchants() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.GetMerchants")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.MerchantListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		pager := req.RequestPagination.ToPagination()
		filter := &domain.MerchantListFilter{
			MerchantName:     req.MerchantName,
			AdminPhoneNumber: req.AdminPhoneNumber,
			MerchantType:     req.MerchantType,
			Status:           req.Status,
			Province:         req.Province,
		}

		if req.CreatedAtGte != "" || req.CreatedAtLte != "" {
			var err error
			startTime, endTime := util.GetShortcutDate("custom", req.CreatedAtGte, req.CreatedAtLte)
			filter.CreatedAtGte, err = util.ParseDateToPtr(startTime)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.TimeFormatInvalid, fmt.Errorf("invalid CreatedAtGte: %w", err)))
				return
			}
			filter.CreatedAtLte, err = util.ParseDateToPtr(endTime)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.TimeFormatInvalid, fmt.Errorf("invalid CreatedAtLte: %w", err)))
				return
			}
		}

		domainMerchants, total, err := h.MerchantInteractor.GetMerchants(ctx, pager, filter)
		if err != nil {
			err = fmt.Errorf("failed to list merchants: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, &types.MerchantListResp{Merchants: domainMerchants, Total: total})
	}
}

// MerchantRenewal 商户续期
//
//	@Summary		商户续期
//	@Description	商户续费续期
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.MerchantRenewalReq	true	"商户续期请求"
//	@Success		200		"No Content"
//	@Router			/merchant/merchant/renewal [post]
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
		user := domain.FromAdminUserContext(ctx)

		merchantRenewal := &domain.MerchantRenewal{
			MerchantID:           req.MerchantID,
			PurchaseDuration:     req.PurchaseDuration,
			PurchaseDurationUnit: req.PurchaseDurationUnit,
			OperatorAccount:      user.ID.String(),
			OperatorName:         user.Username,
		}
		err := h.MerchantInteractor.MerchantRenewal(ctx, merchantRenewal, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Enable 启用商户
//
//	@Summary		启用商户
//	@Description	将商户状态置为激活
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"商户ID"
//	@Success		200	"No Content"
//	@Router			/merchant/merchant/{id}/enable [put]
func (h *MerchantHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		updateParams := &domain.Merchant{ID: id, Status: domain.MerchantStatusActive}
		user := domain.FromAdminUserContext(ctx)
		err = h.MerchantInteractor.MerchantSimpleUpdate(ctx, domain.MerchantSimpleUpdateTypeStatus, updateParams, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用商户
//
//	@Summary		禁用商户
//	@Description	将商户状态置为禁用
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"商户ID"
//	@Success		200	"No Content"
//	@Router			/merchant/merchant/{id}/disable [put]
func (h *MerchantHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		updateParams := &domain.Merchant{ID: id, Status: domain.MerchantStatusDisabled}
		user := domain.FromAdminUserContext(ctx)
		err = h.MerchantInteractor.MerchantSimpleUpdate(ctx, domain.MerchantSimpleUpdateTypeStatus, updateParams, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// CountMerchant 商户数量统计
//
//	@Summary		商户数量统计
//	@Description	获取商户数量统计
//	@Tags			商户管理-商户
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=types.MerchantCount}
//	@Router			/merchant/merchant/count [get]
func (h *MerchantHandler) CountMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("MerchantHandler.CountMerchant")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		merchantCount, err := h.MerchantInteractor.CountMerchant(ctx)
		if err != nil {
			err = fmt.Errorf("failed to count merchants: %w", err)
			c.Error(err)
			return
		}
		resp := &types.MerchantCount{}
		if merchantCount != nil {
			resp.MerchantTypeBrand = merchantCount.MerchantTypeBrand
			resp.MerchantTypeStore = merchantCount.MerchantTypeStore
			resp.Expired = merchantCount.Expired
		}
		response.Ok(c, resp)
	}
}

func (h *MerchantHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserExists):
		return errorx.New(http.StatusConflict, errcode.UserNameExists, err)
	case errors.Is(err, domain.ErrMerchantNameExists):
		return errorx.New(http.StatusConflict, errcode.MerchantNameExists, err)
	case errors.Is(err, domain.ErrStoreNameExists):
		return errorx.New(http.StatusConflict, errcode.StoreNameExists, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("merchant handler error: %w", err)
	}
}
