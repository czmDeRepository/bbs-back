package handler

import (
	"bbs-back/base/common"
	"fmt"
	"runtime"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func Filter(ctx *context.Context) {

	//request := ctx.Request
	//uri := request.RequestURI
	//if strings.HasPrefix(uri, "/v1/user") {
	//	token := request.Header.Get("token")
	//	if token == "" {
	//		//http.Error(ctx.ResponseWriter, "Method Not Allowed", http.StatusForbidden)
	//		ctx.Output.JSON(common.ErrorDetail(nil, common.ERROR_TOKEN_NOEXIST, common.ERROR_MESSAGE[common.ERROR_TOKEN_NOEXIST]), beego.BConfig.RunMode != beego.PROD, false)
	//		return
	//	}
	//	log.Println("过滤器：toekn", token)
	//	parseToken, err := controllers.ValidateToken(token)
	//	if err != nil || !parseToken.Valid {
	//		//http.Error(ctx.ResponseWriter, err.Error(), http.StatusForbidden)
	//		var data interface{} = nil
	//		if err != nil {
	//			data = err.Error()
	//		}
	//		ctx.Output.JSON(common.ErrorDetail(data, common.ERROR_TOKEN_PARSE, common.ERROR_MESSAGE[common.ERROR_TOKEN_PARSE]), beego.BConfig.RunMode != beego.PROD, false)
	//	} else {
	//		claims := parseToken.Claims.(jwt.MapClaims)
	//		request.Header.Add(enum.CUR_USER_ID, fmt.Sprintf("%v", claims["id"]))
	//		request.Header.Add(enum.CUR_USER_ROLD, fmt.Sprintf("%v", claims["role"]))
	//	}
	//}
}

func RecoverPanic(ctx *context.Context, cfg *beego.Config) {
	if err := recover(); err != nil {
		logs.Critical("the request url is ", ctx.Input.URL())
		logs.Critical("Handler crashed with error", err)
		bbsError, ok := err.(*common.BBSError)
		if ok {
			ctx.Output.JSON(common.ErrorMeWithCode(bbsError.Error(), bbsError.ErrorCode), beego.BConfig.RunMode != beego.PROD, false)
			return
		}
		if err == beego.ErrAbort {
			return
		}
		if !cfg.RecoverPanic {
			panic(err)
		}
		if cfg.EnableErrorsShow {
			if _, ok := beego.ErrorMaps[fmt.Sprint(err)]; ok {
				parseInt, _ := strconv.ParseUint(fmt.Sprint(err), 10, 64)
				beego.Exception(parseInt, ctx)
				return
			}
		}

		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
		}
		ctx.Output.JSON(common.ErrorWithMe(err.(error).Error(), "未知异常！！！"), beego.BConfig.RunMode != beego.PROD, false)
	}
}
