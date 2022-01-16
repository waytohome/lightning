package ginx

import "github.com/gin-gonic/gin"

type Handler interface {
	Group() Group
	Method(r gin.IRoutes) Method
	Path() string
	Handle() gin.HandlerFunc
}
