package information

import (
	"fmt"
	"time"
)

type Information struct {
	ActiveVisitorNum   int64 `json:"activeVisitorNum"`
	TotalArticleNum    int64 `json:"totalArticleNum"`
	TotalChatPersonNum int64 `json:"totalChatPersonNum"`
	TotalReadNum       int64 `json:"totalReadNum"`
}

const (
	ACTIVE_VISITOR_NUM = "activeVisitorNum"
	TOTAL_ARTICLE_NUM  = "totalArticleNum"
	TOTAL_READ_NUM     = "totalReadNum"
)

func GetActiveVisitorKey(times ...time.Time) string {
	var now time.Time
	if len(times) > 0 {
		now = times[0]
	} else {
		now = time.Now()
	}
	year, month, day := now.Date()
	return fmt.Sprintf("%d-%d-%d_%s", year, month, day, ACTIVE_VISITOR_NUM)
}
