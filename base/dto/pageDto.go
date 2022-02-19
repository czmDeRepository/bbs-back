package dto

import (
	"bbs-back/base/common"
	"bbs-back/models/dao"
)

type UserDto struct {
	dao.User
	common.PageDto
}

type CategoryDto struct {
	dao.Category
	common.PageDto
}

type ArticleDto struct {
	dao.Article
	common.PageDto
	IsAsc      bool  `json:"isAsc,omitempty" form:"isAsc"`
	OrderIndex int32 `json:"orderIndex,omitempty" form:"orderIndex"`
	common.RangeTime
}

type LabelDto struct {
	dao.Label
	common.PageDto
}

type CommentDto struct {
	dao.Comment
	common.PageDto
}
