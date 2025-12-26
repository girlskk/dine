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
)

// StoreHandler handles store related routes.
type StoreHandler struct {
	StoreInteractor domain.StoreInteractor
}

func NewStoreHandler(storeInteractor domain.StoreInteractor) *StoreHandler {
	return &StoreHandler{StoreInteractor: storeInteractor}
}

func (h *StoreHandler) Routes(r gin.IRouter) {
	r = r.Group("merchant/store")
	r.POST("", h.CreateStore())
	r.PUT("/:id", h.UpdateStore())
	r.DELETE("/:id", h.DeleteStore())
	r.GET("/:id", h.GetStore())
	r.GET("/list", h.GetStores())
	r.PATCH("/:id", h.StoreSimpleUpdate())
}

// CreateStore 创建门店
//
//	@Summary		创建门店
//	@Description	创建单个门店
//	@Tags			商户管理-门店
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.CreateStoreReq	true	"创建门店请求"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/merchant/store [post]
func (h *StoreHandler) CreateStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.CreateStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CreateStoreReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		domainStore := &domain.CreateStoreParams{
			MerchantID:              req.MerchantID,
			AdminPhoneNumber:        req.AdminPhoneNumber,
			StoreName:               req.StoreName,
			StoreShortName:          req.StoreShortName,
			StoreCode:               req.StoreCode,
			Status:                  req.Status,
			BusinessModel:           req.BusinessModel,
			BusinessTypeID:          req.BusinessTypeID,
			LocationNumber:          req.LocationNumber,
			ContactName:             req.ContactName,
			ContactPhone:            req.ContactPhone,
			UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
			StoreLogo:               req.StoreLogo,
			BusinessLicenseURL:      req.BusinessLicenseURL,
			StorefrontURL:           req.StorefrontURL,
			CashierDeskURL:          req.CashierDeskURL,
			DiningEnvironmentURL:    req.DiningEnvironmentURL,
			FoodOperationLicenseURL: req.FoodOperationLicenseURL,
			LoginAccount:            req.LoginAccount,
			LoginPassword:           req.LoginPassword,
			BusinessHours:           req.BusinessHours,
			DiningPeriods:           req.DiningPeriods,
			ShiftTimes:              req.ShiftTimes,
		}
		if req.Address.CountryID != uuid.Nil {
			domainStore.Address = &domain.Address{
				CountryID:  req.Address.CountryID,
				ProvinceID: req.Address.ProvinceID,
				CityID:     req.Address.CityID,
				DistrictID: req.Address.DistrictID,
				Address:    req.Address.Address,
				Lng:        req.Address.Lng,
				Lat:        req.Address.Lat,
			}
		}

		if err := h.StoreInteractor.CreateStore(ctx, domainStore); err != nil {
			c.Error(h.checkEditErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// UpdateStore 更新门店
//
//	@Summary		更新门店
//	@Description	更新单个门店
//	@Tags			商户管理-门店
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"门店ID"
//	@Param			data	body	types.UpdateStoreReq	true	"更新门店请求"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/merchant/store/{id} [put]
func (h *StoreHandler) UpdateStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.UpdateStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeIDStr := c.Param("id")
		storeID, err := uuid.Parse(storeIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.UpdateStoreReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		domainStore := &domain.UpdateStoreParams{
			ID:                      storeID,
			AdminPhoneNumber:        req.AdminPhoneNumber,
			StoreName:               req.StoreName,
			StoreShortName:          req.StoreShortName,
			StoreCode:               req.StoreCode,
			Status:                  req.Status,
			BusinessModel:           req.BusinessModel,
			BusinessTypeID:          req.BusinessTypeID,
			LocationNumber:          req.LocationNumber,
			ContactName:             req.ContactName,
			ContactPhone:            req.ContactPhone,
			UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
			StoreLogo:               req.StoreLogo,
			BusinessLicenseURL:      req.BusinessLicenseURL,
			StorefrontURL:           req.StorefrontURL,
			CashierDeskURL:          req.CashierDeskURL,
			DiningEnvironmentURL:    req.DiningEnvironmentURL,
			FoodOperationLicenseURL: req.FoodOperationLicenseURL,
			LoginPassword:           req.LoginPassword,
			BusinessHours:           req.BusinessHours,
			DiningPeriods:           req.DiningPeriods,
			ShiftTimes:              req.ShiftTimes,
		}
		if req.Address.CountryID != uuid.Nil {
			domainStore.Address = &domain.Address{
				CountryID:  req.Address.CountryID,
				ProvinceID: req.Address.ProvinceID,
				CityID:     req.Address.CityID,
				DistrictID: req.Address.DistrictID,
				Address:    req.Address.Address,
				Lng:        req.Address.Lng,
				Lat:        req.Address.Lat,
			}
		}

		if err := h.StoreInteractor.UpdateStore(ctx, domainStore); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(h.checkEditErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// DeleteStore 删除门店
//
//	@Summary		删除门店
//	@Description	删除单个门店
//	@Tags			商户管理-门店
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"门店ID"
//	@Success		200	"No Content"
//	@Success		204	"No Content"
//	@Failure		400	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/merchant/store/{id} [delete]
func (h *StoreHandler) DeleteStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.DeleteStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeIDStr := c.Param("id")
		storeID, err := uuid.Parse(storeIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.StoreInteractor.DeleteStore(ctx, storeID); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNoContent, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to delete store: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}

// GetStore 获取门店
//
//	@Summary		获取门店
//	@Description	根据门店ID获取门店信息
//	@Tags			商户管理-门店
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"门店ID"
//	@Success		200	{object}	response.Response{data=domain.Store}
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/merchant/store/{id} [get]
func (h *StoreHandler) GetStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.GetStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeIDStr := c.Param("id")
		storeID, err := uuid.Parse(storeIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		domainStore, err := h.StoreInteractor.GetStore(ctx, storeID)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to get store: %w", err))
			return
		}

		response.Ok(c, domainStore)
	}
}

// GetStores 门店列表
//
//	@Summary		门店列表
//	@Description	分页查询门店列表
//	@Tags			商户管理-门店
//	@Security		BearerAuth
//	@Produce		json
//	@Param			data	query		types.StoreListReq	true	"门店列表查询参数"
//	@Success		200		{object}	response.Response{data=types.StoreListResp}
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/merchant/store/list [get]
func (h *StoreHandler) GetStores() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.GetStores")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		pager := req.RequestPagination.ToPagination()
		filter := &domain.StoreListFilter{
			StoreName:        req.StoreName,
			MerchantID:       req.MerchantID,
			BusinessTypeID:   req.BusinessTypeID,
			AdminPhoneNumber: req.AdminPhoneNumber,
			Status:           req.Status,
			BusinessModel:    req.BusinessModel,
			CreatedAtGte:     &req.CreatedAtGte,
			CreatedAtLte:     &req.CreatedAtLte,
			ProvinceID:       req.ProvinceID,
		}
		if req.CreatedAtGte.IsZero() {
			filter.CreatedAtGte = nil
		}
		if req.CreatedAtLte.IsZero() {
			filter.CreatedAtLte = nil
		}

		stores, total, err := h.StoreInteractor.GetStores(ctx, pager, filter)
		if err != nil {
			err = fmt.Errorf("failed to list stores: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, &types.StoreListResp{Stores: stores, Total: total})
	}
}

// StoreSimpleUpdate 更新门店单个字段信息
//
//	@Summary		更新门店单个字段信息
//	@Description	修改门店状态，
//	@Tags			商户管理-门店
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.StoreSimpleUpdateReq	true	"更新门店单个字段信息请求"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/merchant/store/{id} [patch]
func (h *StoreHandler) StoreSimpleUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.StoreSimpleUpdate")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreSimpleUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		storeIDStr := c.Param("id")
		storeID, err := uuid.Parse(storeIDStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		updateParams := &domain.UpdateStoreParams{ID: storeID, Status: req.Status}
		if err := h.StoreInteractor.StoreSimpleUpdate(ctx, req.SimpleUpdateType, updateParams); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to simple update store: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}

func (h *StoreHandler) checkEditErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrStoreBusinessHoursConflict):
		return errorx.New(http.StatusBadRequest, errcode.StoreBusinessHoursConflict, err)
	case errors.Is(err, domain.ErrStoreBusinessHoursTimeInvalid):
		return errorx.New(http.StatusBadRequest, errcode.StoreBusinessHoursTimeInvalid, err)
	case errors.Is(err, domain.ErrStoreDiningPeriodConflict):
		return errorx.New(http.StatusBadRequest, errcode.StoreDiningPeriodConflict, err)
	case errors.Is(err, domain.ErrStoreDiningPeriodTimeInvalid):
		return errorx.New(http.StatusBadRequest, errcode.StoreDiningPeriodTimeInvalid, err)
	case errors.Is(err, domain.ErrStoreDiningPeriodNameExists):
		return errorx.New(http.StatusBadRequest, errcode.StoreDiningPeriodNameExists, err)
	case errors.Is(err, domain.ErrStoreShiftTimeConflict):
		return errorx.New(http.StatusBadRequest, errcode.StoreShiftTimeConflict, err)
	case errors.Is(err, domain.ErrStoreShiftTimeTimeInvalid):
		return errorx.New(http.StatusBadRequest, errcode.StoreShiftTimeTimeInvalid, err)
	case errors.Is(err, domain.ErrStoreShiftTimeNameExists):
		return errorx.New(http.StatusBadRequest, errcode.StoreShiftTimeNameExists, err)
	case domain.IsConflict(err):
		return errorx.New(http.StatusConflict, errcode.StoreNameExists, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("failed to process store: %w", err)
	}
}
