package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// DeviceHandler handles device APIs.
type DeviceHandler struct {
	DeviceInteractor domain.DeviceInteractor
}

func NewDeviceHandler(deviceInteractor domain.DeviceInteractor) *DeviceHandler {
	return &DeviceHandler{DeviceInteractor: deviceInteractor}
}

func (h *DeviceHandler) Routes(r gin.IRouter) {
	r = r.Group("/device")
	r.GET("/:id", h.Get())
	r.GET("", h.List())
}

// Get 获取设备详情
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		获取设备详情
//	@Description	根据设备ID获取详情
//	@Param			id	path		string	true	"设备ID"
//	@Success		200	{object}	response.Response{data=domain.Device}
//	@Router			/device/{id} [get]
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

		user := domain.FromFrontendUserContext(ctx)
		device, err := h.DeviceInteractor.GetDevice(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrDeviceNotExists) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.DeviceNotExists, err))
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
//	@Router			/device [get]
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

		user := domain.FromFrontendUserContext(ctx)

		filter := &domain.DeviceListFilter{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			DeviceType: req.DeviceType,
			Status:     req.Status,
		}

		pager := upagination.New(1, upagination.MaxSize)
		devices, total, err := h.DeviceInteractor.GetDevices(ctx, pager, filter)
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

// SyncStoreDeviceStatus 同步门店设备状态
//
//	@Tags			设备管理
//	@Security		BearerAuth
//	@Summary		同步门店设备状态
//	@Description	根据设备编码列表同步门店设备在线状态
//	@Param			data	body		types.SyncStoreDeviceStatusReq	true	"同步门店设备状态请求参数"
//	@Success		200		{object}	response.Response{data=object}
//	@Router			/device/sync_store_device_status [post]
func (h *DeviceHandler) SyncStoreDeviceStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("DeviceHandler.SyncStoreDeviceStatus")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.SyncStoreDeviceStatusReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromFrontendUserContext(ctx)

		err := h.DeviceInteractor.SyncStoreDeviceStatus(ctx, user.MerchantID, user.StoreID, req.DeviceCodes...)
		if err != nil {
			err = fmt.Errorf("failed to sync store device status: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}
