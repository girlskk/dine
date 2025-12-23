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
	// r.PUT("/:id", h.Update())
	// r.DELETE("/:id", h.Delete())
	// r.GET("/:id", h.GetDetail())
	// r.GET("", h.List())
}

func (h *ProductHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		商品
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

		user := domain.FromBackendUserContext(ctx)

		// 构建 domain.Product
		product := &domain.Product{
			ID:                uuid.New(),
			Name:              req.Name,
			MerchantID:        user.MerchantID,
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
