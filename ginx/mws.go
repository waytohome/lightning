package ginx

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/waytohome/lightning/logx"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		duration := time.Since(start)
		logx.Info("handle request",
			logx.String("method", ctx.Request.Method),
			logx.String("path", path),
			logx.Int("status", ctx.Writer.Status()),
			logx.String("query", query),
			logx.String("ip", ctx.ClientIP()),
			logx.String("user-agent", ctx.Request.UserAgent()),
			logx.Duration("duration", duration),
		)
	}
}

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
			if brokenPipe {
				logx.Error(ctx.Request.URL.Path,
					logx.Any("error", err),
					logx.String("request", string(httpRequest)),
				)
				// If the connection is dead, we can't write a status to it.
				_ = ctx.Error(err.(error)) // nolint: errcheck
				ctx.Abort()
				return
			}
			logx.Error("recovery from panic",
				logx.Any("error", err),
				logx.String("request", string(httpRequest)),
			)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}()
		ctx.Next()
	}
}
