package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzfei/go-helper/slices"
	"yema.dev/api"
	"yema.dev/api/environment"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/service"
)

type EnvironmentHandler struct {
	*Handler
	environmentService service.EnvironmentService
}

func NewEnvironmentHandler(handler *Handler, environmentService service.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{
		Handler:            handler,
		environmentService: environmentService,
	}
}

func (h *EnvironmentHandler) Create(ctx *gin.Context) {
	req := environment.CreateReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.environmentService.Create(ctx, &req), nil)
}

func (h *EnvironmentHandler) List(ctx *gin.Context) {
	req := environment.ListReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.environmentService.List(ctx, &req)
	api.PageData(ctx, total, items, err)
}

func (h *EnvironmentHandler) Update(ctx *gin.Context) {
	req := environment.UpdateReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.environmentService.Update(ctx, &req), nil)
}

func (h *EnvironmentHandler) Delete(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.environmentService.Delete(ctx, spaceAndId), nil)
}

func (h *EnvironmentHandler) Options(ctx *gin.Context) {
	params := environment.ListReq{SpaceId: GetSpaceId(ctx)}

	if err := ctx.ShouldBind(&params); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.environmentService.List(ctx, &params)
	if err != nil {
		api.Fail(ctx, err)
		return
	}
	res := api.DataOptions{Total: total, Options: slices.Map(items, func(item *model.Environment, k int) api.DataOption {
		return api.DataOption{
			Text:   item.Name,
			Value:  item.ID,
			Status: item.Status,
		}
	})}
	api.Success(ctx, res)
}
