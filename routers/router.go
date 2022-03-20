// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"bbs-back/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func Init() {

	ns := beego.NewNamespace("/v1",
		beego.NSInclude(
			&controllers.PublicController{},
		),
		beego.NSNamespace("/article",
			beego.NSInclude(
				&controllers.ArticleController{},
			),
		),
		beego.NSNamespace("/category",
			beego.NSInclude(
				&controllers.CategoryController{},
			),
		),
		beego.NSNamespace("/comment",
			beego.NSInclude(
				&controllers.CommentController{},
			),
		),
		beego.NSNamespace("/label",
			beego.NSInclude(
				&controllers.LabelController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&controllers.WebSocketController{},
			),
		),
		beego.NSNamespace("/message",
			beego.NSInclude(
				&controllers.MessageController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
