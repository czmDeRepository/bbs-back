package controllers

import (
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
	role := controller.getCurUserRole()
	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		controller.end(common.ErrorMeWithCode(common.ERROR_MESSAGE[common.ERROR_POWER], common.ERROR_POWER))
		return
	}
	result, err := monitor.GetResult()
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
	result, err := monitor.GetResult()
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.SuccessWithData(result))
}
