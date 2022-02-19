package controllers

import (
	"bbs-back/base/common"
	"bbs-back/base/dto"
	"bbs-back/models/dao"
)

type CategoryController struct {
	BaseController
}

// @Title Get
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/:id [get]
func (controller *CategoryController) Get() {
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.end(common.ErrorWithMe("参数id解析错误"))
		return
	}
	category := new(dao.Category)
	category.Id = id
	res, err := category.Read()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.SuccessWithData(res))
}

// @Title GetAll
// @Param 	name	query	string	false	"category name"
// @Param	pageNum	query	int32	false
// @Param	pageSize	query	int32	false
// @Param	createTime	query	Time	false
// @Param	updateTime	query	Time	false
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [get]
func (controller *CategoryController) GetAll() {
	c := new(dto.CategoryDto)
	err := controller.ParseForm(c)
	if err != nil {
		controller.paramError(common.NewError(err.Error()))
		return
	}
	isDesc, _ := controller.GetBool("isDesc", false)
	orderIndex, _ := controller.GetInt("orderIndex", 0)
	list, err := c.Category.Find(&c.Page, orderIndex, isDesc)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	total, _ := c.Category.Count()
	res := new(common.PageDto)
	res.Total = total
	res.Data = list
	res.PageNum = c.GetPageNum()
	res.PageSize = int32(len(list))
	controller.end(common.SuccessWithData(res))
}

// @Title Put
// @Description 更新主题信息
// @Param 	name	body	string	true	"name"
// @Param	createTime	body	string	false	"创建时间"
// @Param	updateTime	body	string	false	"更新时间"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [put]
func (controller *CategoryController) Put() {
	role := controller.getCurUserRole()
	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		controller.end(common.ErrorDetail(common.ERROR_POWER, common.ERROR_MESSAGE[common.ERROR_POWER]))
		return
	}
	category := new(dao.Category)
	err := controller.ParseForm(category)
	if err != nil {
		controller.paramError(common.NewError(err.Error()))
		return
	}
	err = category.Update()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Post
// @Param 	name	body	string	true	"name"
// @Success 0 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/ [post]
func (controller *CategoryController) Post() {
	role := controller.getCurUserRole()

	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		controller.end(common.ErrorDetail(common.ERROR_POWER, common.ERROR_MESSAGE[common.ERROR_POWER]))
		return
	}
	category := new(dao.Category)
	controller.ParseForm(category)
	err := category.Insert()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Delete
// @Success 200 {object} dto.Result
// @Failure 1000 :参数错误
// @router	/:id [delete]
func (controller *CategoryController) Delete() {
	role := controller.getCurUserRole()
	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		controller.end(common.ErrorDetail(common.ERROR_POWER, common.ERROR_MESSAGE[common.ERROR_POWER]))
		return
	}
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.paramError(common.NewError(err.Error()))
		return
	}
	category := new(dao.Category)
	category.Id = id
	_, err = category.Delete()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}
