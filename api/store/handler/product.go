package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/store/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

type ProductHandler struct {
	ProductInteractor domain.ProductInteractor
}

func NewProductHandler(productInteractor domain.ProductInteractor) *ProductHandler {
	return &ProductHandler{
		ProductInteractor: productInteractor,
	}
}

func (h *ProductHandler) Routes(r gin.IRouter) {
	r = r.Group("product")
	r.POST("", h.Create())
	r.POST("/setmeal", h.CreateSetMeal())
	r.GET("", h.List())
	r.PUT("/:id", h.Update())
	r.PUT("/setmeal/:id", h.UpdateSetMeal())
	r.DELETE("/:id", h.Delete())
	r.PUT("/:id/off-sale", h.OffSale())
	r.PUT("/:id/on-sale", h.OnSale())
	r.GET("/:id", h.GetDetail())
}

func (h *ProductHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	创建普通商品
//	@Param		data	body	types.ProductCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/product [post]
func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)

		// 构建 domain.Product
		product := &domain.Product{
			ID:                uuid.New(),
			Type:              domain.ProductTypeNormal,
			Name:              req.Name,
			MerchantID:        user.MerchantID,
			StoreID:           user.StoreID,
			CategoryID:        req.CategoryID,
			UnitID:            req.UnitID,
			Mnemonic:          req.Mnemonic,
			ShelfLife:         req.ShelfLife,
			SupportTypes:      req.SupportTypes,
			SaleStatus:        req.SaleStatus,
			SaleChannels:      req.SaleChannels,
			EffectiveDateType: req.EffectiveDateType,
			MinSaleQuantity:   req.MinSaleQuantity,
			AddSaleQuantity:   req.AddSaleQuantity,
			InheritTaxRate:    req.InheritTaxRate,
			InheritStall:      req.InheritStall,
			MainImage:         req.MainImage,
			DetailImages:      req.DetailImages,
			Description:       req.Description,
		}

		// 可选字段
		if req.MenuID != nil {
			product.MenuID = *req.MenuID
		}
		if req.EffectiveStartTime != nil {
			product.EffectiveStartTime = req.EffectiveStartTime
		}
		if req.EffectiveEndTime != nil {
			product.EffectiveEndTime = req.EffectiveEndTime
		}
		if req.TaxRateID != nil {
			product.TaxRateID = *req.TaxRateID
		}
		if req.StallID != nil {
			product.StallID = *req.StallID
		}

		// 转换规格关联
		product.SpecRelations = make(domain.ProductSpecRelations, 0, len(req.SpecRelations))
		for _, specRelReq := range req.SpecRelations {
			specRel := &domain.ProductSpecRelation{
				ID:           uuid.New(),
				ProductID:    product.ID,
				SpecID:       specRelReq.SpecID,
				BasePrice:    specRelReq.BasePrice,
				PackingFeeID: specRelReq.PackingFeeID,
				Barcode:      specRelReq.Barcode,
				IsDefault:    specRelReq.IsDefault,
			}

			if specRelReq.MemberPrice != nil {
				specRel.MemberPrice = specRelReq.MemberPrice
			}
			if specRelReq.EstimatedCostPrice != nil {
				specRel.EstimatedCostPrice = specRelReq.EstimatedCostPrice
			}
			if specRelReq.OtherPrice1 != nil {
				specRel.OtherPrice1 = specRelReq.OtherPrice1
			}
			if specRelReq.OtherPrice2 != nil {
				specRel.OtherPrice2 = specRelReq.OtherPrice2
			}
			if specRelReq.OtherPrice3 != nil {
				specRel.OtherPrice3 = specRelReq.OtherPrice3
			}

			product.SpecRelations = append(product.SpecRelations, specRel)
		}

		// 转换口味做法关联
		if len(req.AttrRelations) > 0 {
			product.AttrRelations = make(domain.ProductAttrRelations, 0, len(req.AttrRelations))
			for _, attrRelReq := range req.AttrRelations {
				attrRel := &domain.ProductAttrRelation{
					ID:         uuid.New(),
					ProductID:  product.ID,
					AttrID:     attrRelReq.AttrID,
					AttrItemID: attrRelReq.AttrItemID,
					IsDefault:  attrRelReq.IsDefault,
				}
				product.AttrRelations = append(product.AttrRelations, attrRel)
			}
		}

		// 转换标签
		if len(req.TagIDs) > 0 {
			product.Tags = make(domain.ProductTags, 0, len(req.TagIDs))
			for _, tagID := range req.TagIDs {
				product.Tags = append(product.Tags, &domain.ProductTag{
					ID: tagID,
				})
			}
		}

		err := h.ProductInteractor.Create(ctx, product)

		if err != nil {
			if errors.Is(err, domain.ErrProductNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create product: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// CreateSetMeal
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	创建套餐商品
//	@Param		data	body	types.SetMealCreateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/setmeal [post]
func (h *ProductHandler) CreateSetMeal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.CreateSetMeal")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.SetMealCreateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		// 构建 domain.Product
		product := &domain.Product{
			ID:                uuid.New(),
			Type:              domain.ProductTypeSetMeal,
			Name:              req.Name,
			MerchantID:        user.MerchantID,
			StoreID:           user.StoreID,
			CategoryID:        req.CategoryID,
			UnitID:            req.UnitID,
			Mnemonic:          req.Mnemonic,
			ShelfLife:         req.ShelfLife,
			SupportTypes:      req.SupportTypes,
			SaleStatus:        req.SaleStatus,
			SaleChannels:      req.SaleChannels,
			EffectiveDateType: req.EffectiveDateType,
			MinSaleQuantity:   req.MinSaleQuantity,
			AddSaleQuantity:   req.AddSaleQuantity,
			InheritTaxRate:    req.InheritTaxRate,
			InheritStall:      req.InheritStall,
			MainImage:         req.MainImage,
			DetailImages:      req.DetailImages,
			Description:       req.Description,
		}

		// 可选字段
		if req.MenuID != nil {
			product.MenuID = *req.MenuID
		}
		if req.EffectiveStartTime != nil {
			product.EffectiveStartTime = req.EffectiveStartTime
		}
		if req.EffectiveEndTime != nil {
			product.EffectiveEndTime = req.EffectiveEndTime
		}
		if req.TaxRateID != nil {
			product.TaxRateID = *req.TaxRateID
		}
		if req.StallID != nil {
			product.StallID = *req.StallID
		}
		// 套餐属性
		if req.ComboEstimatedCostPrice != nil {
			product.EstimatedCostPrice = req.ComboEstimatedCostPrice
		}
		if req.ComboDeliveryCostPrice != nil {
			product.DeliveryCostPrice = req.ComboDeliveryCostPrice
		}

		// 转换规格关联
		product.SpecRelations = make(domain.ProductSpecRelations, 0, len(req.SpecRelations))
		for _, specRelReq := range req.SpecRelations {
			specRel := &domain.ProductSpecRelation{
				ID:           uuid.New(),
				ProductID:    product.ID,
				SpecID:       specRelReq.SpecID,
				BasePrice:    specRelReq.BasePrice,
				PackingFeeID: specRelReq.PackingFeeID,
				Barcode:      specRelReq.Barcode,
				IsDefault:    specRelReq.IsDefault,
			}

			if specRelReq.MemberPrice != nil {
				specRel.MemberPrice = specRelReq.MemberPrice
			}
			if specRelReq.EstimatedCostPrice != nil {
				specRel.EstimatedCostPrice = specRelReq.EstimatedCostPrice
			}
			if specRelReq.OtherPrice1 != nil {
				specRel.OtherPrice1 = specRelReq.OtherPrice1
			}
			if specRelReq.OtherPrice2 != nil {
				specRel.OtherPrice2 = specRelReq.OtherPrice2
			}
			if specRelReq.OtherPrice3 != nil {
				specRel.OtherPrice3 = specRelReq.OtherPrice3
			}

			product.SpecRelations = append(product.SpecRelations, specRel)
		}

		// 转换标签
		if len(req.TagIDs) > 0 {
			product.Tags = make(domain.ProductTags, 0, len(req.TagIDs))
			for _, tagID := range req.TagIDs {
				product.Tags = append(product.Tags, &domain.ProductTag{
					ID: tagID,
				})
			}
		}

		// 转换套餐组
		groups := make(domain.SetMealGroups, 0, len(req.Groups))
		for _, groupReq := range req.Groups {
			group := &domain.SetMealGroup{
				ID:            uuid.New(),
				ProductID:     product.ID,
				Name:          groupReq.Name,
				SelectionType: groupReq.SelectionType,
			}
			for _, detailReq := range groupReq.Details {
				detail := &domain.SetMealDetail{
					ID:                 uuid.New(),
					GroupID:            group.ID,
					ProductID:          detailReq.ProductID,
					Quantity:           detailReq.Quantity,
					IsDefault:          detailReq.IsDefault,
					OptionalProductIDs: detailReq.OptionalProductIDs,
				}
				group.Details = append(group.Details, detail)
			}
			groups = append(groups, group)
		}
		product.Groups = groups

		err := h.ProductInteractor.CreateSetMeal(ctx, product)

		if err != nil {
			if errors.Is(err, domain.ErrProductNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to create product: %w", err)
			c.Error(err)
			return
		}
		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	查询商品列表
//	@Param		data	query		types.ProductListReq	true	"请求信息"
//	@Success	200		{object}	domain.ProductSearchRes	"成功"
//	@Router		/product [get]
func (h *ProductHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ProductListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		page := upagination.New(req.Page, req.Size)
		user := domain.FromStoreUserContext(ctx)

		startAt, err := time.Parse(time.DateOnly, req.StartAt)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		endAt, err := time.Parse(time.DateOnly, req.EndAt)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := domain.ProductSearchParams{
			MerchantID: user.MerchantID,
			StoreID:    user.StoreID,
			Name:       req.Name,
			StartAt:    &startAt,
			EndAt:      &endAt,
		}

		// 转换UUID
		if req.CategoryID != "" {
			categoryID, err := uuid.Parse(req.CategoryID)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.CategoryID = categoryID
		}

		if req.StallID != "" {
			stallID, err := uuid.Parse(req.StallID)
			if err != nil {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			params.StallID = stallID
		}

		// 售卖状态
		if req.SaleStatus != "" {
			params.SaleStatus = domain.ProductSaleStatus(req.SaleStatus)
		}

		// 商品类型
		if req.Type != "" {
			params.Type = domain.ProductType(req.Type)
		}

		res, err := h.ProductInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				err = fmt.Errorf("failed to list products: %w", err)
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}

// Update
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	更新普通商品
//	@Param		id		path	string					true	"商品ID"
//	@Param		data	body	types.ProductUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/{id} [put]
func (h *ProductHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取商品ID
		idStr := c.Param("id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid product id: %w", err)))
			return
		}

		var req types.ProductUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)

		// 构建 domain.Product
		product := &domain.Product{
			ID:                productID,
			Type:              domain.ProductTypeNormal,
			Name:              req.Name,
			MerchantID:        user.MerchantID,
			StoreID:           user.StoreID,
			CategoryID:        req.CategoryID,
			UnitID:            req.UnitID,
			Mnemonic:          req.Mnemonic,
			ShelfLife:         req.ShelfLife,
			SupportTypes:      req.SupportTypes,
			SaleStatus:        req.SaleStatus,
			SaleChannels:      req.SaleChannels,
			EffectiveDateType: req.EffectiveDateType,
			MinSaleQuantity:   req.MinSaleQuantity,
			AddSaleQuantity:   req.AddSaleQuantity,
			InheritTaxRate:    req.InheritTaxRate,
			InheritStall:      req.InheritStall,
			MainImage:         req.MainImage,
			DetailImages:      req.DetailImages,
			Description:       req.Description,
		}

		// 可选字段
		if req.MenuID != nil {
			product.MenuID = *req.MenuID
		}
		if req.EffectiveStartTime != nil {
			product.EffectiveStartTime = req.EffectiveStartTime
		}
		if req.EffectiveEndTime != nil {
			product.EffectiveEndTime = req.EffectiveEndTime
		}
		if req.TaxRateID != nil {
			product.TaxRateID = *req.TaxRateID
		}
		if req.StallID != nil {
			product.StallID = *req.StallID
		}

		// 转换规格关联
		product.SpecRelations = make(domain.ProductSpecRelations, 0, len(req.SpecRelations))
		for _, specRelReq := range req.SpecRelations {
			specRel := &domain.ProductSpecRelation{
				ID:           uuid.New(),
				ProductID:    product.ID,
				SpecID:       specRelReq.SpecID,
				BasePrice:    specRelReq.BasePrice,
				PackingFeeID: specRelReq.PackingFeeID,
				Barcode:      specRelReq.Barcode,
				IsDefault:    specRelReq.IsDefault,
			}

			if specRelReq.MemberPrice != nil {
				specRel.MemberPrice = specRelReq.MemberPrice
			}
			if specRelReq.EstimatedCostPrice != nil {
				specRel.EstimatedCostPrice = specRelReq.EstimatedCostPrice
			}
			if specRelReq.OtherPrice1 != nil {
				specRel.OtherPrice1 = specRelReq.OtherPrice1
			}
			if specRelReq.OtherPrice2 != nil {
				specRel.OtherPrice2 = specRelReq.OtherPrice2
			}
			if specRelReq.OtherPrice3 != nil {
				specRel.OtherPrice3 = specRelReq.OtherPrice3
			}

			product.SpecRelations = append(product.SpecRelations, specRel)
		}

		// 转换口味做法关联
		if len(req.AttrRelations) > 0 {
			product.AttrRelations = make(domain.ProductAttrRelations, 0, len(req.AttrRelations))
			for _, attrRelReq := range req.AttrRelations {
				attrRel := &domain.ProductAttrRelation{
					ID:         uuid.New(),
					ProductID:  product.ID,
					AttrID:     attrRelReq.AttrID,
					AttrItemID: attrRelReq.AttrItemID,
					IsDefault:  attrRelReq.IsDefault,
				}
				product.AttrRelations = append(product.AttrRelations, attrRel)
			}
		}

		// 转换标签
		if len(req.TagIDs) > 0 {
			product.Tags = make(domain.ProductTags, 0, len(req.TagIDs))
			for _, tagID := range req.TagIDs {
				product.Tags = append(product.Tags, &domain.ProductTag{
					ID: tagID,
				})
			}
		}

		err = h.ProductInteractor.Update(ctx, product, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update product: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// UpdateSetMeal
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	更新套餐商品
//	@Param		id		path	string					true	"商品ID"
//	@Param		data	body	types.SetMealUpdateReq	true	"请求信息"
//	@Success	200
//	@Router		/product/setmeal/{id} [put]
func (h *ProductHandler) UpdateSetMeal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.UpdateSetMeal")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取商品ID
		idStr := c.Param("id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, fmt.Errorf("invalid product id: %w", err)))
			return
		}

		var req types.SetMealUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)

		// 构建 domain.Product
		product := &domain.Product{
			ID:                productID,
			Type:              domain.ProductTypeSetMeal,
			Name:              req.Name,
			MerchantID:        user.MerchantID,
			StoreID:           user.StoreID,
			CategoryID:        req.CategoryID,
			UnitID:            req.UnitID,
			Mnemonic:          req.Mnemonic,
			ShelfLife:         req.ShelfLife,
			SupportTypes:      req.SupportTypes,
			SaleStatus:        req.SaleStatus,
			SaleChannels:      req.SaleChannels,
			EffectiveDateType: req.EffectiveDateType,
			MinSaleQuantity:   req.MinSaleQuantity,
			AddSaleQuantity:   req.AddSaleQuantity,
			InheritTaxRate:    req.InheritTaxRate,
			InheritStall:      req.InheritStall,
			MainImage:         req.MainImage,
			DetailImages:      req.DetailImages,
			Description:       req.Description,
		}

		// 可选字段
		if req.MenuID != nil {
			product.MenuID = *req.MenuID
		}
		if req.EffectiveStartTime != nil {
			product.EffectiveStartTime = req.EffectiveStartTime
		}
		if req.EffectiveEndTime != nil {
			product.EffectiveEndTime = req.EffectiveEndTime
		}
		if req.TaxRateID != nil {
			product.TaxRateID = *req.TaxRateID
		}
		if req.StallID != nil {
			product.StallID = *req.StallID
		}
		// 套餐属性
		if req.ComboEstimatedCostPrice != nil {
			product.EstimatedCostPrice = req.ComboEstimatedCostPrice
		}
		if req.ComboDeliveryCostPrice != nil {
			product.DeliveryCostPrice = req.ComboDeliveryCostPrice
		}

		// 转换规格关联
		product.SpecRelations = make(domain.ProductSpecRelations, 0, len(req.SpecRelations))
		for _, specRelReq := range req.SpecRelations {
			specRel := &domain.ProductSpecRelation{
				ID:           uuid.New(),
				ProductID:    product.ID,
				SpecID:       specRelReq.SpecID,
				BasePrice:    specRelReq.BasePrice,
				PackingFeeID: specRelReq.PackingFeeID,
				Barcode:      specRelReq.Barcode,
				IsDefault:    specRelReq.IsDefault,
			}

			if specRelReq.MemberPrice != nil {
				specRel.MemberPrice = specRelReq.MemberPrice
			}
			if specRelReq.EstimatedCostPrice != nil {
				specRel.EstimatedCostPrice = specRelReq.EstimatedCostPrice
			}
			if specRelReq.OtherPrice1 != nil {
				specRel.OtherPrice1 = specRelReq.OtherPrice1
			}
			if specRelReq.OtherPrice2 != nil {
				specRel.OtherPrice2 = specRelReq.OtherPrice2
			}
			if specRelReq.OtherPrice3 != nil {
				specRel.OtherPrice3 = specRelReq.OtherPrice3
			}

			product.SpecRelations = append(product.SpecRelations, specRel)
		}

		// 转换标签
		if len(req.TagIDs) > 0 {
			product.Tags = make(domain.ProductTags, 0, len(req.TagIDs))
			for _, tagID := range req.TagIDs {
				product.Tags = append(product.Tags, &domain.ProductTag{
					ID: tagID,
				})
			}
		}

		// 转换套餐组
		groups := make(domain.SetMealGroups, 0, len(req.Groups))
		for _, groupReq := range req.Groups {
			group := &domain.SetMealGroup{
				ID:            uuid.New(),
				ProductID:     product.ID,
				Name:          groupReq.Name,
				SelectionType: groupReq.SelectionType,
			}
			for _, detailReq := range groupReq.Details {
				detail := &domain.SetMealDetail{
					ID:                 uuid.New(),
					GroupID:            group.ID,
					ProductID:          detailReq.ProductID,
					Quantity:           detailReq.Quantity,
					IsDefault:          detailReq.IsDefault,
					OptionalProductIDs: detailReq.OptionalProductIDs,
				}
				group.Details = append(group.Details, detail)
			}
			groups = append(groups, group)
		}
		product.Groups = groups

		err = h.ProductInteractor.UpdateSetMeal(ctx, product, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductNameExists) {
				c.Error(errorx.New(http.StatusConflict, errcode.ProductNameExists, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to update set meal: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// Delete
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	删除商品
//	@Param		id	path	string	true	"商品ID"
//	@Success	200
//	@Router		/product/{id} [delete]
func (h *ProductHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取商品ID
		idStr := c.Param("id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductInteractor.Delete(ctx, productID, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductBelongToSetMeal) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.ProductBelongToSetMeal, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to delete product: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// OffSale
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	停售商品
//	@Param		id	path	string	true	"商品ID"
//	@Success	200
//	@Router		/product/{id}/off-sale [put]
func (h *ProductHandler) OffSale() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.OffSale")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取商品ID
		idStr := c.Param("id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductInteractor.OffSale(ctx, productID, user)
		if err != nil {
			if errors.Is(err, domain.ErrProductBelongToSetMeal) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.ProductBelongToSetMeal, err))
				return
			}
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to off sale product: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// OnSale
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	启售商品
//	@Param		id	path	string	true	"商品ID"
//	@Success	200
//	@Router		/product/{id}/on-sale [put]
func (h *ProductHandler) OnSale() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.OnSale")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取商品ID
		idStr := c.Param("id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		err = h.ProductInteractor.OnSale(ctx, productID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to on sale product: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, nil)
	}
}

// GetDetail
//
//	@Tags		商品管理
//	@Security	BearerAuth
//	@Summary	获取商品详情
//	@Param		id	path		string			true	"商品ID"
//	@Success	200	{object}	domain.Product	"成功"
//	@Router		/product/{id} [get]
func (h *ProductHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("ProductHandler.GetDetail")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		// 从路径参数获取商品ID
		idStr := c.Param("id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		user := domain.FromStoreUserContext(ctx)
		product, err := h.ProductInteractor.GetDetail(ctx, productID, user)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			err = fmt.Errorf("failed to get product detail: %w", err)
			c.Error(err)
			return
		}

		response.Ok(c, product)
	}
}
