package ginx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/waytohome/lightning/confx"
	"github.com/waytohome/lightning/logx"
)

var (
	handlerMapping = make(map[string]Handler)
	groupMapping   = make(map[string]gin.IRoutes)

	modeMapping = map[string]string{
		"debug": gin.DebugMode,
		"info":  gin.ReleaseMode,
		"warn":  gin.ReleaseMode,
		"error": gin.ReleaseMode,
	}
)

func InitRoutersWithConfigure(c confx.Configure) {
	port, _ := c.GetString("server.port", ":8080")
	if port == "" {
		panic("config server.port is empty")
	}
	timeout, _ := c.GetInt("server.timeout", 3)
	needPprof, _ := c.GetBool("server.pprof", false)
	reqMaxSize, _ := c.GetInt("server.request-max-size", 10)
	allowOrigins, _ := c.GetString("server.allow-origins", "")
	InitRouters(port, Config{
		Timeout:      timeout,
		NeedPprof:    needPprof,
		ReqSizeLimit: reqMaxSize,
		AllowOrigins: strings.Split(allowOrigins, ","),
	})
}

type Config struct {
	Timeout      int      // 超时时间
	ReqSizeLimit int      // 请求体最大值(M)
	NeedPprof    bool     // 是否开启pprof
	AllowOrigins []string // 跨域配置
}

func InitRouters(port string, conf Config) {
	gin.SetMode(modeMapping[logx.GetLevel()])

	r := gin.New()

	// 重写路由日志
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, _ string, _ int) {
		logx.Debug("ginx register handler", logx.String("method", httpMethod), logx.String("url", absolutePath))
	}

	// 内建中间件
	var buildInMws []gin.HandlerFunc
	buildInMws = append(buildInMws, Logger())
	buildInMws = append(buildInMws, Recovery())
	if len(conf.AllowOrigins) > 0 && conf.AllowOrigins[0] != "" {
		buildInMws = append(buildInMws, CORS(conf.AllowOrigins))
	}
	buildInMws = append(buildInMws, SizeLimit(conf.ReqSizeLimit))
	buildInMws = append(buildInMws, Timeout(conf.Timeout))
	r.Use(buildInMws...)

	if conf.NeedPprof {
		pprof.Register(r)
	}

	for _, handler := range handlerMapping {
		var router gin.IRoutes
		group := handler.Group()
		if group != nil {
			// check group
			if strings.TrimSpace(group.Name()) == "" {
				panic("group name is blank which is not allowed")
			}
			if _, ok := groupMapping[group.Name()]; !ok {
				groupMapping[group.Name()] = r.Group(group.Name(), group.Middlewares()...)
			}
			router = groupMapping[group.Name()]
		} else {
			router = r
		}
		execMethod := handler.Method().getExecMethod(router)
		var handleFuncs []gin.HandlerFunc
		handleFuncs = append(handleFuncs, handler.Middlewares()...)
		handleFuncs = append(handleFuncs, handler.Handle())
		execMethod(handler.Path(), handleFuncs...)
	}

	// 优雅关停
	srv := &http.Server{Addr: port, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logx.Warn("server shutdown", logx.String("msg", err.Error()))
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logx.Warn("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logx.Error("server forced to shutdown", logx.String("err", err.Error()))
	}
}

func RegisterHandler(handler Handler) {
	key := getHandlerKey(handler)
	if _, ok := handlerMapping[key]; ok {
		panic("duplicate handler found " + key)
	}
	handlerMapping[key] = handler
}

func getHandlerKey(handler Handler) string {
	key := ""
	if handler.Group() != nil && handler.Path() != "" {
		key = fmt.Sprintf("/%s/%s", handler.Group().Name(), handler.Path())
	} else if handler.Path() != "" {
		key = fmt.Sprintf("/%s", handler.Path())
	} else {
		key = fmt.Sprintf("/%s", handler.Group().Name())
	}
	key = strings.ReplaceAll(key, "//", "/")
	key = fmt.Sprintf("%s-%s", handler.Method(), key)
	return key
}
