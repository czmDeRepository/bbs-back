package main

import (
	"net/http"
	"time"

	"bbs-back/base"
	"bbs-back/base/handler"
	"bbs-back/models"
	"bbs-back/routers"

	"github.com/beego/beego/v2/core/logs"
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
	beego.RunWithMiddleWares("", mw)
}

// Init 统一基础初始化
func Init() {
	base.Init()
	models.Init()
	routers.Init()
}

func mw(handle http.Handler) http.Handler {
	return timeLog{handle}
}

type timeLog struct {
	handle http.Handler
}

func (l timeLog) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if beego.BConfig.RunMode == "dev" {
		//logs.Info("remoteAddr:%s, method: %s, url: %s, header%v, body:%v", req.RemoteAddr, req.Method, req.URL.Path, req.Header, req.Body)
		startTime := time.Now()
		defer func() {
			delayedTime := time.Now().Sub(startTime)
			logs.Info("remoteAddr:%s, method: %s, url: %s, delayedTime: %s", req.RemoteAddr, req.Method, req.URL.Path, delayedTime)
		}()
	}
	l.handle.ServeHTTP(resp, req)
}
