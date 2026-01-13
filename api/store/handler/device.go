package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/store/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

// DeviceHandler handles device APIs.
type DeviceHandler struct {
	DeviceInteractor domain.DeviceInteractor
}

func NewDeviceHandler(deviceInteractor domain.DeviceInteractor) *DeviceHandler {
	return &DeviceHandler{DeviceInteractor: deviceInteractor}
}

func (h *DeviceHandler) Routes(r gin.IRouter) {
	r = r.Group("/restaurant/device")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("/:id", h.Get())
	r.GET("", h.List())
	r.PUT("/:id/enable", h.Enable())
	r.PUT("/:id/disable", h.Disable())
}

// Create 创建设备
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		创建设备
//	@Description	创建设备
//	@Accept			json
//	@Produce		json
//	@Param			data	body	types.DeviceCreateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Router			/restaurant/device [post]
func (h *DeviceHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.DeviceCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		device := &domain.Device{
			ID:          uuid.New(),
			MerchantID:  user.MerchantID,
			StoreID:     user.StoreID,
			Name:        req.Name,
			DeviceType:  req.DeviceType,
			DeviceCode:  req.DeviceCode,
			DeviceBrand: req.DeviceBrand,
			DeviceModel: req.DeviceModel,
			Location:    req.Location,
			Enabled:     req.Enabled,
			SortOrder:   req.SortOrder,
		}
		switch req.DeviceType {
		case domain.DeviceTypeCashier:
			device.OpenCashDrawer = req.DeviceCashier.OpenCashDrawer
		case domain.DeviceTypePrinter:
			device.IP = req.DevicePrint.IP
			device.PaperSize = req.DevicePrint.PaperSize
			device.ConnectType = req.DevicePrint.ConnectType
			device.StallID = req.DevicePrint.StallID
			device.OrderChannels = req.DevicePrint.OrderChannels
			device.DiningWays = req.DevicePrint.DiningWays
			device.DeviceStallPrintType = req.DevicePrint.DeviceStallPrintType
			device.DeviceStallReceiptType = req.DevicePrint.DeviceStallReceiptType
		default:
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("device type %s is not supported", req.DeviceType)))
			return
		}

		if err := h.DeviceInteractor.Create(ctx, device, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Update 更新设备
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		更新设备
//	@Description	更新设备
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"设备ID"
//	@Param			data	body	types.DeviceUpdateReq	true	"请求信息"
//	@Success		200		"No Content"
//	@Router			/restaurant/device/{id} [put]
func (h *DeviceHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.DeviceUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		device := &domain.Device{
			ID:          id,
			Name:        req.Name,
			DeviceType:  req.DeviceType,
			DeviceCode:  req.DeviceCode,
			DeviceBrand: req.DeviceBrand,
			DeviceModel: req.DeviceModel,
			Location:    req.Location,
			Enabled:     req.Enabled,
			SortOrder:   req.SortOrder,
		}
		switch req.DeviceType {
		case domain.DeviceTypeCashier:
			device.OpenCashDrawer = req.DeviceCashier.OpenCashDrawer
		case domain.DeviceTypePrinter:
			device.IP = req.DevicePrint.IP
			device.PaperSize = req.DevicePrint.PaperSize
			device.ConnectType = req.DevicePrint.ConnectType
			device.StallID = req.DevicePrint.StallID
			device.OrderChannels = req.DevicePrint.OrderChannels
			device.DiningWays = req.DevicePrint.DiningWays
			device.DeviceStallPrintType = req.DevicePrint.DeviceStallPrintType
			device.DeviceStallReceiptType = req.DevicePrint.DeviceStallReceiptType
		default:
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("device type %s is not supported", req.DeviceType)))
			return
		}
		if err := h.DeviceInteractor.Update(ctx, device, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Delete 删除设备
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		删除设备
//	@Description	删除设备
//	@Param			id	path	string	true	"设备ID"
//	@Success		200	"No Content"
//	@Success		204	"No Content"
//	@Router			/restaurant/device/{id} [delete]
func (h *DeviceHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromStoreUserContext(ctx)
		err = h.DeviceInteractor.Delete(ctx, id, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Get 获取设备详情
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		获取设备详情
//	@Description	根据设备ID获取详情
//	@Param			id	path		string	true	"设备ID"
//	@Success		200	{object}	response.Response{data=domain.Device}
//	@Router			/restaurant/device/{id} [get]
func (h *DeviceHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		device, err := h.DeviceInteractor.GetDevice(ctx, id, user)
		if err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, device)
	}
}

// List 获取设备列表
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		获取设备列表
//	@Description	分页查询设备列表
//	@Param			data	query		types.DeviceListReq	true	"设备列表查询参数"
//	@Success		200		{object}	response.Response{data=types.DeviceListResp}
//	@Router			/restaurant/device [get]
func (h *DeviceHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.DeviceListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		pager := req.RequestPagination.ToPagination()
		filter := &domain.DeviceListFilter{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			DeviceType: req.DeviceType,
			Status:     req.Status,
			Name:       req.Name,
		}

		devices, total, err := h.DeviceInteractor.GetDevices(ctx, pager, filter, domain.NewDeviceOrderByCreatedAt(true))
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to get devices: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, types.DeviceListResp{Devices: devices, Total: total})
	}
}

// Enable 启用设备
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		启用设备
//	@Description	将设备置为启用
//	@Produce		json
//	@Param			id	path	string	true	"设备ID"
//	@Success		200	"No Content"
//	@Router			/restaurant/device/{id}/enable [put]
func (h *DeviceHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		device := &domain.Device{ID: id, Enabled: true}
		if err := h.DeviceInteractor.DeviceSimpleUpdate(ctx, domain.DeviceSimpleUpdateTypeEnabled, device, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}

// Disable 禁用设备
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		禁用设备
//	@Description	将设备置为禁用
//	@Produce		json
//	@Param			id	path	string	true	"设备ID"
//	@Success		200	"No Content"
//	@Router			/restaurant/device/{id}/disable [put]
func (h *DeviceHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		device := &domain.Device{ID: id, Enabled: false}
		if err := h.DeviceInteractor.DeviceSimpleUpdate(ctx, domain.DeviceSimpleUpdateTypeEnabled, device, user); err != nil {
			c.Error(h.checkErr(err))
			return
		}

		response.Ok(c, nil)
	}
}
func (h *DeviceHandler) checkErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrDeviceNotExists):
		return errorx.New(http.StatusBadRequest, errcode.DeviceNotExists, err)
	case errors.Is(err, domain.ErrDeviceNameExists):
		return errorx.New(http.StatusConflict, errcode.DeviceNameExists, err)
	case errors.Is(err, domain.ErrDeviceCodeExists):
		return errorx.New(http.StatusConflict, errcode.DeviceCodeExists, err)
	case domain.IsNotFound(err):
		return errorx.New(http.StatusNotFound, errcode.NotFound, err)
	case domain.IsParamsError(err):
		return errorx.New(http.StatusBadRequest, errcode.InvalidParams, err)
	default:
		return fmt.Errorf("device handler error: %w", err)
	}
}
