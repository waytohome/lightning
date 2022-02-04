package ginx

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"

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

func CORS(origins []string) gin.HandlerFunc {
	logx.Info("gin allow origins", logx.Strings("origins", origins))
	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	config.AllowMethods = []string{
		MethodGet.String(),
		MethodPost.String(),
		MethodPut.String(),
		MethodDelete.String(),
		MethodOptions.String(),
		MethodPatch.String(),
		MethodHead.String(),
	}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	return cors.New(config)
}

func Timeout(max int) gin.HandlerFunc {
	if max < 3 {
		logx.Warn("gin middleware set timeout is too short, reset to default", logx.Int("expect", max), logx.Int("default", 3))
		max = 3
	}
	handlerFunc := timeout.New(
		timeout.WithTimeout(time.Duration(max)*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, gin.H{"msg": "handle request timeout"})
		},
		))
	return handlerFunc
}

func SizeLimit(max int) gin.HandlerFunc {
	return limits.RequestSizeLimiter(int64(max * 1024 * 1024))
}

func RateLimit(rdb redis.Cmdable, key string, limit redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := redis_rate.NewLimiter(rdb)
		resp, err := limiter.Allow(c, key, limit)
		if err != nil {
			panic(err)
		}
		if resp.Remaining <= 0 {
			c.JSON(403, gin.H{"msg": "request limited"})
			return
		}
		c.Next()
	}
}
