package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/api/backend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type ProfitDistributionRuleHandler struct {
	ProfitDistributionRuleInteractor domain.ProfitDistributionRuleInteractor
}

func NewProfitDistributionRuleHandler(profitDistributionRuleInteractor domain.ProfitDistributionRuleInteractor) *ProfitDistributionRuleHandler {
	return &ProfitDistributionRuleHandler{
		ProfitDistributionRuleInteractor: profitDistributionRuleInteractor,
	}
}

func (h *ProfitDistributionRuleHandler) Routes(r gin.IRouter) {
	r = r.Group("profit/distribution/rule")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.POST("/:id/enable", h.Enable())
	r.POST("/:id/disable", h.Disable())
	r.GET("", h.List())
}

func (h *ProfitDistributionRuleHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		分账方案
//	@Security	BearerAuth
//	@Summary	创建分账方案
//	@Param		data	body	types.ProfitDistributionRuleCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/profit/distribution/rule [post]
func (h *ProfitDistributionRuleHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionRuleHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProfitDistributionRuleCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		req.StoreIDs = lo.Uniq(req.StoreIDs)

		// 构建 domain.ProfitDistributionRule
		rule := &domain.ProfitDistributionRule{
			ID:            uuid.New(),
			MerchantID:    user.GetMerchantID(),
			Name:          req.Name,
			SplitRatio:    req.SplitRatio,
			BillingCycle:  req.BillingCycle,
			EffectiveDate: req.EffectiveDate,
			ExpiryDate:    req.ExpiryDate,
			Status:        domain.ProfitDistributionRuleStatusDisabled, // 默认禁用状态
			StoreCount:    len(req.StoreIDs),
			Stores: lo.Map(req.StoreIDs, func(storeID uuid.UUID, _ int) *domain.StoreSimple {
				return &domain.StoreSimple{
					ID: storeID,
				}
			}),
		}

		err := h.ProfitDistributionRuleInteractor.Create(ctx, rule, user)
		if err != nil {
			if errors.Is(err, domain.ErrProfitDistributionRuleNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if errors.Is(err, domain.ErrProfitDistributionRuleStoreBound) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create profit distribution rule: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		分账方案
//	@Security	BearerAuth
//	@Summary	更新分账方案
//	@Param		id		path	string									true	"分账方案ID"
//	@Param		data	body	types.ProfitDistributionRuleUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/profit/distribution/rule/{id} [put]
func (h *ProfitDistributionRuleHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionRuleHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取分账方案ID
		idStr := c.Param("id")
		ruleID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.ProfitDistributionRuleUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		req.StoreIDs = lo.Uniq(req.StoreIDs)

		// 构建 domain.ProfitDistributionRule
		rule := &domain.ProfitDistributionRule{
			ID:            ruleID,
			Name:          req.Name,
			SplitRatio:    req.SplitRatio,
			BillingCycle:  req.BillingCycle,
			EffectiveDate: req.EffectiveDate,
			ExpiryDate:    req.ExpiryDate,
			StoreCount:    len(req.StoreIDs),
			Stores: lo.Map(req.StoreIDs, func(storeID uuid.UUID, _ int) *domain.StoreSimple {
				return &domain.StoreSimple{
					ID: storeID,
				}
			}),
		}

		err = h.ProfitDistributionRuleInteractor.Update(ctx, rule, user)
		if err != nil {
			if errors.Is(err, domain.ErrProfitDistributionRuleNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if errors.Is(err, domain.ErrProfitDistributionRuleStoreBound) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update profit distribution rule: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		分账方案
//	@Security	BearerAuth
//	@Summary	删除分账方案
//	@Param		id	path	string	true	"分账方案ID"
//	@Success	200
//	@Router		/profit/distribution/rule/{id} [delete]
func (h *ProfitDistributionRuleHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionRuleHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取分账方案ID
		idStr := c.Param("id")
		ruleID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.ProfitDistributionRuleInteractor.Delete(ctx, ruleID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete profit distribution rule: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Enable
//
//	@Tags		分账方案
//	@Security	BearerAuth
//	@Summary	启用分账方案
//	@Param		id	path	string	true	"分账方案ID"
//	@Success	200
//	@Router		/profit/distribution/rule/{id}/enable [post]
func (h *ProfitDistributionRuleHandler) Enable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionRuleHandler.Enable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取分账方案ID
		idStr := c.Param("id")
		ruleID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.ProfitDistributionRuleInteractor.Enable(ctx, ruleID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to enable profit distribution rule: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Disable
//
//	@Tags		分账方案
//	@Security	BearerAuth
//	@Summary	禁用分账方案
//	@Param		id	path	string	true	"分账方案ID"
//	@Success	200
//	@Router		/profit/distribution/rule/{id}/disable [post]
func (h *ProfitDistributionRuleHandler) Disable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionRuleHandler.Disable")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取分账方案ID
		idStr := c.Param("id")
		ruleID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromBackendUserContext(ctx)
		err = h.ProfitDistributionRuleInteractor.Disable(ctx, ruleID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to disable profit distribution rule: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		分账方案
//	@Security	BearerAuth
//	@Summary	查询分账方案列表
//	@Param		data	query		types.ProfitDistributionRuleListReq		true	"请求信息"
//	@Success	200		{object}	domain.ProfitDistributionRuleSearchRes	"成功"
//	@Router		/profit/distribution/rule [get]
func (h *ProfitDistributionRuleHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProfitDistributionRuleHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProfitDistributionRuleListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromBackendUserContext(ctx)

		params := domain.ProfitDistributionRuleSearchParams{
			MerchantID: user.MerchantID,
			Name:       req.Name,
			Status:     req.Status,
		}

		res, err := h.ProfitDistributionRuleInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list profit distribution rules: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
