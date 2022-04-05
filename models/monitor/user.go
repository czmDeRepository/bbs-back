package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bbs-back/base/dto/information"

	"github.com/beego/beego/v2/core/logs"

	"bbs-back/base/storage"
	"bbs-back/models/dao"
)

// 设置用户性别比
func userGender() (err error) {
	user := new(dao.User)
	user.Status = 1
	user.Gender = "男"
	boyCount, err := user.Count()
	if err != nil {
		logs.Error("count boy user fail: %s", err.Error())
		return
	}
	user.Gender = "女"
	girlCount, err := user.Count()
	if err != nil {
		logs.Error("count girl user fail: %s", err.Error())
		return
	}
	storage.GetRedisPool().SetJson(USER_GENDER, []Result{
		{
			Name:  "男",
			Value: boyCount,
		},
		{
			Name:  "女",
			Value: girlCount,
		},
	})
	return
}

func getUserGender() ([]Result, error) {
	var res []Result
	err := storage.GetRedisPool().GetJson(USER_GENDER, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Result struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// 设置用户年代占比
func userYears() error {
	year := time.Now().Year()
	// 取整当前年代
	years := year / 10 * 10
	user := new(dao.User)
	var userMonitorList []Result
	// 统计至80年代
	for ; years >= 1980; years -= 10 {
		qs := user.GetQS(context.Background())
		location, _ := time.LoadLocation("Asia/Shanghai")
		qs = qs.Filter("birthday__gte", time.Date(years, 1, 1, 0, 0, 0, 0, location))
		qs = qs.Filter("birthday__lt", time.Date(years+10, 1, 1, 0, 0, 0, 0, location))
		count, err := qs.Count()
		if err != nil {
			logs.Error("count by birthday fail: %s", err.Error())
			return err
		}
		userMonitorList = append(userMonitorList, Result{
			Name:  fmt.Sprintf("%d0年代", (years%100)/10),
			Value: count,
		})
	}
	return storage.GetRedisPool().SetJson(USER_YEARS, userMonitorList)
}

func GetUserYears() ([]Result, error) {
	var userMonitorList []Result
	userMonitorListStr, err := storage.GetRedisPool().Get(USER_YEARS)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(userMonitorListStr), &userMonitorList)
	if err != nil {
		logs.Error("Unmarshal userYears fail userMonitorListStr: %s, error: %s", userMonitorListStr, err.Error())
		return nil, err
	}
	return userMonitorList, nil
}

// 用户增长
func userIncrease() error {
	var increaseList []Result
	//err := storage.GetRedisPool().GetJson(USER_INCREASE, &increaseList)
	//if err != nil {
	//	logs.Error("[task] get USER_INCREASE fail: %s", err.Error())
	//	return err
	//}
	for i := MONITOR_DAYS - 1; i >= 0; i-- {
		// 拿到上i天日期
		lastIDay := time.Now().Add(-time.Duration(i*24) * time.Hour)
		year, month, day := lastIDay.Date()
		nameKey := fmt.Sprintf("%d-%d-%d", year, month, day)
		qs := new(dao.User).GetQS(context.Background())
		location, _ := time.LoadLocation("Asia/Shanghai")
		qs = qs.Filter("create_time__gte", time.Date(year, month, day, 0, 0, 0, 0, location))
		// 前i-1天
		nextDay := lastIDay.Add(24 * time.Hour)
		year, month, day = nextDay.Date()
		qs = qs.Filter("create_time__lt", time.Date(year, month, day, 0, 0, 0, 0, location))
		increaseUserNum, err := qs.Count()
		if err != nil {
			logs.Error("count by create_time fail: %s", err.Error())
			return err
		}
		increaseList = append(increaseList, Result{
			Name:  nameKey,
			Value: increaseUserNum,
		})
	}
	err := storage.GetRedisPool().SetJson(USER_INCREASE_NUM, increaseList)
	return err
}

func getUserIncrease() ([]Result, error) {
	var increaseList []Result
	err := storage.GetRedisPool().GetJson(USER_INCREASE_NUM, &increaseList)
	if err != nil {
		return nil, err
	}
	return increaseList, nil
}

func clearActiveVisitor() error {
	return storage.GetRedisPool().Del(information.GetActiveVisitorKey(time.Now().Add(-(1 + MONITOR_DAYS) * 24 * time.Hour)))
}

func getActiveVisitors() ([]Result, error) {
	var activeVisitors []Result
	for i := MONITOR_DAYS - 1; i >= 0; i-- {
		current := time.Now().Add(-time.Duration(24*i) * time.Hour)
		year, month, day := current.Date()
		active, err := storage.GetRedisPool().PFCOUNT(information.GetActiveVisitorKey(current))
		if err != nil {
			return nil, err
		}
		activeVisitors = append(activeVisitors, Result{
			Name:  fmt.Sprintf("%d-%d-%d", year, month, day),
			Value: active,
		})
	}
	return activeVisitors, nil
}
