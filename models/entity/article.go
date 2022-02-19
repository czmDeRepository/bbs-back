package entity

import "bbs-back/base/common"

type Article struct {
	Id            int64           `json:"id" form:"id" orm:"pk"`
	Title         string          `json:"title" form:"title"`
	CategoryId    int64           `json:"categoryId" form:"categoryId"`
	UserId        int64           `json:"userId" form:"userId"`
	Status        int32           `json:"status" form:"status"` // -1 已删除 1-未发布，2-已经发布
	ReadCount     int32           `json:"readCount" form:"readCount"`
	CreateTime    common.DateTime `json:"createTime" form:"createTime" orm:"auto_now_add;type(datetime);"`
	UpdateTime    common.DateTime `json:"updateTime" form:"updateTime" orm:"column(update_time)"`
	Content       string          `json:"content" form:"content"`
	FollowingFlag bool            `json:"followingFlag" form:"followingFlag" orm:"-"` // 是否关注
	FollowCount   int64           `json:"followCount" form:"followCount" orm:"-"`     // 关注数量
}
