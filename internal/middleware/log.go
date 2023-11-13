package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

func RequestLogMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		st := time.Now()
		var reqBody []byte
		if ctx.Request.Body != nil {
			reqBody, _ = ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		zapFields := make([]zap.Field, 0, 5)
		zapFields = append(zapFields, zap.String("method", ctx.Request.Method))
		zapFields = append(zapFields, zap.String("url", ctx.Request.URL.String()))
		ctx.Next()
		zapFields = append(zapFields, zap.Duration("runTime", time.Now().Sub(st)))
		zapFields = append(zapFields, zap.ByteString("reqBody", reqBody))
		zapFields = append(zapFields, zap.ByteString("resBody", blw.body.Bytes()))
		log.Debug("Request", zapFields...)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
