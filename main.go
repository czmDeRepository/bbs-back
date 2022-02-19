package main

import (
	"bbs-back/base"
	"bbs-back/base/handler"
	"bbs-back/models"
	_ "bbs-back/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	Init()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	// 注册过滤器
	beego.InsertFilter("/v1/*", beego.BeforeRouter, handler.Filter)
	//// 跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	// 错误页面显示，默认true
	//beego.BConfig.EnableErrorsRender = false

	// 修改系统错误输出
	beego.BConfig.RecoverFunc = handler.RecoverPanic
	beego.Run()
}

// Init 统一基础初始化
func Init() {
	base.Init()
	models.Init()
}
