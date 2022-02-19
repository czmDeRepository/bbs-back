package models

import (
	"bbs-back/base/dto/information"
	"bbs-back/base/storage"
	"bbs-back/models/dao"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func Init() {
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	orm.RegisterModel(new(dao.Article), new(dao.Category), new(dao.Comment), new(dao.Label), new(dao.User))
	dao.ORM = orm.NewOrm()
	totalUserNum, _ := new(dao.User).Count()
	storage.GetRedisPool().Set(information.TOTAL_USER_NUM, totalUserNum)
	article := new(dao.Article)
	totalArticleNum, _ := article.Count(nil)
	storage.GetRedisPool().Set(information.TOTAL_ARTICLE_NUM, totalArticleNum)
	totalReadNum, _ := article.TotalReadNum()
	storage.GetRedisPool().Set(information.TOTAL_READ_NUM, totalReadNum)
}
