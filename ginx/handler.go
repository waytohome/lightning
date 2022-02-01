package ginx

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler interface {
	Group() Group
	Middlewares() []gin.HandlerFunc
	Method() Method
	Path() string
	Handle() gin.HandlerFunc
}

type SwaggerHandler struct{}

func (h *SwaggerHandler) Middlewares() []gin.HandlerFunc {
	return nil
}

func (h *SwaggerHandler) Group() Group {
	return nil
}

func (h *SwaggerHandler) Method() Method {
	return MethodGet
}

func (h *SwaggerHandler) Path() string {
	return "/swagger/*any"
}

func (h *SwaggerHandler) Handle() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
