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

type ProductAttrHandler struct {
	ProductAttrInteractor domain.ProductAttrInteractor
}

func NewProductAttrHandler(productAttrInteractor domain.ProductAttrInteractor) *ProductAttrHandler {
	return &ProductAttrHandler{
		ProductAttrInteractor: productAttrInteractor,
	}
}

func (h *ProductAttrHandler) Routes(r gin.IRouter) {
	r = r.Group("product/attr")
	r.POST("", h.Create())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.DELETE("/item/:id", h.DeleteItem())
	r.GET("", h.List())
}

func (h *ProductAttrHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品口味做法
//	@Security	BearerAuth
//	@Summary	创建商品口味做法
//	@Param		data	body	types.ProductAttrCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/attr [post]
func (h *ProductAttrHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductAttrCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)

		attr := &domain.ProductAttr{
			ID:         uuid.New(),
			Name:       req.Name,
			Channels:   req.Channels,
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
		}

		// 转换口味做法项
		if len(req.Items) > 0 {
			itemNames := make(map[string]bool)
			attr.Items = make([]*domain.ProductAttrItem, 0, len(req.Items))
			for _, itemReq := range req.Items {
				// 验证名称在口味做法下唯一
				if itemNames[itemReq.Name] {
					c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, domain.ErrProductAttrItemNameExists))
					return
				}
				itemNames[itemReq.Name] = true
				item := &domain.ProductAttrItem{
					ID:        uuid.New(),
					AttrID:    attr.ID,
					Name:      itemReq.Name,
					Image:     itemReq.Image,
					BasePrice: itemReq.BasePrice,
				}
				attr.Items = append(attr.Items, item)
			}
		}

		err := h.ProductAttrInteractor.Create(ctx, attr)

		if err != nil {
			if errors.Is(err, domain.ErrProductAttrNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductAttrNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create product attr: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Update
//
//	@Tags		商品口味做法
//	@Security	BearerAuth
//	@Summary	更新商品口味做法
//	@Param		id		path	string						true	"口味做法ID"
//	@Param		data	body	types.ProductAttrUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/attr/{id} [put]
func (h *ProductAttrHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取口味做法ID
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.ProductAttrUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		// 将请求数据映射到 domain.ProductAttr
		attr := &domain.ProductAttr{
			ID:       id,
			Name:     req.Name,
			Channels: req.Channels,
		}

		// 验证口味做法项名称不重复
		if len(req.Items) > 0 {
			itemNames := make(map[string]bool)
			for _, itemReq := range req.Items {
				if itemNames[itemReq.Name] {
					c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, domain.ErrProductAttrItemNameExists))
					return
				}
				itemNames[itemReq.Name] = true
				item := &domain.ProductAttrItem{
					ID:        itemReq.ID,
					AttrID:    id,
					Name:      itemReq.Name,
					Image:     itemReq.Image,
					BasePrice: itemReq.BasePrice,
				}
				attr.Items = append(attr.Items, item)
			}
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductAttrInteractor.Update(ctx, attr, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductAttrNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductAttrNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update product attr: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		商品口味做法
//	@Security	BearerAuth
//	@Summary	删除商品口味做法
//	@Param		id	path	string	true	"口味做法ID"
//	@Success	200	"No Content"
//	@Router		/product/attr/{id} [delete]
func (h *ProductAttrHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductAttrInteractor.Delete(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductAttrDeleteHasItems) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.ProductAttrDeleteHasItems, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}

			err = fmt.Errorf("failed to delete product attr: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// DeleteItem
//
//	@Tags		商品口味做法
//	@Security	BearerAuth
//	@Summary	删除商品口味做法项
//	@Param		id	path	string	true	"口味做法项ID"
//	@Success	200	"No Content"
//	@Router		/product/attr/item/{id} [delete]
func (h *ProductAttrHandler) DeleteItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.DeleteItem")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductAttrInteractor.DeleteItem(ctx, id, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductAttrItemDeleteHasProducts) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.ProductAttrItemDeleteHasProducts, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}

			err = fmt.Errorf("failed to delete product attr item: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		商品口味做法
//	@Security	BearerAuth
//	@Summary	获取商品口味做法列表
//	@Success	200	{object}	domain.ProductAttrs	"成功"
//	@Router		/product/attr [get]
func (h *ProductAttrHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductAttrHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		user := domain.FromStoreUserContext(ctx)

		params := domain.ProductAttrSearchParams{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
		}

		res, err := h.ProductAttrInteractor.ListBySearch(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list product attrs: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
