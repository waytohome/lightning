package ginx

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Endpoints = []gin.HandlerFunc
type Endpoint = gin.HandlerFunc

type Handler interface {
	Middlewares() Endpoints
	Method() Method
	Path() string
	Handle() Endpoint
}

type SwaggerHandler struct{}

func (h *SwaggerHandler) Middlewares() Endpoints {
	return nil
}

func (h *SwaggerHandler) Method() Method {
	return MethodGet
}

func (h *SwaggerHandler) Path() string {
	return "/swagger/*any"
}

func (h *SwaggerHandler) Handle() Endpoint {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
