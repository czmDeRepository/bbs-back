package monitor

import (
	"reflect"
	"runtime"

	"bbs-back/base/dto/information"

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

// GetResult 获取监控入口
func GetResult() (map[string]interface{}, error) {
	// user
	userGender, err := getUserGender()
	if err != nil {
		return nil, err
	}
	userIncrease, err := getUserIncrease()
	if err != nil {
		return nil, err
	}
	userYears, err := GetUserYears()
	if err != nil {
		return nil, err
	}
	// article
	activeVisitors, err := getActiveVisitors()
	if err != nil {
		return nil, err
	}
	articleIncrease, err := getArticleIncrease()
	if err != nil {
		return nil, err
	}
	commentIncrease, err := getCommentIncrease()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		USER_GENDER:                    userGender,
		USER_INCREASE_NUM:              userIncrease,
		USER_YEARS:                     userYears,
		information.ACTIVE_VISITOR_NUM: activeVisitors,
		ARTICLE_INCREASE_NUM:           articleIncrease,
		COMMENT_INCREASE_NUM:           commentIncrease,
	}, nil
}
