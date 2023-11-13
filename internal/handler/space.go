package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yema.dev/api"
	"yema.dev/api/space"
	"yema.dev/internal/errcode"
	"yema.dev/internal/service"
)

type SpaceHandler struct {
	*Handler
	spaceService service.SpaceService
}

func NewSpaceHandler(handler *Handler, spaceService service.SpaceService) *SpaceHandler {
	return &SpaceHandler{
		Handler:      handler,
		spaceService: spaceService,
	}
}

func (h *SpaceHandler) Create(ctx *gin.Context) {
	req := space.CreateReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.spaceService.Create(ctx, &req), nil)
}

func (h *SpaceHandler) List(ctx *gin.Context) {
	req := space.ListReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.spaceService.List(ctx, &req)
	api.PageData(ctx, total, items, err)
}

func (h *SpaceHandler) Update(ctx *gin.Context) {
	req := space.UpdateReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.spaceService.Update(ctx, &req), nil)
}

func (h *SpaceHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.spaceService.Delete(ctx, int64(n)), nil)
}
