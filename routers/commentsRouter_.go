package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Follow",
            Router: "/follow",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "UnFollow",
            Router: "/unfollow",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CommentController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CommentController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CommentController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CommentController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CommentController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CommentController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CommentController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CommentController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CommentController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CommentController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:CommentController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:CommentController"],
        beego.ControllerComments{
            Method: "GetCountByArticleId",
            Router: "/count/:articleId",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:LabelController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:LabelController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:LabelController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:LabelController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:LabelController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:LabelController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:LabelController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:LabelController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:PublicController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:PublicController"],
        beego.ControllerComments{
            Method: "Information",
            Router: "/information",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:PublicController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:PublicController"],
        beego.ControllerComments{
            Method: "Upload",
            Router: "/upload",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "FollowGet",
            Router: "/follow",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "FollowPost",
            Router: "/follow",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "FollowDelete",
            Router: "/follow",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: "/login",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:UserController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:UserController"],
        beego.ControllerComments{
            Method: "Refresh",
            Router: "/refresh",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bbs-back/controllers:WebSocketController"] = append(beego.GlobalControllerRouter["bbs-back/controllers:WebSocketController"],
        beego.ControllerComments{
            Method: "Join",
            Router: "/chat",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
