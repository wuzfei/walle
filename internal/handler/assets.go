package handler

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io"
	"mime"
	"path"
	"strings"
)

// AssetsHandler 静态文件处理
type AssetsHandler struct {
	rootFs   *embed.FS
	assetsFs *embed.FS
}

func NewAssetsHandler(rootFs, assetsFs *embed.FS) *AssetsHandler {
	return &AssetsHandler{
		rootFs:   rootFs,
		assetsFs: assetsFs,
	}
}

// Register 注册到路由
func (h *AssetsHandler) Register(e gin.IRouter) {
	//这三个是静态文件的路由
	e.GET("/", func(ctx *gin.Context) {
		fileHandle(ctx, h.rootFs, "index.html")
	})
	e.GET("/:file", func(ctx *gin.Context) {
		fileHandle(ctx, h.rootFs, ctx.Param("file"))
	})
	e.GET("/:file/*child", func(ctx *gin.Context) {
		fileHandle(ctx, h.rootFs, ctx.Param("file")+ctx.Param("child"))
	})
	e.GET("/assets/*file", func(ctx *gin.Context) {
		fileHandle(ctx, h.assetsFs, "assets"+ctx.Param("file"))
	})
}

// fileHandle 静态文件读取
func fileHandle(ctx *gin.Context, fs *embed.FS, file string) {
	file = "web/dist/" + strings.TrimPrefix(file, "/")
	f, err := fs.Open(file)
	if err != nil {
		ctx.AbortWithStatus(404)
		return
	}
	ctx.Header("Content-Type", mime.TypeByExtension(path.Ext(file)))
	_, err = io.Copy(ctx.Writer, f)
	if err != nil {
		_ = ctx.AbortWithError(500, err)
		return
	}
}
