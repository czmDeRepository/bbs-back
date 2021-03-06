package entity

import "bbs-back/base/common"

type User struct {
	Id              int64           `json:"id" form:"id" orm:"pk"`
	Name            string          `json:"name" form:"name"`
	Password        string          `json:"password" form:"password"`
	Account         string          `json:"account" form:"account"`
	Email           string          `json:"email" form:"email"`
	TelephoneNumber string          `json:"telephoneNumber" form:"telephoneNumber"`
	Birthday        common.DateTime `json:"birthday" form:"birthday"`
	Status          int32           `json:"status" form:"status"`
	Role            int32           `json:"role" form:"role"`
	Gender          string          `json:"gender" form:"gender"`
	ImageUrl        string          `json:"imageUrl" form:"imageUrl"`
	common.TimeModel
}
