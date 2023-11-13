package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yema.dev/api"
	"yema.dev/pkg/jwt"
)

type Handler struct {
	log *zap.Logger
}

func NewHandler(log *zap.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func IsLogin(ctx *gin.Context) (auth *jwt.TokenClaims, err error) {
	if v, ok := ctx.Value("claims").(*jwt.TokenClaims); ok {
		return v, nil
	}
	return
}

// UserId 当前登陆用户id
func UserId(ctx *gin.Context) int64 {
	auth, err := IsLogin(ctx)
	if err != nil {
		return 0
	}
	return auth.UserId
}

// GetSpaceWithId 当前登陆用户id
func GetSpaceWithId(ctx *gin.Context) (*api.SpaceWithId, error) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("not found id")
	}
	return &api.SpaceWithId{SpaceId: GetSpaceId(ctx), ID: int64(n)}, nil
}

// Username 当前登陆用户id
func Username(ctx *gin.Context) string {
	auth, err := IsLogin(ctx)
	if err != nil {
		return ""
	}
	return auth.Username
}

func GetSpaceId(ctx *gin.Context) int64 {
	n, err := strconv.Atoi(strings.Trim(ctx.GetHeader("Space-Id"), " "))
	if err != nil {
		return 0
	}
	return int64(n)
}

func GetBearerToken(ctx *gin.Context) string {
	tokenStr := ctx.Request.Header.Get("Authorization")
	if len(tokenStr) > 0 && tokenStr[0:7] == "Bearer " {
		return tokenStr[7:]
	}
	//websocket
	tokenStr = ctx.Request.Header.Get("Sec-WebSocket-Protocol")
	res := strings.Split(tokenStr, ",")
	if len(res) == 2 {
		tokenStr = res[0]
		ctx.Request.Header.Set("space-Id", res[1])
	}
	return tokenStr
}

func UpGrader(ctx *gin.Context) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{GetBearerToken(ctx), strconv.Itoa(int(GetSpaceId(ctx)))},
	}
	return upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
}
