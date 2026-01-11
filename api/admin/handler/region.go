package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type RegionHandler struct {
	CountryInteractor  domain.CountryInteractor
	ProvinceInteractor domain.ProvinceInteractor
}

func NewRegionHandler(countryInteractor domain.CountryInteractor, provinceInteractor domain.ProvinceInteractor) *RegionHandler {
	return &RegionHandler{
		CountryInteractor:  countryInteractor,
		ProvinceInteractor: provinceInteractor,
	}
}

func (h *RegionHandler) Routes(r gin.IRouter) {
	r = r.Group("/region")
	r.GET("/countries", h.ListCountries())
	r.GET("/:id/provinces", h.ListProvinces())
}

// ListCountries
//
//	@Tags		地区
//	@Security	BearerAuth
//	@Summary	获取国家/地区列表
//	@Success	200	{object}	response.Response{data=types.CountryListResp}
//	@Router		/region/countries [get]
func (h *RegionHandler) ListCountries() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RegionHandler.ListCountries")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		countries, err := h.CountryInteractor.GetAllCountries(ctx)
		if err != nil {
			c.Error(errorx.New(http.StatusInternalServerError, errcode.InternalError, err))
			return
		}

		response.Ok(c, types.CountryListResp{
			Countries: countries,
		})
	}
}

// ListProvinces
//
//	@Tags		地区
//	@Security	BearerAuth
//	@Summary	获取指定国家/地区的省/州列表
//	@Param		id	path		string	true	"国家/地区 ID (UUID)"
//	@Success	200	{object}	response.Response{data=types.ProvinceListResp}
//	@Router		/region/{id}/provinces [get]
func (h *RegionHandler) ListProvinces() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("RegionHandler.ListProvinces")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}
		provinces, err := h.ProvinceInteractor.GetProvinces(ctx, id)
		if err != nil {
			c.Error(errorx.New(http.StatusInternalServerError, errcode.InternalError, err))
			return
		}

		response.Ok(c, types.ProvinceListResp{
			Provinces: provinces,
		})
	}
}
