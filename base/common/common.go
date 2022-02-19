package common

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type TimeModel struct {
	CreateTime DateTime `json:"createTime" form:"createTime" orm:"auto_now_add;type(datetime);"`
	UpdateTime DateTime `json:"updateTime" form:"updateTime" orm:"auto_now;type(datetime)"`
}

func (t *TimeModel) TimeFilter(qs orm.QuerySeter) orm.QuerySeter {
	if !t.CreateTime.IsZero() {
		qs = qs.Filter("create_time", t.CreateTime.Time)
	}
	if !t.UpdateTime.IsZero() {
		qs = qs.Filter("update_time", t.UpdateTime.Time)
	}
	return qs
}

type PageDto struct {
	Data interface{} `json:"data"`
	Page
}

type Page struct {
	PageNum  int32 `json:"pageNum" form:"pageNum"`
	PageSize int32 `json:"pageSize" form:"pageSize"`
	Total    int64 `json:"total"`
}

func NewNotPage() *Page {
	return &Page{PageNum: -1, PageSize: -1}
}
func (page *Page) GetPageSize() int32 {
	if page.PageSize <= 0 {
		return 10
	} else if page.PageSize > 1000 {
		// 最大限制1000
		return 1000
	}
	return page.PageSize
}

func (page *Page) GetPageNum() int32 {
	if page.PageNum <= 0 {
		return 1
	}
	return page.PageNum
}

func (page *Page) GetOffset() int32 {
	return (page.GetPageNum() - 1) * page.GetPageSize()
}
func (page *Page) NeedPage() bool {
	return page.PageSize != -1 && page.PageNum != -1
}

// 时间范围
type RangeTime struct {
	StartTime time.Time `json:"startTime,omitempty" form:"startTime"`
	EndTime   time.Time `json:"endTime,omitempty" form:"endTime"`
}

func (t *RangeTime) GetStartTime() string {
	return t.StartTime.Format(FormatDateTime)
}

func (t *RangeTime) GetEndTime() string {
	return t.EndTime.Format(FormatDateTime)
}
