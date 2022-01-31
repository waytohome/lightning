package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Abort(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(GetHttpCode(code), &resp{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func AbortWithErr(c *gin.Context, err error) {
	err = Adapt(err)
	if err := err; err == nil {
		return
	}
	entry := err.(*errorEntry)
	code := entry.GetCode()
	c.AbortWithStatusJSON(GetHttpCode(code), &resp{
		Code: code,
		Msg:  err.Error(),
		Data: nil,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, resp{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
