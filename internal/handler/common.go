package handler

import (
	"github.com/gin-gonic/gin"
	"yema.dev/api"
	"yema.dev/internal/utils"
	"yema.dev/pkg/version"
)

type CommonHandler struct {
	*Handler
}

func NewCommonHandler(handler *Handler) *CommonHandler {
	return &CommonHandler{
		Handler: handler,
	}
}

func (h *CommonHandler) Version(ctx *gin.Context) {
	api.Success(ctx, version.Build)
}

func (h *CommonHandler) System(ctx *gin.Context) {
	s, err := utils.GetServerInfo()
	api.Response(ctx, err, s)
}
