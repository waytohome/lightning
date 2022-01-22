package ginx

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler interface {
	Group() Group
	Method() string
	Path() string
	Handle() gin.HandlerFunc
}

type SwaggerHandler struct{}

func (h *SwaggerHandler) Group() Group {
	return nil
}

func (h *SwaggerHandler) Method() string {
	return MethodGet
}

func (h *SwaggerHandler) Path() string {
	return "/swagger/*any"
}

func (h *SwaggerHandler) Handle() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
