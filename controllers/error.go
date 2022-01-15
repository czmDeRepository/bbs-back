package controllers

import (
	"bbs-back/base/common"
	beego "github.com/beego/beego/v2/server/web"
)

type ErrorController struct {
	BaseController
}

func init()  {
	// 注册错误处理
	beego.ErrorController(&ErrorController{})
}

func (controller *ErrorController) Error404() {
	controller.end(common.ErrorWithCode(nil,404))
}

func (controller *ErrorController) Error500() {
	controller.end(common.Error(nil))
}


func (controller *ErrorController) ErrorDb() {
	controller.end(common.ErrorWithMe(nil, "database is now down"))
}