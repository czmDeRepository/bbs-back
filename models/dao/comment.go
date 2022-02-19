package dao

import (
	"bytes"
	"strconv"

	"bbs-back/base/common"
	"bbs-back/models/entity"
)

type Comment struct {
	entity.Comment
	// 评论者用户名
	UserName string `json:"userName"`
	// 评论者性别
	UserGender string `json:"userGender"`
	//	评论者用户头像路径
	UserImageUrl string `json:"userImageUrl"`
	// 被回复者用户名
	RepliedUserName string `json:"repliedUserName"`
	// 被回复者用户头像
	RepliedUserImageUrl string `json:"repliedUserImageUrl"`
	//RepliedUserGender			string	`json:"repliedUserGender"`
	// 子回复数组
	ChildrenList []*Comment `json:"childrenList" orm:"-"`
}

func (c *Comment) Insert() error {
	if c.Status == 0 {
		c.Status = 1
	}
	if c.RepliedUserId == 0 {
		c.RepliedUserId = -1
	}
	c.CreateTime = *common.Now()
	_, err := ORM.Raw("insert into comment (article_id, user_id, replied_user_id, create_time, update_time, content, status, comment_id) values (?, ?, ?, ?, ?, ?, ?, ?)",
		c.ArticleId, c.UserId, c.RepliedUserId, c.CreateTime.Time, c.CreateTime.Time, c.Content, c.Status, c.CommentId,
	).Exec()
	return err
}

func (c *Comment) Read() (*Comment, error) {
	var comment = new(Comment)
	err := ORM.Raw("select c.*, u.`name` user_name, u.`image_url` user_image_url, u.gender user_gender, u2.name replied_user_name, u2.`image_url` replied_user_image_url "+
		"from `comment` c inner join `user` u on c.user_id = u.id "+
		"left join `user` u2 on c.replied_user_id = u2.id  "+
		"where c.id = ?", c.Id).QueryRow(comment)
	return comment, err
}

func (c *Comment) Find(page *common.Page, desc bool) ([]*Comment, error) {
	var res []*Comment
	var sql = new(bytes.Buffer)
	sql.WriteString("select c.*, u.`name` user_name, u.`image_url` user_image_url, u.gender user_gender, u2.name replied_user_name, u2.`image_url` replied_user_image_url from " +
		"`comment` c inner join `user` u on c.user_id = u.id " +
		"left join `user` u2 on c.replied_user_id = u2.id " +
		"where c.status = ")
	if c.Status != 0 {
		sql.WriteString(strconv.FormatInt(int64(c.Status), 10))
	} else {
		sql.WriteString("1")
	}
	if c.Id != 0 {
		sql.WriteString(" and c.id = ")
		sql.WriteString(strconv.FormatInt(c.Id, 10))
	}
	if c.ArticleId != 0 {
		sql.WriteString(" and c.article_id = ")
		sql.WriteString(strconv.FormatInt(c.ArticleId, 10))
	}
	if c.UserId != 0 {
		sql.WriteString(" and c.user_id = ")
		sql.WriteString(strconv.FormatInt(c.UserId, 10))
	}
	if c.RepliedUserId != 0 {
		sql.WriteString(" and c.replied_user_id = ")
		sql.WriteString(strconv.FormatInt(c.RepliedUserId, 10))
	}
	if c.CommentId != 0 {
		sql.WriteString(" and c.comment_id = ")
		sql.WriteString(strconv.FormatInt(c.CommentId, 10))
	}
	sql.WriteString(" order by c.create_time ")
	if desc {
		sql.WriteString(" desc ")
	}
	if page.NeedPage() {
		sql.WriteString(" limit ")
		sql.WriteString(strconv.FormatInt(int64(page.GetPageSize()), 10))
		sql.WriteString(" offset ")
		sql.WriteString(strconv.FormatInt(int64(page.GetOffset()), 10))
	}
	_, err := ORM.Raw(sql.String()).QueryRows(&res)
	return res, err
}

func (c *Comment) Count() (int64, error) {
	qs := ORM.QueryTable(c)
	if c.Id != 0 {
		qs = qs.Filter("id", c.Id)
	}
	if c.ArticleId != 0 {
		qs = qs.Filter("article_id", c.ArticleId)
	}
	if c.UserId != 0 {
		qs = qs.Filter("user_id", c.UserId)
	}
	if c.Status != 0 {
		qs = qs.Filter("status", c.Status)
	} else {
		qs = qs.Filter("status", 1)
	}
	if c.RepliedUserId != 0 {
		qs = qs.Filter("replied_user_id", c.RepliedUserId)
	}
	if c.CommentId != 0 {
		qs = qs.Filter("comment_id", c.CommentId)
	}
	return qs.Count()
}

func (c *Comment) Delete() error {
	c.Status = -1
	_, err := ORM.Update(c, "status")
	return err
}

func (c *Comment) Update() error {
	if c.Id == 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "id不应该为0")
	}
	var cols = make([]string, 0, 4)
	if c.ArticleId != 0 {
		cols = append(cols, "article_id")
	}
	if c.Status != 0 {
		cols = append(cols, "status")
	}
	if c.Content != "" {
		cols = append(cols, "content")
	}
	_, err := ORM.Update(c, cols...)
	return err
}
