package controllers

import (
	"bbs-back/base/common"
	"bbs-back/base/dto"
	"bbs-back/models"
)

type LabelController struct {
	BaseController
}

// @Title Get
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/:id [get]
func (controller *LabelController) Get()  {
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.paramError(err)
		return
	}
	label := new(models.Label)
	label.Id = id
	res, err := label.Read()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.SuccessWithData(res))
}

// @Title GetAll
// @Param 	name	string	query	false	"标签名"
// @Param	pageNum	query	int32	false
// @Param	pageSize	query	int32	false
// @Param	createTime	query	Time	false
// @Param	updateTime	query	Time	false
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [get]
func (controller *LabelController) GetAll() {
	labelDto := new(dto.LabelDto)
	err := controller.ParseForm(labelDto)
	if err != nil {
		controller.paramError(err)
		return
	}
	isDesc, _ := controller.GetBool("isDesc", false)
	orderIndex, _ := controller.GetInt("orderIndex", 0)
	list, err := labelDto.Label.Find(&labelDto.Page, orderIndex, isDesc)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	total, _ := labelDto.Label.Count()
	res := new(common.PageDto)
	res.PageNum = labelDto.GetPageNum()
	res.PageSize = int32(len(list))
	res.Data = list
	res.Total = total
	controller.end(common.SuccessWithData(res))
}

// @Title Put
// @Param 	name	string	body	true	"标签名"
// @Param	id		int64	body	true	"id"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [put]
func (controller *LabelController) Put()  {
	role := controller.getCurUserRole()

	if role != models.USER_ROLE_ADMIN && role != models.USER_ROLE_SUPER_ADMIN  {
		controller.end(common.ErrorDetail(nil, common.ERROR_POWER, common.ERROR_MESSAGE[common.ERROR_POWER]))
		return
	}

	label := new(models.Label)
	err := controller.ParseForm(label)
	if err != nil {
		controller.paramError(err)
		return
	}
	if label.Id == 0 {
		controller.end(common.ErrorDetail(nil, common.ERROR_PARAM, "id不允许为空"))
		return
	}
	err = label.Update()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Post
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [post]
func (controller *LabelController) Post()  {
	role := controller.getCurUserRole()

	if role != models.USER_ROLE_ADMIN && role != models.USER_ROLE_SUPER_ADMIN  {
		controller.end(common.ErrorDetail(nil, common.ERROR_POWER, common.ERROR_MESSAGE[common.ERROR_POWER]))
		return
	}
	label := new(models.Label)
	err := controller.ParseForm(label)
	if err != nil {
		controller.paramError(err)
		return
	}
	if label.Name == "" {
		controller.end(common.ErrorDetail(nil, common.ERROR_PARAM, "标签名都不允许为空"))
		return
	}
	err = label.Insert()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

