package models

import (
	"bbs-back/base/common"
	"bbs-back/base/database/bbsRedis"
	"bbs-back/base/dto/information"
	"github.com/beego/beego/v2/client/orm"
	"strings"
)

type User struct {
	Id int64				`json:"id" form:"id" orm:"pk"`
	Name string				`json:"name" form:"name"`
	Password string			`json:"password" form:"password"`
	Account string			`json:"account" form:"account"`
	Email string			`json:"email" form:"email"`
	TelephoneNumber int64	`json:"telephoneNumber" form:"telephoneNumber"`
	Age int32				`json:"age" form:"age"`
	Status int32			`json:"status" form:"status"`
	common.TimeModel
	Role int32				`json:"role" form:"role"`
	Gender string			`json:"gender" form:"gender"`
	ImageUrl	string		`json:"imageUrl" form:"imageUrl"`
}

const (
	tableName = "user"
	USER_STATUS_NORMAL = 1	//状态， 1-正常使用 2-已注销 3-黑名单
	USER_STATUS_CANCELLATION = 2
	USER_STATUS_BLACKLIST = 3
	//	3 -超级管理员，2-管理员，1-普通用户
	USER_ROLE_USER        = 1
	USER_ROLE_ADMIN       = 2
	USER_ROLE_SUPER_ADMIN = 3
)

func (u *User) Read() (*User, error) {
	res := new(User)
	res.Id = u.Id
	err := ORM.Read(res)
	return res, err
}

func (u *User) FindOne() (*User, error) {
	res := new(User)
	err := u.createQsByParam().One(res)
	return res, err
}
var userOrderList = []string {
	"update_time",
	"create_time",
}
func (u *User) Find(page *common.Page, orderIndex int, isDesc bool, cols ...string) ([]*User, error) {
	qs := u.createQsByParam()
	if u.Id == 0 {
		// 用户列表查询时超级管理员不对外暴露
		qs = qs.Exclude("role", USER_ROLE_SUPER_ADMIN)
	}
	if orderIndex < 0 || orderIndex >= len(userOrderList){
		orderIndex = 0
	}
	if isDesc {
		qs = qs.OrderBy("-" + userOrderList[orderIndex])
	} else {
		qs = qs.OrderBy(userOrderList[orderIndex])
	}
	qs = qs.Limit(page.GetPageSize(), page.GetOffset())
	var userList []*User
	var err error
	if len(cols) > 0 {
		_, err = qs.All(&userList, cols...)
	} else {
		_, err = qs.All(&userList)
	}
	return userList, err
}

func (u *User) Count() (int64, error) {
	qs := u.createQsByParam()
	return qs.Count()
}

// 更据参数筛选
func (u *User) createQsByParam() orm.QuerySeter {
	qs := ORM.QueryTable(tableName)
	if u.Id != 0 {
		qs = qs.Filter("id", u.Id)
	}
	if u.Name != "" {
		qs = qs.Filter("name__icontains", u.Name)
	}
	if u.Password != "" {
		qs = qs.Filter("password", u.Password)
	}
	if u.Account != "" {
		qs = qs.Filter("account", u.Account)
	}
	if u.Email != "" {
		qs = qs.Filter("email", u.Email)
	}
	if u.Age != 0 {
		qs = qs.Filter("age", u.Age)
	}
	if u.Status != 0 {
		qs = qs.Filter("status", u.Status)
	}
	//else {
	//	qs = qs.Filter("status", USER_STATUS_NORMAL)
	//}
	if u.Role != 0 {
		qs = qs.Filter("role", u.Role)
	}

	if u.Gender != "" {
		qs = qs.Filter("gender", u.Gender)
	}
	qs = u.TimeFilter(qs)
	return qs
}

func (u *User) Insert() error {
	if u.Role == 0 {
		u.Role = 1
	}
	if u.Status == 0 {
		u.Status = 1
	}
	id, err := ORM.Insert(u)
	if err != nil {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, err.Error())
	}
	bbsRedis.Incr(information.TOTAL_USER_NUM)
	u.Id = id
	return err
}

func (u *User) Delete() error {
	u.Status = USER_STATUS_CANCELLATION
	_, err := ORM.Update(u, "status")
	return err
}

func (u *User) Update() error {
	if u.Id == 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "id不应该为0")
	}
	var cols []string
	// 根据主健更新
	if u.Name != "" {
		cols = append(cols, "name")
	}
	if u.Password != "" {
		cols = append(cols, "password")
	}
	if u.Email != "" {
		cols = append(cols, "email")
	}
	if u.Age != 0 {
		cols = append(cols, "age")
	}
	if u.Role != 0 {
		cols = append(cols, "role")
	}
	if u.Gender != "" {
		cols = append(cols, "gender")
	}
	if u.TelephoneNumber > 0 {
		cols = append(cols, "telephone_number")
	}
	if u.ImageUrl != "" {
		cols = append(cols, "image_url")
	}
	if u.Status != 0 {
		cols = append(cols, "status")
	}
	if len(cols) == 0 {
		return nil
	}
	_, err := ORM.Update(u, cols...)
	return err
}


// 获取我关注者的列表
func (u *User) FollowerList(cols ...string) []*User {
	var colsStr string
	if len(cols) > 0 {
		colsStr = "," + strings.Join(cols, ",")
	}
	sql := "select user.id" + colsStr +
		" from user inner join" +
		" (select followed_id as id from user inner join follower on user.id = follower.follower_id where user.status > 0 and user.id = ?)" +
		" as temp on user.id = temp.id"
	raw := ORM.Raw(sql, u.Id)
	var res []*User
	raw.QueryRows(&res)
	return res
}

// 关注我的列表
func (u *User) FolloweredList(cols ...string) []*User {
	var colsStr string
	if len(cols) > 0 {
		colsStr = "," + strings.Join(cols, ",")
	}
	sql := "select user.id " + colsStr +
		" from user inner join" +
		" (select follower_id as id from user inner join follower on user.id = follower.followed_id where user.status > 0 and user.id = ?)" +
		" as temp on user.id = temp.id"

	raw := ORM.Raw(sql, u.Id)
	var res []*User
	raw.QueryRows(&res)
	return res
}

func (u *User) Follower(id int64) error {
	var num int32
	ORM.Raw("select count(*) from follower where followed_id = ? and follower_id = ?", id, u.Id).QueryRow(&num)
	if num > 0 {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, "重复关注！！！")
	}
	_, err := ORM.Raw("insert into follower (followed_id, follower_id) values(?,?)", id, u.Id).Exec()
	return err
}

func (u *User) UnFollower(id int64) error {
	_, err := ORM.Raw("delete from follower where followed_id = ? and follower_id = ?", id, u.Id).Exec()
	return err
}