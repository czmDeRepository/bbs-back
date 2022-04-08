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
	// 两小时一次
	task.AddFunc("10 0 0/2 * * ?", ExecMonitor)
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
func GetResult(days int) (map[string]interface{}, error) {
	if days > MONITOR_DAYS {
		days = MONITOR_DAYS
	}
	if days <= 0 {
		days = 7
	}
	// user
	userGender, err := getUserGender()
	if err != nil {
		return nil, err
	}
	userYears, err := getUserYears()
	if err != nil {
		return nil, err
	}
	userIncrease, err := getUserIncrease()
	if err != nil {
		return nil, err
	}
	activeVisitors, err := getActiveVisitors()
	if err != nil {
		return nil, err
	}
	// article
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
		USER_YEARS:                     userYears,
		USER_INCREASE_NUM:              userIncrease[MONITOR_DAYS-days:],
		information.ACTIVE_VISITOR_NUM: activeVisitors[MONITOR_DAYS-days:],
		ARTICLE_INCREASE_NUM:           articleIncrease[MONITOR_DAYS-days:],
		COMMENT_INCREASE_NUM:           commentIncrease[MONITOR_DAYS-days:],
	}, nil
}

func GetUser(days int) (map[string]interface{}, error) {
	userIncrease, err := getUserIncrease()
	if err != nil {
		return nil, err
	}
	activeVisitors, err := getActiveVisitors()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		USER_INCREASE_NUM:              userIncrease[MONITOR_DAYS-days:],
		information.ACTIVE_VISITOR_NUM: activeVisitors[MONITOR_DAYS-days:],
	}, nil
}

func GetArticle(days int) (map[string]interface{}, error) {
	// article
	articleIncrease, err := getArticleIncrease()
	if err != nil {
		return nil, err
	}
	commentIncrease, err := getCommentIncrease()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		ARTICLE_INCREASE_NUM: articleIncrease[MONITOR_DAYS-days:],
		COMMENT_INCREASE_NUM: commentIncrease[MONITOR_DAYS-days:],
	}, nil
}
