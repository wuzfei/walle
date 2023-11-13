package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yema.dev/api"
	"yema.dev/api/server"
	"yema.dev/internal/errcode"
	"yema.dev/internal/service"
)

type ServerHandler struct {
	*Handler
	serverService service.ServerService
}

func NewServerHandler(handler *Handler, serverService service.ServerService) *ServerHandler {
	return &ServerHandler{
		Handler:       handler,
		serverService: serverService,
	}
}

func (h *ServerHandler) Create(ctx *gin.Context) {
	req := server.CreateReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.serverService.Create(ctx, &req), nil)
}

func (h *ServerHandler) List(ctx *gin.Context) {
	req := server.ListReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.serverService.List(ctx, &req)
	api.PageData(ctx, total, items, err)
}

func (h *ServerHandler) Update(ctx *gin.Context) {
	req := server.UpdateReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.serverService.Update(ctx, &req), nil)
}

func (h *ServerHandler) Delete(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.serverService.Delete(ctx, spaceAndId), nil)
}

func (h *ServerHandler) Check(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.serverService.Check(ctx, spaceAndId), nil)
}

func (h *ServerHandler) SetAuthorized(ctx *gin.Context) {
	req := server.SetAuthorizedReq{SpaceId: GetSpaceId(ctx)}

	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.serverService.SetAuthorized(ctx, &req), nil)
}

func (h *ServerHandler) Terminal(ctx *gin.Context) {
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
	if err = h.serverService.Terminal(ctx, wsConn, spaceAndId, Username(ctx)); err != nil {
		h.log.Error("terminal error", zap.Error(err))
	}
}
