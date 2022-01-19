package ginx

import "github.com/gin-gonic/gin"

type Handler interface {
	Group() Group
	Method() string
	Path() string
	Handle() gin.HandlerFunc
}
