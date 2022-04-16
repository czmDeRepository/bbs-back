package monitor

import (
	"fmt"
	"time"

	"bbs-back/base/storage"
)

func GetDateKey(suffix string, times ...time.Time) string {
	var now time.Time
	if len(times) > 0 {
		now = times[0]
	} else {
		now = time.Now()
	}
	year, month, day := now.Date()
	return fmt.Sprintf("%d-%d-%d_%s", year, month, day, suffix)
}

// 获取群聊统计数据
func getChatOnlineNum(days ...int) ([]Result, error) {
	var chatOnlineNums []Result
	advanceDates := MONITOR_DAYS
	if len(days) > 0 && days[0] > 0 {
		advanceDates = days[0]
	}
	for i := advanceDates - 1; i >= 0; i-- {
		current := time.Now().Add(-time.Duration(24*i) * time.Hour)
		year, month, day := current.Date()
		maxOnlineNum, err := storage.GetRedisPool().BitCount(GetDateKey(CHART_MAX_ONLINE_NUM, current))
		if err != nil {
			return nil, err
		}
		chatOnlineNums = append(chatOnlineNums, Result{
			Name:  fmt.Sprintf("%d-%d-%d", year, month, day),
			Value: maxOnlineNum,
		})
	}
	return chatOnlineNums, nil
}

// 删除30天前数据
func clearChat() error {
	err := storage.GetRedisPool().Del(GetDateKey(CHART_MAX_ONLINE_NUM, time.Now().Add(-(1+MONITOR_DAYS)*24*time.Hour)))
	if err != nil {
		return err
	}
	return storage.GetRedisPool().Del(GetDateKey(CHAT_MESSAGEG_NUM, time.Now().Add(-(1+MONITOR_DAYS)*24*time.Hour)))
}

// 获取群聊消息数
func getChatMessageNum(days ...int) ([]Result, error) {
	advanceDates := MONITOR_DAYS
	if len(days) > 0 && days[0] > 0 {
		advanceDates = days[0]
	}
	var chatOnlineNums []Result
	for i := advanceDates - 1; i >= 0; i-- {
		current := time.Now().Add(-time.Duration(24*i) * time.Hour)
		year, month, day := current.Date()
		maxOnlineNum, err := storage.GetRedisPool().GetInt64(GetDateKey(CHAT_MESSAGEG_NUM, current))
		if err != nil {
			return nil, err
		}
		chatOnlineNums = append(chatOnlineNums, Result{
			Name:  fmt.Sprintf("%d-%d-%d", year, month, day),
			Value: maxOnlineNum,
		})
	}
	return chatOnlineNums, nil
}
