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
	r.PATCH("/:id", h.DeviceSimpleUpdate())
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
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
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

		user := domain.FromBackendUserContext(ctx)
		device := &domain.Device{
			ID:          uuid.New(),
			MerchantID:  user.MerchantID,
			StoreID:     req.StoreID,
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
			device.StallID = req.DevicePrint.StallID
			device.OrderChannels = req.DevicePrint.OrderChannels
			device.DiningWays = req.DevicePrint.DiningWays
			device.DeviceStallPrintType = req.DevicePrint.DeviceStallPrintType
			device.DeviceStallReceiptType = req.DevicePrint.DeviceStallReceiptType
		default:
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("device type %s is not supported", req.DeviceType)))
			return
		}

		if err := h.DeviceInteractor.Create(ctx, device); err != nil {
			if errors.Is(err, domain.ErrDeviceNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.DeviceNameExists, err))
				return
			}
			if errors.Is(err, domain.ErrDeviceCodeExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.DeviceCodeExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create device: %w", err)
			c.Error(err)
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
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
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
			device.StallID = req.DevicePrint.StallID
			device.OrderChannels = req.DevicePrint.OrderChannels
			device.DiningWays = req.DevicePrint.DiningWays
			device.DeviceStallPrintType = req.DevicePrint.DeviceStallPrintType
			device.DeviceStallReceiptType = req.DevicePrint.DeviceStallReceiptType
		default:
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("device type %s is not supported", req.DeviceType)))
			return
		}
		if err := h.DeviceInteractor.Update(ctx, device); err != nil {
			if errors.Is(err, domain.ErrDeviceNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.DeviceNameExists, err))
				return
			}
			if errors.Is(err, domain.ErrDeviceCodeExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.DeviceCodeExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update device: %w", err)
			c.Error(err)
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
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
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
		err = h.DeviceInteractor.Delete(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNoContent, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete device: %w", err)
			c.Error(err)
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
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
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

		device, err := h.DeviceInteractor.GetDevice(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get device: %w", err)
			c.Error(err)
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
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
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

		user := domain.FromBackendUserContext(ctx)
		pager := req.RequestPagination.ToPagination()
		filter := &domain.DeviceListFilter{
			MerchantID: user.MerchantID,
			StoreID:    req.StoreID,
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

// DeviceSimpleUpdate 更新设备单个字段
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		更新设备单个字段信息
//	@Description	快速切换启用状态
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"设备ID"
//	@Param			data	body	types.DeviceSimpleUpdateReq	true	"更新设备单个字段信息请求信息"
//	@Success		200		"No Content"
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/restaurant/device/{id} [patch]
func (h *DeviceHandler) DeviceSimpleUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.DeviceSimpleUpdate")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.DeviceSimpleUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		device := &domain.Device{ID: id, Enabled: req.Enabled}
		if err := h.DeviceInteractor.DeviceSimpleUpdate(ctx, req.SimpleUpdateType, device); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to simple update device: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}
