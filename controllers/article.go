package controllers

import (
	"strings"
	"time"

	"bbs-back/base/common"
	"bbs-back/base/dto"
	"bbs-back/models/dao"
)

type ArticleController struct {
	BaseController
}

// @Title Get
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/:id [get]
func (controller *ArticleController) Get() {
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.end(common.ErrorWithMe("id不为空"))
		return
	}
	article := new(dao.Article)
	article.Id = id
	res, err := article.Read()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	userId, err := controller.GetInt64("userId")
	if err == nil {
		res.FollowingFlag = res.FollowingCount(userId) > 0
	}
	res.FollowCount = res.FollowingCount()
	controller.end(common.SuccessWithData(res))
}

// @Title GetAll
// @Param	title			query	string	false  	"论贴标题"
// @Param	categoryId		query	int64	false	"论贴分类"
// @Param	userId			query	int64	false	"创建者id"
// @Param	status			query	int32	false	"状态 -1 已删除 1-未发布，2-已经发布"
// @Param	createTime		query	Time	false	"创建时间"
// @Param	updateTime		query	Time	false	"更新时间"
// @Param	pageNum			query	int32	false	"页数"
// @Param	pageSize		query	int32	false	"页面大小"
// @param	labelIdList		query	[]int64 false	"关联标签id查询"
// @Param	isDesc			query	bool	false
// @Param	orderIndex		query	int32	false	"0-更新时间 1-创建时间 2-阅读数量 3-点赞数量"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [get]
func (controller *ArticleController) GetAll() {
	article := new(dto.ArticleDto)
	controller.ParseForm(article)
	labelIdList := controller.GetString("labelIdList")
	var labelIdListParam []string
	if labelIdList != "" {
		labelIdListParam = strings.Split(labelIdList, ",")
	}
	if article.FollowingFlag {
		article.UserId = controller.getCurUserId()
	}
	list, err := article.Find(&article.Page, article.OrderIndex, article.IsAsc, &article.RangeTime, labelIdListParam...)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	for _, item := range list {
		item.FollowCount = item.FollowingCount()
	}
	total, _ := article.Article.Count(&article.RangeTime, labelIdListParam...)
	res := new(common.PageDto)
	res.PageNum = article.GetPageNum()
	res.PageSize = int32(len(list))
	res.Data = list
	res.Total = total
	controller.end(common.SuccessWithData(res))
}

// @Title Put
// @Param	title			body	string	false  	"论贴标题"
// @Param	categoryId		body	int64	false	"论贴分类"
// @Param	userId			body	int64	false	"创建者id"
//@Success 0 {object} dto.Result
//@Failure 1000 :参数错误
// @router	/ [put]
func (controller *ArticleController) Put() {
	article := new(dao.Article)
	controller.ParseForm(article)
	if article.Id == 0 {
		controller.paramError(common.NewError("id不为空"))
		return
	}
	labelStrings := controller.GetString("labelList")
	var labelList []string
	if labelStrings != "" {
		labelList = strings.Split(labelStrings, ",")
	}
	article.UpdateTime.Time = time.Now()
	err := article.Update(labelList...)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Post
//@Success 0 {object} dto.Result
//@Failure 1000 :参数错误
// @router	/ [post]
func (controller *ArticleController) Post() {
	article := new(dao.Article)
	curUserId := controller.getCurUserId()
	labelStrings := controller.GetString("labelList")
	controller.ParseForm(article)
	article.UserId = curUserId
	var labelList []string
	if labelStrings != "" {
		labelList = strings.Split(labelStrings, ",")
	}
	err := article.Insert(labelList...)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}

	controller.end(common.SuccessWithData(article))
}

// @Title Delete
// @Param 	id	query 	int64	true	"key"
//@Success 0 {object} dto.Result
//@Failure 1000 :参数错误
// @router	/ [delete]
func (controller *ArticleController) Delete() {
	id, err := controller.GetInt64("id")
	if err != nil {
		controller.end(common.ErrorWithMe("id不为空"))
		return
	}
	userId := controller.getCurUserId()
	role := controller.getCurUserRole()
	article := new(dao.Article)
	article.Id = id
	if role == dao.USER_ROLE_USER {
		article.UserId = userId
		count, _ := article.Count(nil)
		if count == 0 {
			controller.end(common.ErrorWithMe("非法删除他人帖子"))
			return
		}
	}
	err = article.Delete()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Param	articleId			body	int64	true  	"论贴id"
//@Success 0 {object} dto.Result
//@Failure 1000 :参数错误
// @router	/follow [post]
func (controller *ArticleController) Follow() {
	articleId, err := controller.GetInt64("articleId")
	if err != nil {
		controller.end(common.ErrorDetail(common.ERROR_PARAM, "论贴id不为空"))
		return
	}
	userId := controller.getCurUserId()
	article := new(dao.Article)
	article.Id = articleId
	err = article.Follow(userId)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Param	articleId			body	int64	true  	"论贴id"
//@Success 0 {object} dto.Result
//@Failure 1000 :参数错误
// @router	/unfollow [delete]
func (controller *ArticleController) UnFollow() {
	articleId, err := controller.GetInt64("articleId")
	if err != nil {
		controller.end(common.ErrorDetail(common.ERROR_PARAM, "论贴id不为空"))
		return
	}
	userId := controller.getCurUserId()
	article := new(dao.Article)
	article.Id = articleId
	err = article.UnFollow(userId)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}
