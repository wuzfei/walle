package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yema.dev/api"
	"yema.dev/api/deploy"
	"yema.dev/internal/errcode"
	"yema.dev/internal/service"
)

type DeployHandler struct {
	*Handler
	deployService service.DeployService
}

func NewDeployHandler(handler *Handler, deployService service.DeployService) *DeployHandler {
	return &DeployHandler{
		Handler:       handler,
		deployService: deployService,
	}
}

func (h *DeployHandler) Create(ctx *gin.Context) {
	req := deploy.CreateReq{SpaceId: GetSpaceId(ctx), UserId: UserId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.deployService.Create(ctx, &req), nil)
}

func (h *DeployHandler) List(ctx *gin.Context) {
	req := deploy.ListReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.deployService.List(ctx, &req)
	api.PageData(ctx, total, items, err)
}

func (h *DeployHandler) Detail(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	data, err := h.deployService.Detail(ctx, spaceAndId)
	api.Response(ctx, err, data)
}

// Audit 审核
func (h *DeployHandler) Audit(ctx *gin.Context) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams)
		return
	}
	req := deploy.AuditReq{SpaceId: GetSpaceId(ctx), AuditUserId: UserId(ctx), ID: int64(n)}

	if err = ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.deployService.Audit(ctx, &req), nil)
}

// Release 发布
func (h *DeployHandler) Release(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	err = h.deployService.Release(ctx, spaceAndId, UserId(ctx))
	api.Response(ctx, err, nil)
}

// StopRelease 中止发布
func (h *DeployHandler) Stop(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	err = h.deployService.Stop(ctx, spaceAndId)
	api.Response(ctx, err, nil)
}

// Rollback 回滚
func (h *DeployHandler) Rollback(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	err = h.deployService.Rollback(ctx, spaceAndId)
	api.Response(ctx, err, nil)
}

// Console 发布执行记录
//func (h *DeployHandler) Console(ctx *gin.Context) {
//	spaceAndId, err := GetSpaceWithId(ctx)
//	if err != nil {
//		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
//		return
//	}
//	wsConn, err := UpGrader(ctx)
//	if err != nil {
//		api.Fail(ctx, err)
//		return
//	}
//	defer func() {
//		_ = wsConn.Close()
//	}()
//	h.deployService.Console(ctx, wsConn, spaceAndId)
//}
