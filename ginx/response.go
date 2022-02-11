package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/waytohome/lightning/logx"
)

type resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Abort(c *gin.Context, code int, msg string) {
	httpCode := GetHttpCode(code)
	logx.Error(msg, logx.Int("status", httpCode), logx.String("path", c.Request.RequestURI))
	c.AbortWithStatusJSON(httpCode, &resp{
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
	httpCode := GetHttpCode(code)
	logx.Error(err.Error(), logx.Int("status", httpCode), logx.String("path", c.Request.RequestURI))
	c.AbortWithStatusJSON(httpCode, &resp{
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
