package ginx

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	GetMethod     = "GET"
	PostMethod    = "POST"
	DeleteMethod  = "DELETE"
	PatchMethod   = "PATCH"
	PutMethod     = "PUT"
	OptionsMethod = "OPTIONS"
	HeadMethod    = "HEAD"
)

var (
	handlerMapping = make(map[string]Handler)
	methodMapping  = make(map[string]Method)
)

type Method = func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

type Handler interface {
	Method() string
	Prefix() string
	Path() string
	Handle() gin.HandlerFunc
}

func InitRouter(port string) {
	router := gin.Default()

	initMethodMapping(router)

	for _, handler := range handlerMapping {
		fullPath := getFullPath(handler)
		m := methodMapping[handler.Method()]
		m(fullPath, handler.Handle())
	}

	if err := router.Run(port); err != nil {
		panic(err)
	}
}

func RegisterHandler(handler Handler) {
	key := getFullPath(handler)
	if _, ok := handlerMapping[key]; ok {
		panic("duplicate handler found " + key)
	}
	handlerMapping[key] = handler
}

func getFullPath(handler Handler) string {
	fullPath := fmt.Sprintf("/%s/%s", handler.Prefix(), handler.Path())
	fullPath = strings.ReplaceAll(fullPath, "//", "/")
	return fullPath
}

func initMethodMapping(r *gin.Engine) {
	methodMapping[GetMethod] = r.GET
	methodMapping[PostMethod] = r.POST
	methodMapping[PutMethod] = r.PUT
	methodMapping[DeleteMethod] = r.DELETE
	methodMapping[PatchMethod] = r.PATCH
	methodMapping[OptionsMethod] = r.OPTIONS
	methodMapping[HeadMethod] = r.HEAD
}
