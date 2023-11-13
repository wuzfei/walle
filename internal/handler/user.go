package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuzfei/go-helper/slices"
	"strconv"
	"yema.dev/api"
	"yema.dev/api/user"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/service"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	req := user.LoginReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := h.userService.Login(ctx, &req)
	api.Response(ctx, err, res)
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	api.Response(ctx, h.userService.Logout(ctx, UserId(ctx)), nil)
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	req := user.RefreshTokenReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := h.userService.RefreshToken(ctx, &req)
	api.Response(ctx, err, res)
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	spaceAndId := api.SpaceWithId{
		SpaceId: GetSpaceId(ctx),
		ID:      UserId(ctx),
	}
	res, err := h.userService.GetProfile(ctx, &spaceAndId)
	api.Response(ctx, err, res)
}

func (h *UserHandler) Create(ctx *gin.Context) {
	req := user.CreateReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.userService.Create(ctx, &req), nil)
}

func (h *UserHandler) List(ctx *gin.Context) {
	req := user.ListReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.userService.List(ctx, &req)
	api.PageData(ctx, total, items, err)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	req := user.UpdateReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.userService.Update(ctx, &req), nil)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.userService.Delete(ctx, int64(n)), nil)
}

func (h *UserHandler) Options(ctx *gin.Context) {
	req := user.ListReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.userService.List(ctx, &req)
	if err != nil {
		api.Fail(ctx, err)
		return
	}
	res := api.DataOptions{Total: total, Options: slices.Map(items, func(item *model.User, k int) api.DataOption {
		return api.DataOption{
			Text:   fmt.Sprintf("%s(%s)", item.Username, item.Email),
			Value:  item.ID,
			Status: item.Status,
		}
	})}
	api.Success(ctx, res)
}
