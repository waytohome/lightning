package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/waytohome/lightning/ginx"
)

type SwaggerHandler struct {
}

func (h *SwaggerHandler) Group() ginx.Group {
	return nil
}

func (h *SwaggerHandler) Method() string {
	return ginx.MethodGet
}

func (h *SwaggerHandler) Path() string {
	return "/swagger/*any"
}

func (h *SwaggerHandler) Handle() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
