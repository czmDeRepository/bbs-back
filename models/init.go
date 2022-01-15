package models

import (
	"bbs-back/base/database/bbsRedis"
	"bbs-back/base/dto/information"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

var ORM orm.Ormer

func init() {
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	orm.RegisterModel(new(Article), new(Category), new(Comment), new(Label), new(User))
	ORM = orm.NewOrm()
	totalUserNum, _ := new(User).Count()
	bbsRedis.Set(information.TOTAL_USER_NUM, totalUserNum)
	article := new(Article)
	totalArticleNum, _ := article.Count(nil)
	bbsRedis.Set(information.TOTAL_ARTICLE_NUM, totalArticleNum)
	totalReadNum, _ := article.TotalReadNum()
	bbsRedis.Set(information.TOTAL_READ_NUM, totalReadNum)
}
