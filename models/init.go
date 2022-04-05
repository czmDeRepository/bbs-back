package models

import (
	"bbs-back/base/dto/information"
	"bbs-back/base/storage"
	"bbs-back/models/dao"
	"bbs-back/models/monitor"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func Init() {
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	orm.RegisterModel(new(dao.Article), new(dao.Category), new(dao.Comment), new(dao.Label), new(dao.User), new(dao.Message))
	dao.ORM = orm.NewOrm()
	article := new(dao.Article)
	totalArticleNum, _ := article.Count(nil)
	if err := storage.GetRedisPool().Set(information.TOTAL_ARTICLE_NUM, totalArticleNum); err != nil {
		logs.Error("set totalArticleNum fail, error: %s", err.Error())
	}
	totalReadNum, _ := article.TotalReadNum()
	if err := storage.GetRedisPool().Set(information.TOTAL_READ_NUM, totalReadNum); err != nil {
		logs.Error("set totalReadNum fail, error: %s", err.Error())
	}
	// 监控大盘信息
	monitor.Init()
}
