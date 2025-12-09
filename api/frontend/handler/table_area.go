package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

// TableAreaHandler 处理台桌区域相关请求
type TableAreaHandler struct {
	AreaInteractor domain.TableAreaInteractor
}

// NewTableAreaHandler 创建台桌区域处理器
func NewTableAreaHandler(areaInteractor domain.TableAreaInteractor) *TableAreaHandler {
	return &TableAreaHandler{
		AreaInteractor: areaInteractor,
	}
}

func (h *TableAreaHandler) Routes(r gin.IRouter) {
	r = r.Group("/table-area")
	r.POST("/list", h.List())
}

// List 获取台桌区域列表
//
//	@Tags		台桌管理
//	@Security	BearerAuth
//	@Summary	区域列表
//	@Param		data	body		types.TableAreaListReq	true	"请求信息"
//	@Success	200		{object}	domain.AreaSearchRes	"成功"
//	@Router		/table-area/list [post]
func (h *TableAreaHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableAreaHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableAreaHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TableAreaListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromFrontendUserContext(ctx)
		params := domain.AreaSearchParams{
			StoreID: user.StoreID,
		}
		res, err := h.AreaInteractor.PagedListBySearch(ctx, page, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(uerr.BadRequest(err.Error()))
			} else {
				c.Error(err)
			}
			return
		}

		response.Ok(c, res)
	}
}
