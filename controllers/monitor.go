package controllers

import (
	"fmt"

	"bbs-back/base/common"
	"bbs-back/models/dao"
	"bbs-back/models/monitor"
)

type Monitor struct {
	BaseController
}

// @Title	Get
// @router  / [get]
func (controller *Monitor) Get() {
	days, err := controller.GetInt("days")
	if err != nil {
		controller.paramError(err)
		return
	}
	role := controller.getCurUserRole()
	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		controller.end(common.ErrorMeWithCode(common.ERROR_MESSAGE[common.ERROR_POWER], common.ERROR_POWER))
		return
	}
	result, err := monitor.GetResult(days)
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.SuccessWithData(result))
}

// @Title	put
// @router  / [put]
func (controller *Monitor) Put() {
	role := controller.getCurUserRole()
	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		controller.end(common.ErrorMeWithCode(common.ERROR_MESSAGE[common.ERROR_POWER], common.ERROR_POWER))
		return
	}
	monitor.ExecMonitor()
	result, err := monitor.GetResult(7)
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.SuccessWithData(result))
}

// @Title	GetByType
// @router  /:type [get]
func (controller *Monitor) GetByType() {
	days, err := controller.GetInt("days")
	if err != nil {
		controller.paramError(err)
		return
	}
	role := controller.getCurUserRole()
	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		controller.end(common.ErrorMeWithCode(common.ERROR_MESSAGE[common.ERROR_POWER], common.ERROR_POWER))
		return
	}
	monitorType := controller.GetString(":type")
	var result map[string]interface{}
	switch monitorType {
	case "article":
		result, err = monitor.GetArticle(days)
	case "user":
		result, err = monitor.GetUser(days)
	case "chat":
		result, err = monitor.GetChat(days)
	default:
		controller.end(common.ErrorWithMe(fmt.Sprintf("not type【%s】 exist", monitorType)))
		return
	}
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.SuccessWithData(result))
}
