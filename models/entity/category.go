package entity

import "bbs-back/base/common"

type Category struct {
	Id   int64  `json:"id" form:"id" orm:"pk"`
	Name string `json:"name" form:"name"`
	common.TimeModel
}
