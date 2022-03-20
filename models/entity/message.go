package entity

import (
	"bbs-back/base/common"
)

type Message struct {
	Id int64 `json:"id" form:"id" orm:"pk"`
	// 状态 -1 已删除， 1 正常
	Status int `json:"status" form:"status"`
	// 评论id
	CommentId int64 `json:"comment_id" form:"commentId"`
	common.TimeModel
	// 消息类型 1 未读 2 已读
	Type int `json:"type" form:"type"`
	// 类别 1 评论消息 2 回复消息
	Kind int `json:"kind" form:"kind"`
	// 创建者id
	CreatorId int64 `json:"creator_id" form:"creatorId"`
	// 用户id 目标用户id
	TargetId int64 `json:"target_id" form:"targetId"`
}

//// 系统消息
//type SystemMessage struct {
//	Id int64 `json:"id"`
//	// 状态 -1 已删除， 1 正常
//	Status int `json:"status"`
//	// 消息内容
//	Content string `json:"content"`
//	common.TimeModel
//	// 创建者id
//	CreatorId int64 `json:"creator_id"`
//	// 生效时间范围
//	StartTime time.Time `json:"start_time"`
//	EndTime   time.Time `json:"end_time"`
//}
