package monitor

import (
	"reflect"
	"runtime"

	"github.com/beego/beego/v2/core/logs"
	"github.com/robfig/cron"
)

func Init() {
	//monitor()
	task := cron.New()
	// 1点开始两小时一次
	task.AddFunc("0 0 1/2 * * ?", ExecMonitor)
	ExecMonitor()
	task.Start()
}

func ExecMonitor() {
	// 统计用户性别占比
	if userGender() != nil {
		return
	}
	// 统计用户年代
	if isFail(userYears) {
		return
	}
	// 用户增长
	if isFail(userIncrease) {
		return
	}
	// 清除过期日活
	if isFail(clearActiveVisitor) {
		return
	}
	// 论贴新增
	if isFail(articleIncrease) {
		return
	}
	// 评论新增
	if isFail(commentIncrease) {
		return
	}
}

func isFail(function func() error) bool {
	if err := function(); err != nil {
		logs.Error("[task] %s error: %s", runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name(), err.Error())
		return true
	}
	return false
}
