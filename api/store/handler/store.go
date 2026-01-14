package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiguang.dev/pos-dine/dine/api/store/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

// StoreHandler handles store related routes for store backend.
type StoreHandler struct {
	StoreInteractor domain.StoreInteractor
}

func NewStoreHandler(storeInteractor domain.StoreInteractor) *StoreHandler {
	return &StoreHandler{StoreInteractor: storeInteractor}
}

func (h *StoreHandler) Routes(r gin.IRouter) {
	r = r.Group("store")
	r.PUT("", h.UpdateStore())
	r.GET("", h.GetStore())
	r.PUT("/enable", h.Enable())
	r.PUT("/disable", h.Disable())
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
//	@Router			/store [put]
func (h *StoreHandler) UpdateStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.UpdateStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromStoreUserContext(ctx)

		var req types.UpdateStoreReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		domainStore := &domain.UpdateStoreParams{
			ID:                      user.StoreID,
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

		if err := h.StoreInteractor.UpdateStore(ctx, domainStore, user); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(h.checkErr(err))
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
//	@Router			/store [get]
func (h *StoreHandler) GetStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.GetStore")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromStoreUserContext(ctx)

		domainStore, err := h.StoreInteractor.GetStore(ctx, user.StoreID, user)
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

// Enable 启用门店
//
//	@Summary		启用门店
//	@Description	将门店状态置为营业
//	@Tags			门店管理
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	"No Content"
//	@Router			/store/enable [put]
func (h *StoreHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromStoreUserContext(ctx)

		updateParams := &domain.UpdateStoreParams{ID: user.StoreID, Status: domain.StoreStatusOpen}
		if err := h.StoreInteractor.StoreSimpleUpdate(ctx, domain.StoreSimpleUpdateFieldStatus, updateParams, user); err != nil {
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
//	@Success		200	"No Content"
//	@Router			/store/disable [put]
func (h *StoreHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("StoreHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromStoreUserContext(ctx)

		updateParams := &domain.UpdateStoreParams{ID: user.StoreID, Status: domain.StoreStatusClosed}
		if err := h.StoreInteractor.StoreSimpleUpdate(ctx, domain.StoreSimpleUpdateFieldStatus, updateParams, user); err != nil {
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

func (h *StoreHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrStoreNotExists):
		return errorx.New(http.StatusBadRequest, errcode.StoreNotExists, err)
	case errors.Is(err, domain.ErrStoreNameExists):
		return errorx.New(http.StatusConflict, errcode.StoreNameExists, err)
	case errors.Is(err, domain.ErrUserExists):
		return errorx.New(http.StatusConflict, errcode.UserNameExists, err)
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
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("failed to process store: %w", err)
	}
}
