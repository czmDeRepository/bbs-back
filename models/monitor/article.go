package monitor

import (
	"context"
	"fmt"
	"time"

	"bbs-back/base/storage"
	"bbs-back/models/dao"

	"github.com/beego/beego/v2/core/logs"
)

func articleIncrease() error {
	var increaseList []Result
	for i := MONITOR_DAYS - 1; i >= 0; i-- {
		// 拿到前i天日期
		lastIDay := time.Now().Add(-time.Duration(i*24) * time.Hour)
		year, month, day := lastIDay.Date()
		nameKey := fmt.Sprintf("%d-%d-%d", year, month, day)
		qs := new(dao.Article).GetQS(context.Background())
		// 只统计已经发布的
		qs.Filter("status", 2)
		location, _ := time.LoadLocation("Asia/Shanghai")
		qs = qs.Filter("create_time__gte", time.Date(year, month, day, 0, 0, 0, 0, location))
		// 前i-1天
		nextDay := lastIDay.Add(24 * time.Hour)
		year, month, day = nextDay.Date()
		qs = qs.Filter("create_time__lt", time.Date(year, month, day, 0, 0, 0, 0, location))
		increaseArticleNum, err := qs.Count()
		if err != nil {
			logs.Error("count by create_time fail: %s", err.Error())
			return err
		}
		increaseList = append(increaseList, Result{
			Name:  nameKey,
			Value: increaseArticleNum,
		})
	}
	return storage.GetRedisPool().SetJson(ARTICLE_INCREASE_NUM, increaseList)
}

func getArticleIncrease(days ...int) ([]Result, error) {
	var increaseList []Result
	err := storage.GetRedisPool().GetJson(ARTICLE_INCREASE_NUM, &increaseList)
	if err != nil {
		return nil, err
	}
	if len(days) > 0 && days[0] > 0 && len(increaseList) > MONITOR_DAYS-days[0] {
		increaseList = increaseList[MONITOR_DAYS-days[0]:]
	}
	return increaseList, nil
}
