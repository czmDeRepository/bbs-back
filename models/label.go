package models

import (
	"bbs-back/base/common"
	"bbs-back/models/entity"
	"github.com/beego/beego/v2/client/orm"
)

type Label entity.Label

func (l *Label) Read() (*Label, error) {
	res := new(Label)
	res.Id = l.Id
	err := ORM.Read(res)
	return res, err
}

var labelOrderList = []string{
	"update_time",
	"create_time",
}

func (l *Label) Find(page *common.Page, orderIndex int, isDesc bool, cols ...string) ([]*Label, error) {
	var qs orm.QuerySeter
	if page == nil || !page.NeedPage() {
		qs = l.createQsByParam()
	} else {
		qs = l.createQsByParam().Limit(page.GetPageSize(), page.GetOffset())
	}
	if orderIndex < 0 || orderIndex >= len(labelOrderList) {
		orderIndex = 0
	}
	if isDesc {
		qs = qs.OrderBy("-" + labelOrderList[orderIndex])
	} else {
		qs = qs.OrderBy(labelOrderList[orderIndex])
	}
	var res []*Label
	_, err := qs.All(&res, cols...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (l *Label) Count() (int64, error) {
	return l.createQsByParam().Count()
}

func (l *Label) createQsByParam() orm.QuerySeter {
	qs := ORM.QueryTable(l)
	if l.Id != 0 {
		qs = qs.Filter("id", l.Id)
	}
	if l.Name != "" {
		qs = qs.Filter("name__icontains", l.Name)
	}
	if l.Status == 0 {
		qs = qs.Filter("status", 1)
	} else {
		qs = qs.Filter("status", l.Status)
	}
	qs = l.TimeFilter(qs)
	return qs
}

func (l *Label) Insert() error {
	l.Status = 1
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	isSuccess, _, err := ORM.ReadOrCreate(l, "name")
	if err != nil {
		return err
	}
	if isSuccess {
		return nil
	}
	return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "标签名已经存在")
}

func (l *Label) Update() error {
	if l.Name == "" {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "标签名非空！")
	}
	if l.Id == 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "id非空！")
	}
	num, err := new(Label).createQsByParam().Filter("name", l.Name).Count()
	if err != nil {
		return err
	}
	if num != 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "标签名已经存在！")
	}
	var cols []string
	if l.Status != 0 {
		cols = append(cols, "status")
	}
	if l.Name != "" {
		cols = append(cols, "name")
	}
	_, err = ORM.Update(l, cols...)
	return err
}

func (l *Label) Delete() error {
	if l.Id == 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "id不可为空")
	}
	l.Status = -1
	_, err := ORM.Update(l, "status")
	return err
}

func (l *Label) FindByArticleId(articleId int64) []*Label {
	var res []*Label
	ORM.Raw("select l.* from label l inner join label_article la on l.`id` = la.`label_id` "+
		"where l.`status` = 1 and la.`article_id` = ? ", articleId).QueryRows(&res)
	return res
}
