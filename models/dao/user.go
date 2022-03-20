package dao

import (
	"strings"

	"bbs-back/base/common"
	"bbs-back/base/dto/information"
	"bbs-back/base/storage"
	"bbs-back/base/util"
	"bbs-back/models/entity"

	"github.com/beego/beego/v2/client/orm"
)

type User entity.User

const (
	tableName                = "user"
	USER_STATUS_NORMAL       = 1 //状态， 1-正常使用 2-已注销 3-黑名单
	USER_STATUS_CANCELLATION = 2
	USER_STATUS_BLACKLIST    = 3
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
	qs := u.createQsByParam()
	if u.Password != "" {
		// 密码加密
		password, err := util.DecryptPassword(u.Password)
		if err != nil {
			return nil, err
		}
		u.Password = password
		qs = qs.Filter("password", u.Password)
	}
	err := qs.One(res)
	return res, err
}

var userOrderList = []string{
	"update_time",
	"create_time",
}

func (u *User) Find(page *common.Page, orderIndex int, isDesc bool, cols ...string) ([]*User, error) {
	qs := u.createQsByParam()
	if u.Id == 0 {
		// 用户列表查询时超级管理员不对外暴露
		qs = qs.Exclude("role", USER_ROLE_SUPER_ADMIN)
	}
	if orderIndex < 0 || orderIndex >= len(userOrderList) {
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
	// 密码加密
	password, err := util.DecryptPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = password
	id, err := ORM.Insert(u)
	if err != nil {
		return common.NewErrorWithCode(common.ERROR_DB_LIMIT, err.Error())
	}
	storage.GetRedisPool().Incr(information.TOTAL_USER_NUM)
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
		// 密码加密
		password, err := util.DecryptPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = password
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
	// 电话允许置空
	cols = append(cols, "telephone_number")
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
