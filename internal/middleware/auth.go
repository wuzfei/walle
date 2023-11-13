package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"yema.dev/pkg/jwt"
)

func Auth(j *jwt.JWT, log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) > 0 && tokenStr[0:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}
		//websocket
		tokenStr = ctx.Request.Header.Get("Sec-WebSocket-Protocol")
		res := strings.Split(tokenStr, ",")
		if len(res) == 2 {
			tokenStr = res[0]
			//ctx.Request.Header.Set(SpaceHeaderName, res[1])
		}
		if tokenStr == "" {
			log.Error("Authorization token is empty", zap.String("url", ctx.Request.URL.String()))
			ctx.AbortWithError(401, errors.New("authorization token error"))
			return
		}
		tp, err := j.ValidateToken(tokenStr)
		if err != nil {
			log.Error("Authorization token error", zap.String("url", ctx.Request.URL.String()), zap.String("token", tokenStr), zap.Error(err))
			ctx.AbortWithError(401, err)
			return
		}
		ctx.Set("claims", tp)
		ctx.Next()
	}
}
