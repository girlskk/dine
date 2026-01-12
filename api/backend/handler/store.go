package handler

import (
	"errors"
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
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

// StoreHandler handles store related routes for brand backend.
type StoreHandler struct {
	StoreInteractor domain.StoreInteractor
}

func NewStoreHandler(storeInteractor domain.StoreInteractor) *StoreHandler {
	return &StoreHandler{StoreInteractor: storeInteractor}
}

func (h *StoreHandler) Routes(r gin.IRouter) {
	r = r.Group("store")
	r.POST("", h.CreateStore())
	r.PUT("/:id", h.UpdateStore())
	r.DELETE("/:id", h.DeleteStore())
	r.GET("/:id", h.GetStore())
	r.GET("/list", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// CreateStore 创建门店
//
//	@Summary		创建门店
//	@Description	创建单个门店
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.CreateStoreReq	true	"创建门店请求"
//	@Success		200		"No Content"
//	@Router			/store [post]
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

		user := domain.FromBackendUserContext(ctx)

		domainStore := &domain.CreateStoreParams{
			MerchantID:              user.MerchantID,
			AdminPhoneNumber:        req.AdminPhoneNumber,
			StoreName:               req.StoreName,
			StoreShortName:          req.StoreShortName,
			StoreCode:               req.StoreCode,
			Status:                  req.Status,
			BusinessModel:           req.BusinessModel,
			BusinessTypeCode:        req.BusinessTypeCode,
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
		domainStore.Address = &domain.Address{
			Country:  req.Address.Country,
			Province: req.Address.Province,
			Address:  req.Address.Address,
			Lng:      req.Address.Lng,
			Lat:      req.Address.Lat,
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
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"门店ID"
//	@Param			data	body	types.UpdateStoreReq	true	"更新门店请求"
//	@Success		200		"No Content"
//	@Router			/store/{id} [put]
func (h *StoreHandler) UpdateStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.UpdateStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeID, err := uuid.Parse(c.Param("id"))
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
			BusinessTypeCode:        req.BusinessTypeCode,
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
			BusinessHours:           req.BusinessHours,
			DiningPeriods:           req.DiningPeriods,
			ShiftTimes:              req.ShiftTimes,
		}
		domainStore.Address = &domain.Address{
			Country:  req.Address.Country,
			Province: req.Address.Province,
			Address:  req.Address.Address,
			Lng:      req.Address.Lng,
			Lat:      req.Address.Lat,
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
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"门店ID"
//	@Success		200	"No Content"
//	@Success		204	"No Content"
//	@Router			/store/{id} [delete]
func (h *StoreHandler) DeleteStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.DeleteStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeID, err := uuid.Parse(c.Param("id"))
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
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"门店ID"
//	@Success		200	{object}	response.Response{data=domain.Store}
//	@Router			/store/{id} [get]
func (h *StoreHandler) GetStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.GetStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeID, err := uuid.Parse(c.Param("id"))
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

// List 门店列表
//
//	@Summary		门店列表
//	@Description	分页查询门店列表
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Produce		json
//	@Param			data	query		types.StoreListReq	true	"门店列表查询参数"
//	@Success		200		{object}	response.Response{data=types.StoreListResp}
//	@Router			/store/list [get]
func (h *StoreHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.StoreListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)

		pager := req.RequestPagination.ToPagination()
		filter := &domain.StoreListFilter{
			StoreName:        req.StoreName,
			MerchantID:       user.MerchantID,
			BusinessTypeCode: req.BusinessTypeCode,
			AdminPhoneNumber: req.AdminPhoneNumber,
			Status:           req.Status,
			BusinessModel:    req.BusinessModel,
		}
		// parse MerchantID if provided
		if req.MerchantID != "" {
			mid, err := uuid.Parse(req.MerchantID)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			filter.MerchantID = mid
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

		stores, total, err := h.StoreInteractor.GetStores(ctx, pager, filter)
		if err != nil {
			c.Error(fmt.Errorf("failed to list stores: %w", err))
			return
		}

		response.Ok(c, &types.StoreListResp{Stores: stores, Total: total})
	}
}

// Enable 启用门店
//
//	@Summary		启用门店
//	@Description	将门店状态置为营业
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"门店ID"
//	@Success		200	"No Content"
//	@Router			/store/{id}/enable [put]
func (h *StoreHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		updateParams := &domain.UpdateStoreParams{ID: storeID, Status: domain.StoreStatusOpen}
		if err := h.StoreInteractor.StoreSimpleUpdate(ctx, domain.StoreSimpleUpdateFieldStatus, updateParams); err != nil {
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

// Disable 禁用门店
//
//	@Summary		禁用门店
//	@Description	将门店状态置为停业
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"门店ID"
//	@Success		200	"No Content"
//	@Router			/store/{id}/disable [put]
func (h *StoreHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		storeID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		updateParams := &domain.UpdateStoreParams{ID: storeID, Status: domain.StoreStatusClosed}
		if err := h.StoreInteractor.StoreSimpleUpdate(ctx, domain.StoreSimpleUpdateFieldStatus, updateParams); err != nil {
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
	case errors.Is(err, domain.ErrUserExists):
		return errorx.New(http.StatusConflict, errcode.UserNameExists, err)
	case errors.Is(err, domain.ErrStoreNameExists):
		return errorx.New(http.StatusConflict, errcode.StoreNameExists, err)
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
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("failed to process store: %w", err)
	}
}
