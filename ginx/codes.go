package ginx

// code 400
const (
	CodeParamValidateFailed = 400001 // 参数校验失败

	Code400Begin = 400101 // 业务状态码起点
)

// code 401
const (
	CodeAuthFailed = 401001 // 身份认证失败

	Code401Begin = 401101
)

// code 403
const (
	CodeOperateForbidden = 403001 // 禁止操作

	Code403Begin = 403101
)

// code 404
const (
	CodeResourceNotFound = 404001 // 找不到资源

	Code404Begin = 404101
)

// code 500
const (
	CodeServerError = 500001 // 服务器出错

	Code500Begin = 500101
	CodeUnknown  = 500999 // 未知错误
)

func GetHttpCode(code int) int {
	return code / 1000
}
