package ginx

import (
	"fmt"
	"strings"

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

type Method = func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

func InitRoutersWithConfigure(c confx.Configure) {
	port, _ := c.GetString("server.port", ":8080")
	if port == "" {
		panic("config server.port is empty")
	}
	InitRouters(port)
}

func InitRouters(port string) {
	gin.SetMode(modeMapping[logx.GetLevel()])

	r := gin.New()

	r.Use(Logger(), Recovery())

	pprof.Register(r)

	for _, handler := range handlerMapping {
		var router gin.IRoutes
		group := handler.Group()
		if group != nil {
			// check group
			if _, ok := groupMapping[group.Name()]; !ok {
				groupMapping[group.Name()] = r.Group(group.Name(), group.Middlewares()...)
			}
			router = groupMapping[group.Name()]
		} else {
			router = r
		}
		method := handler.Method(router)
		method(handler.Path(), handler.Handle())
		logx.Info("ginx handler is ready", logx.String("path", getHandlerPath(handler)))
	}

	if err := r.Run(port); err != nil {
		panic(err)
	}
}

func RegisterHandler(handler Handler) {
	key := getHandlerPath(handler)
	if _, ok := handlerMapping[key]; ok {
		panic("duplicate handler found " + key)
	}
	handlerMapping[key] = handler
}

func getHandlerPath(handler Handler) string {
	key := fmt.Sprintf("/%s", handler.Path())
	if handler.Group() != nil {
		key = fmt.Sprintf("/%s/%s", handler.Group().Name(), handler.Path())
	}
	key = strings.ReplaceAll(key, "//", "/")
	return key
}
