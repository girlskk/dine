package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type OssHandler struct {
	client *oss.Client
}

func NewOssHandler(client *oss.Client) *OssHandler {
	return &OssHandler{
		client: client,
	}
}

func (h *OssHandler) Routes(r gin.IRouter) {
	r = r.Group("/oss")
	r.POST("/token", h.Token())
}

// Token	获取OSS Token
//
//	@Tags		阿里OSS
//	@Security	BearerAuth
//	@Summary	获取OSS Token
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.OssTokenReq	true	"请求参数"
//	@Success	200		{object}	types.OssTokenResp	"成功"
//	@Router		/oss/token [post]
func (h *OssHandler) Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OssHandler.Token")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.OssTokenReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		token, err := h.client.GeneratePolicyToken()
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			} else {
				c.Error(fmt.Errorf("failed to generate policy token: %w", err))
			}
			return
		}

		key := domain.GenerateObjectKey(req.Scene, req.Filename)

		response.Ok(c, &types.OssTokenResp{
			PolicyToken: *token,
			Key:         key,
			ContentDisposition: lo.TernaryF(
				req.ForDownload,
				func() string { return oss.ContentDispositionAttachmentFilename(req.Filename) },
				oss.ContentDispositionInline,
			),
		})
	}
}
