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

// TableHandler 处理台相关请求
type TableHandler struct {
	TableInteractor domain.TableInteractor
}

func NewTableHandler(tableInteractor domain.TableInteractor) *TableHandler {
	return &TableHandler{
		TableInteractor: tableInteractor,
	}
}

func (h *TableHandler) Routes(r gin.IRouter) {
	r = r.Group("/table")
	r.POST("/list", h.List())
}

// List 获取台桌列表
//
//	@Tags		台桌管理
//	@Security	BearerAuth
//	@Summary	台桌列表
//	@Param		data	body		types.TableListReq		true	"请求信息"
//	@Success	200		{object}	domain.TableSearchRes	"成功"
//	@Router		/table/list [post]
func (h *TableHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span, ctx := opentracing.StartSpanFromContext(ctx, "TableHandler.List")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("TableHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.TableListReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		page := upagination.New(req.Page, req.Size)
		user := domain.FromFrontendUserContext(ctx)
		params := domain.TableSearchParams{
			StoreID: user.Store.ID,
			AreaID:  req.AreaID,
		}
		if req.Status != nil {
			params.Status = *req.Status
		}

		res, err := h.TableInteractor.PagedListBySearch(ctx, page, params)
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
