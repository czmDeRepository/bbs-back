package models

import (
	"bbs-back/base/common"
	"bbs-back/models/entity"
	"github.com/beego/beego/v2/client/orm"
)

type Category entity.Category

func (c *Category) Insert() error {
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	isSuccess, _, err := ORM.ReadOrCreate(c, "name")
	if err != nil {
		return err
	}
	if isSuccess {
		return nil
	}
	return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "分类名已经存在")
}
func (c *Category) Read() (*Category, error) {
	res := new(Category)
	res.Id = c.Id
	err := ORM.Read(res)
	return res, err
}

var categoryOrderList = []string{
	"update_time",
	"create_time",
}

func (c *Category) Find(page *common.Page, orderIndex int, isDesc bool, cols ...string) ([]*Category, error) {
	var qs orm.QuerySeter
	if page == nil || !page.NeedPage() {
		qs = c.createQsByParam()
	} else {
		qs = c.createQsByParam().Limit(page.GetPageSize(), page.GetOffset())
	}
	// 防止错误索引
	if orderIndex >= len(categoryOrderList) || orderIndex < 0 {
		orderIndex = 0
	}
	if isDesc {
		qs = qs.OrderBy("-" + categoryOrderList[orderIndex])
	} else {
		qs = qs.OrderBy(categoryOrderList[orderIndex])
	}

	var res []*Category
	_, err := qs.All(&res, cols...)
	return res, err
}

func (c *Category) Count() (int64, error) {
	return c.createQsByParam().Count()
}

func (c *Category) createQsByParam() orm.QuerySeter {
	qs := ORM.QueryTable(Category{})
	if c.Id != 0 {
		qs = qs.Filter("id", c.Id)
	}
	if c.Name != "" {
		qs = qs.Filter("name__icontains", c.Name)
	}
	qs = c.TimeFilter(qs)
	return qs
}

func (c *Category) Update() error {
	if c.Name == "" {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "分类名非空！")
	}
	if c.Id == 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "id非空！")
	}
	num, err := new(Category).createQsByParam().Filter("name", c.Name).Count()
	if err != nil {
		return err
	}
	if num != 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "分类名已经存在！")
	}
	_, err = ORM.Update(c, "name")
	return err
}

func (c *Category) Delete() (int64, error) {
	return c.createQsByParam().Delete()
}
