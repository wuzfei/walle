package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzfei/go-helper/slices"
	"go.uber.org/zap"
	"yema.dev/api"
	"yema.dev/api/project"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/service"
)

type ProjectHandler struct {
	*Handler
	projectService service.ProjectService
}

func NewProjectHandler(handler *Handler, projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		Handler:        handler,
		projectService: projectService,
	}
}

func (h *ProjectHandler) Create(ctx *gin.Context) {
	req := project.CreateReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.projectService.Create(ctx, &req), nil)
}

func (h *ProjectHandler) List(ctx *gin.Context) {
	req := project.ListReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.projectService.List(ctx, &req)
	api.PageData(ctx, total, items, err)
}

func (h *ProjectHandler) Update(ctx *gin.Context) {
	req := project.UpdateReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.projectService.Update(ctx, &req), nil)
}

func (h *ProjectHandler) Delete(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.projectService.Delete(ctx, spaceAndId), nil)
}

func (h *ProjectHandler) Detail(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := h.projectService.Detail(ctx, spaceAndId)
	api.Response(ctx, err, res)
}

func (h *ProjectHandler) Options(ctx *gin.Context) {
	req := project.ListReq{SpaceId: GetSpaceId(ctx)}

	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.projectService.List(ctx, &req)
	if err != nil {
		api.Fail(ctx, err)
		return
	}
	res := api.DataOptions{Total: total, Options: slices.Map(items, func(item *model.Project, k int) api.DataOption {
		return api.DataOption{
			Text:   item.Name,
			Value:  item.ID,
			Status: item.Status,
		}
	})}
	api.Success(ctx, res)
}

// Detection 检测项目, websocket连接
func (h *ProjectHandler) Detection(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}

	wsConn, err := UpGrader(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrServer.Wrap(err))
	}
	defer func() {
		_ = wsConn.Close()
	}()
	if err = h.projectService.DetectionWs(ctx, wsConn, spaceAndId); err != nil {
		h.log.Error("DetectionWs return error", zap.Error(err))
	}
}

// Branches 分支列表
func (h *ProjectHandler) Branches(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := h.projectService.GetBranches(ctx, spaceAndId)
	api.Response(ctx, err, res)
}

// Tags tags列表
func (h *ProjectHandler) Tags(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := h.projectService.GetTags(ctx, spaceAndId)
	api.Response(ctx, err, res)
}

// Commits 提交记录
func (h *ProjectHandler) Commits(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	branch := ctx.Query("branch")
	if err != nil || branch == "" {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := h.projectService.GetCommits(ctx, spaceAndId, ctx.Query("branch"))
	api.Response(ctx, err, res)
}
