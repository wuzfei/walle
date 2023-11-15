package handler

import (
	"github.com/gin-gonic/gin"
	"yema.dev/api"
	"yema.dev/api/member"
	"yema.dev/internal/errcode"
	"yema.dev/internal/service"
)

type MemberHandler struct {
	*Handler
	memberService service.MemberService
}

func NewMemberHandler(handler *Handler, memberService service.MemberService) *MemberHandler {
	return &MemberHandler{
		Handler:       handler,
		memberService: memberService,
	}
}

func (h *MemberHandler) Store(ctx *gin.Context) {
	req := member.StoreReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.memberService.Store(ctx, &req), nil)
}

func (h *MemberHandler) List(ctx *gin.Context) {
	req := member.ListReq{SpaceId: GetSpaceId(ctx)}
	if err := ctx.ShouldBind(&req); err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := h.memberService.List(ctx, &req)
	api.PageData(ctx, total, items, err)
}

func (h *MemberHandler) Delete(ctx *gin.Context) {
	spaceAndId, err := GetSpaceWithId(ctx)
	if err != nil {
		api.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	api.Response(ctx, h.memberService.Delete(ctx, spaceAndId), nil)
}
