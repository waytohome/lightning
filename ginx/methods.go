package ginx

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ExecMethod func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

type Method int

const (
	_ Method = iota
	MethodGet
	MethodPost
	MethodPut
	MethodDelete
	MethodOptions
	MethodPatch
	MethodHead
)

func (m Method) String() string {
	switch m {
	case MethodGet:
		return "GET"
	case MethodPost:
		return "POST"
	case MethodPut:
		return "PUT"
	case MethodDelete:
		return "DELETE"
	case MethodOptions:
		return "OPTIONS"
	case MethodPatch:
		return "PATCH"
	case MethodHead:
		return "HEAD"
	default:
		panic(fmt.Sprintf("unrecognized method found, method = %d \n", m))
	}
}

func (m Method) getExecMethod(r gin.IRoutes) ExecMethod {
	switch m {
	case MethodGet:
		return r.GET
	case MethodPost:
		return r.POST
	case MethodPut:
		return r.PUT
	case MethodDelete:
		return r.DELETE
	case MethodOptions:
		return r.OPTIONS
	case MethodPatch:
		return r.PATCH
	case MethodHead:
		return r.HEAD
	default:
		panic(fmt.Sprintf("unrecognized method found, method = %d \n", m))
	}
}
