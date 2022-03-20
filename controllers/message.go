package controllers

import (
	"encoding/json"

	"bbs-back/base/common"
	"bbs-back/base/dto"
	"bbs-back/models/dao"

	beego "github.com/beego/beego/v2/server/web"
)

type MessageController struct {
	BaseController
}

// @Title Get
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/:id [get]
func (controller *MessageController) Get() {
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.end(common.ErrorWithMe("id不为空"))
		return
	}
	controller.mustLogin()
	message := new(dao.Message)
	message.Id = id
	res, err := message.Read(controller.context())
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.SuccessWithData(res))
}

// @Title GetStatistics
// @Description 未读获取消息数量
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/statistics [get]
func (controller *MessageController) GetStatistics() {

	curUserId := controller.getCurUserId()
	message := new(dao.Message)
	message.TargetId = curUserId
	message.Status = 1
	// 未读
	message.Type = 1
	// 评论数量
	message.Kind = 1
	commentCount, err := message.Count(controller.context())
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	// 回复数量
	message.Kind = 2
	replyCount, err := message.Count(controller.context())
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}

	controller.end(common.SuccessWithData(beego.M{
		"commentCount": commentCount,
		"replyCount":   replyCount,
	}))
}

// @Title GetAll
// @Param 	status	query	int	 false
// @Param	pageNum	query	int32	false
// @Param	pageSize	query	int32	false
// @Param	createTime	query	Time	false
// @Param	updateTime	query	Time	false
// @Param	creatorId	query	int64	false
// @Param	type	query	int32	false
// @Param	targetId	query	int32	false
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [get]
func (controller *MessageController) GetAll() {
	messageDto := new(dto.MessageDto)
	err := controller.ParseForm(messageDto)
	if err != nil {
		controller.paramError(common.NewError(err.Error()))
		return
	}
	list, err := messageDto.Message.Find(controller.context(), &messageDto.Page)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	total, _ := messageDto.Message.Count(controller.context())
	res := new(common.PageDto)
	res.Total = total
	res.Data = list
	res.PageNum = messageDto.GetPageNum()
	res.PageSize = int32(len(list))
	controller.end(common.SuccessWithData(res))
}

// @Title   Post
// @Param 	status	body	int	 false
// @Param	pageNum	body	int32	false
// @Param	pageSize	body	int32	false
// @Param	creatorId	body	int64	false
// @Param	type	body	int32	false
// @Param	targetId	body	int32	false
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [post]
func (controller *MessageController) Post() {
	message := new(dao.Message)
	err := controller.ParseForm(message)
	if err != nil {
		controller.paramError(common.NewError(err.Error()))
		return
	}
	curUserId := controller.getCurUserId()
	message.CreatorId = curUserId
	err = message.Insert(controller.context())
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title   Put
// @Param 	ids	body	int[]	 false
// @Param 	status	body	int	 false
// @Param	type	body	int	false
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [put]
func (controller *MessageController) Put() {
	message := new(dao.Message)
	err := controller.ParseForm(message)
	if err != nil {
		controller.paramError(err)
		return
	}
	var ids []int64
	err = json.Unmarshal([]byte(controller.GetString("ids")), &ids)
	if err != nil {
		controller.paramError(err)
		return
	}
	// 保证登录状态
	controller.mustLogin()
	for _, id := range ids {
		message.Id = id
		err = message.Update(controller.context())
		if err != nil {
			controller.end(common.HandleError(err))
			return
		}
	}
	controller.end(common.Success())
}

type messageParam struct {
	Ids []int64 `json:"id" form:"ids"`
	// 状态 -1 已删除， 1 正常
	Status int `json:"status" form:"status"`
	// 消息类型 1 未读 2 已读
	Type int `json:"type" form:"type"`
}

// @Title Delete
// @Success 200 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/:id [delete]
func (controller *MessageController) Delete() {
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.paramError(common.NewError(err.Error()))
		return
	}
	message := new(dao.Message)
	message.Id = id
	err = message.Delete(controller.context())
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}
