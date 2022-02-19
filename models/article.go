package models

import (
	"bbs-back/base/common"
	"bbs-back/base/dto/information"
	"bbs-back/base/storage"
	"bytes"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
	"time"
)

type Article struct {
	Id            int64           `json:"id" form:"id" orm:"pk"`
	Title         string          `json:"title" form:"title"`
	CategoryId    int64           `json:"categoryId" form:"categoryId"`
	UserId        int64           `json:"userId" form:"userId"`
	Status        int32           `json:"status" form:"status"` // -1 已删除 1-未发布，2-已经发布
	ReadCount     int32           `json:"readCount" form:"readCount"`
	CreateTime    common.DateTime `json:"createTime" form:"createTime" orm:"auto_now_add;type(datetime);"`
	UpdateTime    common.DateTime `json:"updateTime" form:"updateTime" orm:"column(update_time)"`
	Content       string          `json:"content" form:"content"`
	LabelList     []*Label        `json:"labelList" form:"labelList" orm:"-"`         // 关联标签
	FollowingFlag bool            `json:"followingFlag" form:"followingFlag" orm:"-"` // 是否关注
	FollowCount   int64           `json:"followCount" form:"followCount" orm:"-"`     // 关注数量
}

func (a *Article) Read() (*Article, error) {
	res := new(Article)
	res.Id = a.Id
	err := ORM.Read(res)
	if err != nil {
		return nil, err
	}
	res.ReadCount++
	_, err = ORM.Update(res, "read_count")
	if err != nil {
		return nil, err
	}
	storage.GetRedisPool().Incr(information.TOTAL_READ_NUM)
	res.LabelList = new(Label).FindByArticleId(res.Id)
	return res, err
}

func (a *Article) createQsByParam() orm.QuerySeter {
	qs := ORM.QueryTable(a)
	if a.Id != 0 {
		qs = qs.Filter("id", a.Id)
	}
	if a.Title != "" {
		qs = qs.Filter("title__icontains", a.Title)
	}
	if a.UserId != 0 {
		qs = qs.Filter("user_id", a.UserId)
	}
	if a.Status != 0 {
		qs = qs.Filter("status", a.Status)
	} else {
		qs = qs.Filter("status", 1)
	}
	if a.ReadCount != 0 {
		qs = qs.Filter("read_count", a.ReadCount)
	}
	if !a.CreateTime.IsZero() {
		qs = qs.Filter("create_time", a.CreateTime.Time)
	}
	if !a.UpdateTime.IsZero() {
		qs = qs.Filter("update_time", a.UpdateTime.Time)
	}
	return qs
}

var ORDER_BY_LIST = []string{
	"update_time",
	"create_time",
	"read_count",
}

func (a *Article) Find(page *common.Page, orderIndex int32, isAsc bool, rangeTime *common.RangeTime, labelIdList ...string) ([]*Article, error) {
	var res []*Article
	//var qs orm.QuerySeter
	//if page.NeedPage() {
	//	qs = a.createQsByParam().Limit(page.GetPageSize(), page.GetPageNum() - 1)
	//} else {
	//	qs = a.createQsByParam()
	//}
	//"id", "title", "category_id", "user_id", "status", "support_count", "read_count","create_time", "update_time"
	//_, err := qs.All(&res,cols...)

	sql := a.createRawSql("a.id, a.title, a.category_id, a.user_id, a.status, a.support_count, a.read_count, a.create_time, a.update_time", rangeTime, labelIdList...)
	sql.WriteString(" order by ")
	// 防止错误索引
	if orderIndex >= int32(len(ORDER_BY_LIST)) || orderIndex < 0 {
		orderIndex = 0
	}

	sql.WriteString(ORDER_BY_LIST[orderIndex])

	if !isAsc {
		sql.WriteString(" desc ")
	}
	sql.WriteString(" limit ? offset ?")
	_, err := ORM.Raw(sql.String(), page.GetPageSize(), page.GetOffset()).QueryRows(&res)
	if err != nil {
		return nil, err
	}
	l := new(Label)
	for _, item := range res {
		item.LabelList = l.FindByArticleId(item.Id)
	}
	return res, err
}

func (a *Article) Count(rangeTime *common.RangeTime, labelIdList ...string) (int64, error) {
	//qs := a.createQsByParam()
	//return qs.Count()
	var res int64
	sql := a.createRawSql("count(*)", rangeTime, labelIdList...)
	err := ORM.Raw(sql.String()).QueryRow(&res)
	return res, err
}

func (a *Article) createRawSql(cols string, rangeTime *common.RangeTime, labelIdList ...string) *bytes.Buffer {
	var sql = new(bytes.Buffer)
	sql.WriteString("select ")
	sql.WriteString(cols)
	sql.WriteString(" from article a  ")
	if a.FollowingFlag {
		sql.WriteString(" inner join article_follower af on a.id = af.article_id ")
	}
	if a.Status == 0 {
		sql.WriteString(" where a.status = 2")
	} else {
		sql.WriteString(" where a.status = ")
		sql.WriteString(strconv.FormatInt(int64(a.Status), 10))
	}

	if a.FollowingFlag {
		sql.WriteString(" and af.user_id = ")
		sql.WriteString(strconv.FormatInt(a.UserId, 10))
	} else if a.UserId != 0 {
		sql.WriteString(" and a.user_id = ")
		sql.WriteString(strconv.FormatInt(a.UserId, 10))
	}

	if a.Id != 0 {
		sql.WriteString(" and a.id = ")
		sql.WriteString(strconv.FormatInt(a.Id, 10))
	}

	if a.CategoryId != 0 {
		sql.WriteString(" and a.category_id = ")
		sql.WriteString(strconv.FormatInt(a.CategoryId, 10))
	}
	if a.Title != "" {
		sql.WriteString(" and a.title like \"%")
		sql.WriteString(a.Title)
		sql.WriteString("%\" ")
	}
	// 创建时间范围查询
	if rangeTime != nil && !rangeTime.StartTime.IsZero() && !rangeTime.EndTime.IsZero() {
		sql.WriteString(" and a.create_time between '")
		sql.WriteString(rangeTime.GetStartTime())
		sql.WriteString("' and '")
		sql.WriteString(rangeTime.GetEndTime())
		sql.WriteString("' ")
	}

	if len(labelIdList) > 0 {
		sql.WriteString(" and a.id in (select distinct la.article_id from label_article la  where  la.label_id in ( ")
		sql.WriteString(labelIdList[0])
		for i := 1; i < len(labelIdList); i++ {
			sql.WriteString(",")
			sql.WriteString(labelIdList[i])
		}
		sql.WriteString(" )) ")
	}
	return sql
}
func (a *Article) Insert(labelIdList ...string) error {
	if a.Status == 0 {
		a.Status = 2
	}
	id, err := ORM.Insert(a)
	a.Id = id
	if len(labelIdList) > 0 {
		a.insertLabelArticle(labelIdList...)
	}
	storage.GetRedisPool().Incr(information.TOTAL_ARTICLE_NUM)
	return err
}

func (a *Article) Delete() error {
	a.Status = -1
	a.UpdateTime.Time = time.Now()
	_, err := ORM.Update(a, "status")
	return err
}

func (a *Article) Update(labelIdList ...string) error {
	var cols []string
	if a.Id != 0 {
		cols = append(cols, "id")
	}
	if a.Title != "" {
		cols = append(cols, "title")
	}
	if a.UserId != 0 {
		cols = append(cols, "user_id")
	}
	if a.Status != 0 {
		cols = append(cols, "status")
	}
	if a.ReadCount != 0 {
		cols = append(cols, "read_count")
	}
	if a.CategoryId != 0 {
		cols = append(cols, "category_id")
	}
	if a.Content != "" {
		cols = append(cols, "content")
	}
	if !a.UpdateTime.IsZero() {
		cols = append(cols, "update_time")
	}
	_, err := ORM.Update(a, cols...)
	if err != nil {
		return err
	}
	if len(labelIdList) > 0 {
		a.deleteLabelArticle()
		a.insertLabelArticle(labelIdList...)
	}
	return err
}

// 新建标签关联关系
func (a *Article) insertLabelArticle(labelIdList ...string) error {
	if len(labelIdList) == 0 {
		return nil
	}
	var sql = new(bytes.Buffer)
	sql.WriteString("insert into label_article (label_id, article_id) values ")
	sql.WriteString(" (")
	sql.WriteString(labelIdList[0])
	sql.WriteString(" , ")
	articleId := strconv.FormatInt(a.Id, 10)
	sql.WriteString(articleId)
	sql.WriteString(" )")
	for index, length := 1, len(labelIdList); index < length; index++ {
		sql.WriteString(", (")
		sql.WriteString(labelIdList[index])
		sql.WriteString(" , ")
		sql.WriteString(articleId)
		sql.WriteString(" )")
	}
	_, err := ORM.Raw(sql.String()).Exec()
	return err
}

// 删除标签关联关系
func (a *Article) deleteLabelArticle() error {
	var sql = new(bytes.Buffer)
	sql.WriteString("delete from label_article where article_id = ? ")
	_, err := ORM.Raw(sql.String(), a.Id).Exec()
	return err
}

func (a *Article) Follow(userId int64) error {
	if a.Id == 0 || userId == 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "论贴或用户id不为空")
	}
	_, err := ORM.Raw("insert into article_follower (article_id, user_id) values (?, ?)", a.Id, userId).Exec()

	if err != nil {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, err.Error())
	}

	return err
}

func (a *Article) FollowingCount(userId ...int64) int64 {
	var res = new(int64)
	if len(userId) == 1 {
		ORM.Raw("select count(*) from article_follower where article_id = ? and user_id = ?", a.Id, userId).QueryRow(res)
	} else {
		ORM.Raw("select count(*) from article_follower where article_id = ? ", a.Id).QueryRow(res)
	}
	return *res
}

func (a *Article) UnFollow(userId int64) error {
	if a.Id == 0 || userId == 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "论贴或用户id不为空")
	}
	_, err := ORM.Raw("delete from article_follower where article_id = ? and user_id = ?", a.Id, userId).Exec()
	return err
}

func (a *Article) TotalReadNum() (int64, error) {
	var res = new(int64)
	err := ORM.Raw("select sum(read_count) from article where status = 2").QueryRow(res)
	return *res, err
}