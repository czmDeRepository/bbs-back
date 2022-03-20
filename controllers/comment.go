package controllers

import (
	"bbs-back/base/common"
	"bbs-back/base/dto"
	"bbs-back/models/dao"
)

type CommentController struct {
	BaseController
}

// @Title Get
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/:id [get]
func (controller *CommentController) Get() {
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.paramError(err)
		return
	}
	comment := new(dao.Comment)
	comment.Id = id
	res, err := comment.Read()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.SuccessWithData(res))
}

// @Title GetAll
// @Param 	articleId		query	int64		false	"主论贴id"
// @Param 	userId			query 	int64		false	"评论者id"
// @Param 	repliedUserId	query 	int64		false	"被回复者id -1 表示评论论贴"
// @Param 	supportCount	query 	int64		false	"点赞数量"
// @Param	status			query 	int32		false	"状态 1：正常 -1：已删除"
// @Param	pageNum			query	int32		false
// @Param	pageSize		query	int32		false
// @Param	isDetail		query	bool		false	"是否包含子回复"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [get]
func (controller *CommentController) GetAll() {
	commentDto := new(dto.CommentDto)
	err := controller.ParseForm(commentDto)
	if err != nil {
		controller.paramError(err)
		return
	}
	desc, err := controller.GetBool("desc", true)
	list, err := commentDto.Comment.Find(&commentDto.Page, desc)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	isDetail, _ := controller.GetBool("isDetail", true)
	// 查询子回复
	if isDetail {
		m := new(dao.Comment)
		notPage := common.NewNotPage()
		for _, item := range list {
			m.CommentId = item.Id
			item.ChildrenList, _ = m.Find(notPage, true)
		}
	}
	total, _ := commentDto.Comment.Count()
	res := new(common.PageDto)
	res.PageNum = commentDto.GetPageNum()
	res.PageSize = int32(len(list))
	res.Data = list
	res.Total = total
	controller.end(common.SuccessWithData(res))
}

// @Title Put
// @Param	id				int64	body	true	"评论id"
// @Param 	supportCount	int64	body	false	"点赞数量"
// @Param	status			int32	body	false	"状态 1：正常 -1：已删除"
// @Param	content			string	body	true	"内容"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [put]
func (controller *CommentController) Put() {
	curUserId := controller.getCurUserId()
	comment := new(dao.Comment)
	err := controller.ParseForm(comment)
	if err != nil {
		controller.paramError(err)
		return
	}
	if comment.Id == 0 || comment.Content == "" {
		controller.end(common.ErrorDetail(common.ERROR_PARAM, "id和评论内容都不允许为空"))
		return
	}
	// start 只允许修改自己评论
	param := new(dao.Comment)
	param.Id = comment.Id
	param.UserId = curUserId
	count, err := param.Count()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	if count == 0 {
		controller.end(common.ErrorDetail(common.ERROR_CURRENT_USER, "无法修改他人评论！"))
		return
	}
	// end
	comment.UserId = curUserId
	err = comment.Update()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Post
// @Param 	articleId		int64	body	false	"主论贴id"
// @Param 	repliedUserId	int64	body	false	"被回复者id -1 表示评论论贴"
// @Param 	supportCount	int64	body	false	"点赞数量"
// @Param	status			int32	body	false	"状态 1：正常 -1：已删除"
// @Param	content			string	body	true	"内容"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [post]
func (controller *CommentController) Post() {
	curUserId := controller.getCurUserId()

	comment := new(dao.Comment)
	err := controller.ParseForm(comment)
	if err != nil {
		controller.paramError(err)
		return
	}
	if comment.Content == "" || comment.ArticleId == 0 {
		controller.end(common.ErrorDetail(common.ERROR_PARAM, "评论内容或论贴id不允许为空"))
		return
	}
	comment.UserId = curUserId
	id, err := comment.Insert()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	// 自己回复自己，不触发消息
	if curUserId == comment.RepliedUserId {
		controller.end(common.Success())
		return
	}
	// 保持消息
	article := new(dao.Article)
	article.Id = comment.ArticleId
	article, err = article.Read()

	message := &dao.Message{}
	message.CommentId = id
	message.CreatorId = curUserId
	if comment.RepliedUserId != -1 {
		// 自己回复自己，不触发消息
		if comment.RepliedUserId == curUserId {
			controller.end(common.Success())
			return
		}
		message.TargetId = comment.RepliedUserId
		message.Kind = 2
	} else {
		// 自己评论自己论坛，不触发消息
		if article.UserId == curUserId {
			controller.end(common.Success())
			return
		}
		message.TargetId = article.UserId
		message.Kind = 1
	}
	err = message.Insert(controller.context())
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Delete
// @Param	id					int64	body	true	"消息id"
// @Param 	articleId			int64	body	true	"主论贴id"
// @Param 	repliedUserId		int64	body	false	"被回复者id"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [delete]
func (controller *CommentController) Delete() {
	curUserId := controller.getCurUserId()
	curUserRole := controller.getCurUserRole()

	param := new(dao.Comment)
	controller.ParseForm(param)
	if param.Id == 0 {
		controller.paramError(common.NewError("id非空！！！"))
		return
	}
	comment := new(dao.Comment)
	comment.Id = param.Id
	// 回复消息须带上
	comment.RepliedUserId = param.RepliedUserId
	if curUserRole != dao.USER_ROLE_ADMIN && curUserRole != dao.USER_ROLE_SUPER_ADMIN {
		// start 普通用户只允许删除自己评论
		comment.UserId = curUserId
		count, err := comment.Count()
		if err != nil {
			controller.end(common.HandleError(err))
			return
		}
		if count == 0 {
			controller.end(common.ErrorDetail(common.ERROR_CURRENT_USER, common.ERROR_MESSAGE[common.ERROR_CURRENT_USER]))
			return
		}
		// end
	}
	err := comment.Delete()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Get
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/count/:articleId [get]
func (controller *CommentController) GetCountByArticleId() {
	articleId, err := controller.GetInt64(":articleId")
	if err != nil {
		controller.paramError(err)
		return
	}
	comment := new(dao.Comment)
	comment.ArticleId = articleId
	count, err := comment.Count()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.SuccessWithData(count))
}
