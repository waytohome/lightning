package ginx

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/waytohome/lightning/logx"
)

type handler struct{}

func (h *handler) Group() Group {
	return nil
}

func (h *handler) Method() string {
	return MethodGet
}

func (h *handler) Path() string {
	return "hello"
}

func (h *handler) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "hello"})
	}
}

func TestHandlerWithoutGroup(t *testing.T) {
	RegisterHandler(&handler{})
	InitRouters(":8080", Config{})
}

type group struct{}

func (g *group) Name() string {
	return "group"
}

func (g *group) Middlewares() []gin.HandlerFunc {
	return nil
}

type handler2 struct{}

func (h *handler2) Group() Group {
	return &group{}
}

func (h *handler2) Method() string {
	return MethodPost
}

func (h *handler2) Path() string {
	return "hello"
}

func (h *handler2) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "hello where in group"})
	}
}

func TestHandlerWithGroup(t *testing.T) {
	RegisterHandler(&handler2{})
	InitRouters(":8080", Config{})
}

func init() {
	logx.SetLevel("info")
}
