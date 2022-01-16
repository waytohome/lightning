package ginx

import "github.com/gin-gonic/gin"

type Group interface {
	Name() string
	Middlewares() []gin.HandlerFunc
}
