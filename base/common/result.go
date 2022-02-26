package common

type Result struct {
	Success bool   `json:"success"`
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

const (
	FAIL                = -1
	SUCCESS             = 0
	ERROR_PARAM         = 1000
	ERROR_TOKEN_CREATE  = 2000
	ERROR_TOKEN_PARSE   = 2001
	ERROR_TOKEN_NOEXIST = 2002
	ERROR_TOKEN_EXPIRE  = 2003
	ERROR_CURRENT_USER  = 2004
	ERROR_DB_LIMIT      = 3000
	ERROR_EMAIL_EXISTED = 3030
	ERROR_POWER         = 5000
	ERROR_TIME_LIMIT    = 6000
)

var ERROR_MESSAGE = map[int32]string{
	FAIL:                "系统错误",
	ERROR_PARAM:         "请求参数错误",
	ERROR_TOKEN_CREATE:  "创建TOKEN异常",
	ERROR_TOKEN_PARSE:   "TOKEN解析失败",
	ERROR_TOKEN_NOEXIST: "TOKEN不存在",
	ERROR_TOKEN_EXPIRE:  "TOKEN已经失效",
	ERROR_CURRENT_USER:  "非法用户操作",
	ERROR_DB_LIMIT:      "数据操作限制",
	ERROR_EMAIL_EXISTED: "邮箱已存在",
	ERROR_POWER:         "权限不足",
	ERROR_TIME_LIMIT:    "操作太频繁",
}

func SuccessWithData(data interface{}) (res Result) {
	res.Data = data
	res.Code = SUCCESS
	res.Success = true
	return
}

func Success() (res Result) {
	return SuccessWithData(nil)
}

func Error(err error) (res Result) {
	return ErrorDetail(FAIL, err.Error())
}

func ErrorMeWithCode(message string, errCode int32) Result {
	return ErrorDetail(errCode, message)
}

func ErrorWithCode(errCode int32) Result {
	return ErrorDetail(errCode, ERROR_MESSAGE[errCode])
}

func ErrorWithMe(message string) Result {
	return ErrorDetail(FAIL, message)
}

func ErrorDetail(errCode int32, message string) (res Result) {
	res.Code = errCode
	res.Message = message
	return res
}
