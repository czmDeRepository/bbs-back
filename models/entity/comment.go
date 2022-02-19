package entity

import "bbs-back/base/common"

type Comment struct {
	Id int64 `json:"id" form:"id" orm:"pk"`
	// 主论贴id
	ArticleId int64 `json:"articleId" form:"articleId"`
	// 评论者id
	UserId int64 `json:"userId" form:"userId"`
	// 被回复者id
	RepliedUserId int64  `json:"repliedUserId" form:"repliedUserId"`
	Status        int32  `json:"status" form:"status"`
	Content       string `json:"content" form:"content"`
	// 父级评论id
	CommentId int64 `json:"commentId" form:"commentId"`
	common.TimeModel
}
