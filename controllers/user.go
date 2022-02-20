package controllers

import (
	"strconv"
	"time"

	"bbs-back/base/common"
	"bbs-back/base/dto"
	"bbs-back/base/storage"
	"bbs-back/base/util"
	"bbs-back/models/dao"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
)

// Operations about Users
type UserController struct {
	BaseController
}

// @Title Get
// @Description get user by id
// @Param	id		path 	int64	true "The key for staticblock"
// @Success 200 {object} dto.Result
// @Failure 403 :id is empty
// @router /:id [get]
func (controller *UserController) Get() {
	id, err := controller.GetInt64(":id")
	if err != nil {
		controller.end(common.ErrorMeWithCode(err.Error(), common.ERROR_PARAM))
		return
	}
	// 当前token的用户id
	curUserId := controller.getCurUserId()
	curUserRole := controller.getCurUserRole()
	if curUserId != id && curUserRole == 1 {
		controller.end(common.ErrorMeWithCode(common.ERROR_MESSAGE[common.ERROR_CURRENT_USER], common.ERROR_CURRENT_USER))
		return
	}
	u := new(dao.User)
	u.Id = id
	res, err := u.Read()
	if err != nil {
		controller.end(common.Error(err))
	}
	controller.end(common.SuccessWithData(res))
}

// @Title GetAll
// @Description 获取用户列表
// @Param 	name	query	string	false	"user name"
// @Param	account	query	string	false	"user account"
// @Param	role	query	int32	false	"user role -1: 超级管理员 0：管理员 1：普通用户"
// @Param	status	query	int32	false   "user status 1-正常使用 2-已注销 3-黑名单"
// @Param	pageNum	query	int32	false
// @Param	pageSize	query	int32	false
// @Success 200 {object} dto.Result
// @Failure 403 :pageNum or pageSize is empty
// @router	/ [get]
func (controller *UserController) GetAll() {

	param := new(dto.UserDto)
	controller.ParseForm(param)

	// 多用户获取需要权限
	if param.Id == 0 {
		role := controller.getCurUserRole()
		if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
			controller.end(common.ErrorDetail(common.ERROR_POWER, common.ERROR_MESSAGE[common.ERROR_POWER]))
			return
		}
	}

	orderIndex, _ := controller.GetInt("orderIndex", 0)
	isDesc, _ := controller.GetBool("isDesc", false)
	userList, err := param.User.Find(&param.Page, orderIndex, isDesc)
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	total, _ := param.User.Count()
	res := new(common.PageDto)
	res.PageNum = param.GetPageNum()
	res.PageSize = int32(len(userList))
	res.Data = userList
	res.Total = total
	controller.end(common.SuccessWithData(res))
}

// @Title Put
// @Description 更新用户信息
// @Param	id		body	int64	true	"user id"
// @Param 	name	body	string	true	"user name"
// @Param	password	body	string	true	"user password"
// @Param	account	body	string	true	"user account"
// @Param	email	body	string	false	"邮箱"
// @Param	telephoneNumber	body	string	false	"手机号"
// @Param	age		body	int32	false	"年龄"
// @Param	status	body	int32	false   "user status 1-正常使用 2-已注销 3-黑名单"
// @Param	role	body	int32	false	"user role -1: 超级管理员 1：管理员 2：普通用户"
// @Param	createTime	body	string	false	"创建时间"
// @Param	updateTime	body	string	false	"更新时间"
// @Param	gender	body	string	false	"性别"
// @Success 200 {object} dto.Result
// @router	/ [put]
func (controller *UserController) Put() {

	param := new(dao.User)
	controller.ParseForm(param)
	if param.Id == 0 {
		controller.end(common.ErrorDetail(common.ERROR_PARAM, "用户id为空"))
		return
	}
	role := controller.getCurUserRole()

	if role != dao.USER_ROLE_ADMIN && role != dao.USER_ROLE_SUPER_ADMIN {
		curUserId := controller.getCurUserId()
		if curUserId != param.Id {
			controller.end(common.ErrorDetail(common.ERROR_CURRENT_USER, "不可修改其他用户信息！！！"))
			return
		}
	}
	err := param.Update()

	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())
}

// @Title Delete
// @Description 删除用户
// @Param 	id	query 	string	false	"user name"
// @Param	account	query	string	false	"user account"
// @Param	role	query	int32	false	"user role -1: 超级管理员 0：管理员 1：普通用户"
// @Param	status	query	int32	false   "user status 1-正常使用 2-已注销 3-黑名单"
// @Param	pageNum	query	int32	true
// @Param	pageSize	query	int32	true
// @Success 200 {object} dto.Result
// @Failure 403 :pageNum or pageSize is empty
// @router	/ [delete]
func (controller *UserController) Delete() {
	u := new(dao.User)
	id, err := controller.GetInt64("id")
	if err != nil {
		controller.end(common.ErrorMeWithCode(err.Error(), common.ERROR_PARAM))
		return
	}
	u.Id = id
	err = u.Delete()
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.Success())
}

var SecretKey string

func init() {
	SecretKey, _ = beego.AppConfig.String("secretKey")
}

// @Title Post
// @Description 新增用户
// @Param 	name	formData	string	true	"user name"
// @Param	password	formData	string	true	"user password"
// @Param	account	formData	string	true	"user account"
// @Param	email	formDate	string	false	"邮箱"
// @Param	telephoneNumber	fromData	string	false	"手机号"
// @Param	age		formData	int32	false	"年龄"
// @Param	status	formData	int32	false   "user status 1-正常使用 2-已注销 3-黑名单"
// @Param	role	formData	int32	false	"user role -1: 超级管理员 1：管理员 2：普通用户"
// @Param	createTime	fromData	string	false	"创建时间"
// @Param	updateTime	fromData	string	false	"更新时间"
// @Param	gender	formDate	string	false	"性别"
// @Param	ImageUrl	formDate	string	false	"头像路径"
// @Success 200 {object} dto.Result
// @router	/ [post]
func (controller *UserController) Post() {
	u := new(dao.User)
	controller.ParseForm(u)
	if u.Name == "" || u.Account == "" || u.Account == "" {
		controller.end(common.ErrorWithMe("缺少必要参数"))
	}
	err := u.Insert()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	token, createTime, _ := createToken(u)
	storage.GetRedisPool().HSet(USER_TOKEN_CREATE_TIME_KEY, strconv.FormatInt(u.Id, 10), createTime)
	controller.end(common.SuccessWithData(token))
}

// @Title Login
// @Description login system
// @Param	account	formData	string	true	"user account"
// @Param	password	formData	string	true	"user password"
// @Success 200 {object} dto.Result
// @Failure 403 :id is empty
// @router /login [get]
func (controller *UserController) Login() {
	captcha := controller.GetString("captcha")
	if captcha == "" {
		controller.end(common.ErrorWithMe("验证码非空"))
		return
	}
	var user *dao.User
	var err error
	// 邮箱验证登录
	if email := controller.GetString("email"); email != "" {
		user, err = (&dao.User{Email: email}).FindOne()
		if err != nil {
			controller.end(common.Error(err))
			return
		}
		content, err := storage.GetRedisPool().Get(util.GetEmailKey(user.Account))
		if err != nil && err == redis.ErrNil || content == "" {
			controller.end(common.ErrorWithMe("验证码已失效"))
			return
		}
		if captcha != content {
			controller.end(common.ErrorWithMe("验证码错误"))
			return
		}
		storage.GetRedisPool().Del(util.GetEmailKey(user.Account))
	} else {
		// 账号密码登录
		account := controller.GetString("account")
		password := controller.GetString("password")
		if account == "" || password == "" {
			controller.end(common.ErrorWithMe("账号密码非空"))
			return
		}

		if err := util.VerifyCaptcha(controller.GetString("captchaKey"), captcha, true); err != nil {
			controller.end(common.Error(err))
			return
		}
		param := new(dao.User)
		param.Account = account
		param.Password = password
		user, err = param.FindOne()
		if err != nil {
			controller.end(common.ErrorWithMe("用户不存在或密码错误 "))
			return
		}
	}
	switch user.Status {
	case dao.USER_STATUS_BLACKLIST:
		controller.end(common.ErrorWithMe("您目前被限制登陆，请联系管理员解锁！！！"))
		return
	case dao.USER_STATUS_CANCELLATION:
		controller.end(common.ErrorWithMe("当前账号已注销！！！"))
		return
	}

	tokenString, createTime, err := createToken(user)
	if err != nil {
		controller.end(common.ErrorMeWithCode(err.Error(), common.ERROR_TOKEN_CREATE))
		return
	}
	storage.GetRedisPool().HSet(USER_TOKEN_CREATE_TIME_KEY, strconv.FormatInt(user.Id, 10), createTime)
	controller.end(common.SuccessWithData(tokenString))
}

func createToken(user *dao.User) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// 一天过期
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	createTime := time.Now().Unix()
	claims["iat"] = createTime
	claims["id"] = user.Id
	claims["role"] = user.Role
	claims["name"] = user.Name
	claims["gender"] = user.Gender
	claims["imageUrl"] = user.ImageUrl
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(SecretKey))
	return tokenString, createTime, err
}

// @Title Refresh
// @router /refresh [put]
func (controller *UserController) Refresh() {
	userId := controller.getCurUserId()
	user := new(dao.User)
	user.Id = userId
	user.Status = dao.USER_STATUS_NORMAL
	curUser, err := user.FindOne()
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	token, createTime, err := createToken(curUser)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	storage.GetRedisPool().HSet(USER_TOKEN_CREATE_TIME_KEY, strconv.FormatInt(userId, 10), createTime)
	controller.end(common.SuccessWithData(token))
}

// @Title Follow
// @Description 获取关注列表
// @Param	id		Query 	int64	true "The key for staticblock"
// @Param	followFlag	Query	bool	true	"关注我的人和我关注的人标识,true:关注我的人"
// @Success 200 {object} dto.Result
// @Failure 1000 :缺少参数或存在非法参数
// @router /follow [get]
func (controller *UserController) FollowGet() {
	followFlag, err := controller.GetBool("followFlag")
	if err != nil {
		controller.end(common.ErrorMeWithCode(err.Error(), common.ERROR_PARAM))
		return
	}
	u := new(dao.User)
	id := controller.getCurUserId()
	u.Id = id
	var list []*dao.User
	if followFlag {
		list = u.FollowerList("name", "account")
	} else {
		list = u.FolloweredList("name", "account")
	}
	controller.end(common.SuccessWithData(list))
}

// @Title FollowPost
// @Description 获取关注列表
// @Param	targetId	Query	int64	true	"被关注者id"
// @Success 200 {object} dto.Result
// @Failure 1000 :缺少参数或存在非法参数
// @router /follow [post]
func (controller *UserController) FollowPost() {
	targetId, err := controller.GetInt64("targetId")
	if err != nil {
		controller.end(common.ErrorMeWithCode(err.Error(), common.ERROR_PARAM))
		return
	}
	u := new(dao.User)
	id := controller.getCurUserId()
	u.Id = id
	err = u.Follower(targetId)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	controller.end(common.Success())

}

// @Title FollowDelete
// @Description 获取关注列表
// @Param	targetId	Query	int64	true	"被关注者id"
// @Success 200 {object} dto.Result
// @Failure 1000 :缺少参数或存在非法参数
// @router /follow [delete]
func (controller *UserController) FollowDelete() {
	targetId, err := controller.GetInt64("targetId")
	if err != nil {
		controller.end(common.ErrorMeWithCode(err.Error(), common.ERROR_PARAM))
		return
	}
	u := new(dao.User)
	id := controller.getCurUserId()
	u.Id = id
	err = u.UnFollower(targetId)
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.Success())
}
