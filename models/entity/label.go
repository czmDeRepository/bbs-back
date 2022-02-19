package entity

import "bbs-back/base/common"

type Label struct {
	Id     int64  `json:"id" form:"id" orm:"pk"`
	Name   string `json:"name" form:"name"`
	Status int32  `json:"status" form:"status"`
	common.TimeModel
}
