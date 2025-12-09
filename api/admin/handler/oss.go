package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/api/admin/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ali/oss"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	uerr "gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/errors"
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
		span, ctx := opentracing.StartSpanFromContext(ctx, "OssHandler.Token")
		defer span.Finish()
		logger := logging.FromContext(ctx).Named("OssHandler.Token")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.OssTokenReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(uerr.BadRequest(err.Error()))
			return
		}

		token, err := h.client.GeneratePolicyToken()
		if err != nil {
			if msg, ok := domain.GetParamsErrorMessage(err); ok {
				c.Error(uerr.BadRequest(msg))
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
