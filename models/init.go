package models

import (
	"bbs-back/base/dto/information"
	"bbs-back/base/storage"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

var ORM orm.Ormer

func Init() {
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	orm.RegisterModel(new(Article), new(Category), new(Comment), new(Label), new(User))
	ORM = orm.NewOrm()
	totalUserNum, _ := new(User).Count()
	storage.GetRedisPool().Set(information.TOTAL_USER_NUM, totalUserNum)
	article := new(Article)
	totalArticleNum, _ := article.Count(nil)
	storage.GetRedisPool().Set(information.TOTAL_ARTICLE_NUM, totalArticleNum)
	totalReadNum, _ := article.TotalReadNum()
	storage.GetRedisPool().Set(information.TOTAL_READ_NUM, totalReadNum)
}
