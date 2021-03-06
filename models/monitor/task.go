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

	// 清除过期群聊日活
	if isFail(clearChat) {
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
	userIncrease, err := getUserIncrease(days)
	if err != nil {
		return nil, err
	}
	activeVisitors, err := getActiveVisitors(days)
	if err != nil {
		return nil, err
	}
	// article
	articleIncrease, err := getArticleIncrease(days)
	if err != nil {
		return nil, err
	}
	commentIncrease, err := getCommentIncrease(days)
	if err != nil {
		return nil, err
	}

	chatOnlineNum, err := getChatOnlineNum(days)
	if err != nil {
		return nil, err
	}
	chatMessageNums, err := getChatMessageNum(days)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		USER_GENDER:                    userGender,
		USER_YEARS:                     userYears,
		USER_INCREASE_NUM:              userIncrease,
		information.ACTIVE_VISITOR_NUM: activeVisitors,
		ARTICLE_INCREASE_NUM:           articleIncrease,
		COMMENT_INCREASE_NUM:           commentIncrease,
		CHART_MAX_ONLINE_NUM:           chatOnlineNum,
		CHAT_MESSAGEG_NUM:              chatMessageNums,
	}, nil
}

func GetUser(days int) (map[string]interface{}, error) {
	userIncrease, err := getUserIncrease(days)
	if err != nil {
		return nil, err
	}
	activeVisitors, err := getActiveVisitors(days)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		USER_INCREASE_NUM:              userIncrease,
		information.ACTIVE_VISITOR_NUM: activeVisitors,
	}, nil
}

func GetArticle(days int) (map[string]interface{}, error) {
	// article
	articleIncrease, err := getArticleIncrease(days)
	if err != nil {
		return nil, err
	}
	commentIncrease, err := getCommentIncrease(days)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		ARTICLE_INCREASE_NUM: articleIncrease,
		COMMENT_INCREASE_NUM: commentIncrease,
	}, nil
}

func GetChat(days int) (map[string]interface{}, error) {
	chatOnlineNums, err := getChatOnlineNum(days)
	if err != nil {
		return nil, err
	}
	chatMessageNums, err := getChatMessageNum(days)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		CHART_MAX_ONLINE_NUM: chatOnlineNums,
		CHAT_MESSAGEG_NUM:    chatMessageNums,
	}, nil
}
