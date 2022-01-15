package dto

import (
	"bbs-back/base/common"
	"bbs-back/models"
)

type UserDto struct {
	models.User
	common.PageDto
}

type CategoryDto struct {
	models.Category
	common.PageDto
}

type ArticleDto struct {
	models.Article
	common.PageDto
	IsAsc		bool	`json:"isAsc,omitempty" form:"isAsc"`
	OrderIndex	int32	`json:"orderIndex,omitempty" form:"orderIndex"`
	common.RangeTime
}

type LabelDto struct {
	models.Label
	common.PageDto
}

type CommentDto struct {
	models.Comment
	common.PageDto
}