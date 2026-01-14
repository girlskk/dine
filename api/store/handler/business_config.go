package handler

import (
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

type BusinessConfigHandler struct {
	BusinessConfigInteractor domain.BusinessConfigInteractor
}

func NewBusinessConfigHandler(BusinessConfigInteractor domain.BusinessConfigInteractor) *BusinessConfigHandler {
	return &BusinessConfigHandler{
		BusinessConfigInteractor: BusinessConfigInteractor,
	}
}

func (h *BusinessConfigHandler) Routes(r gin.IRouter) {
	r = r.Group("business/config")
	r.GET("", h.List())
}

func (h *BusinessConfigHandler) NoAuths() []string {
	return []string{}
}

// List
//
//	@Tags		经营管理
//	@Security	BearerAuth
//	@Summary	经营设置列表
//	@Param		data	query		types.BusinessConfigListReq		true	"经营设置列表查询参数"
//	@Success	200		{object}	domain.BusinessConfigSearchRes	"成功"
//	@Router		/business/config [get]
func (h *BusinessConfigHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("BusinessConfigHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.BusinessConfigListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromStoreUserContext(ctx)
		params := domain.BusinessConfigSearchParams{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			Name:       req.Name,
			Group:      req.Group,
		}
		res, err := h.BusinessConfigInteractor.ListBySearch(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list businessConfigs: %w", err)
				c.Error(err)
			}
			return
		}
		response.Ok(c, res)
	}
}

// UpsertConfig
//
//	@Tags		经营管理
//	@Security	BearerAuth
//	@Summary	更新经营设置
//	@Param		data	body	types.BusinessConfigUpsertReq	true	"请求信息"
//	@Success	200
//	@Router		/business/config [put]
func (h *BusinessConfigHandler) UpsertConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("BusinessConfigHandler.UpsertConfig")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.BusinessConfigUpsertReq
		if err := c.BindJSON(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		user := domain.FromStoreUserContext(ctx)
		var configs []*domain.BusinessConfig
		for _, config := range req.Configs {
			var (
				configID uuid.UUID
				err      error
			)
			configID, err = uuid.Parse(config.ID)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			configs = append(configs, &domain.BusinessConfig{
				ID:             uuid.New(),
				SourceConfigID: configID,
				MerchantID:     user.MerchantID,
				Name:           config.Name,
				Group:          config.Group,
				ConfigType:     config.ConfigType,
				Key:            config.Key,
				Value:          config.Value,
				IsDefault:      false,
				Status:         true,
				ModifyStatus:   true,
				Sort:           config.Sort,
				Tip:            config.Tip,
			})
		}
		err := h.BusinessConfigInteractor.UpsertConfig(ctx, configs, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update businessConfig: %w", err)
			c.Error(err)
			return
		}
		response.Ok(c, nil)
	}
}
