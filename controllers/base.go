package controllers

import (
	"bbs-back/base/common"
	"bbs-back/base/database/bbsRedis"
	"bbs-back/base/enum"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"strconv"
	"time"
)

type BaseController struct {
	beego.Controller
}

const (
	RESPONSE = "json"
	USER_TOKEN_CREATE_TIME_KEY = "USER_TOKEN_KEY"
)

func (controller *BaseController) Prepare() {
	//log.Println("prepare")
}


// 解析时间参数
func (controller *BaseController) ParseForm(obj interface{}) error {
	form, err := controller.Input()
	if err != nil {
		return err
	}
	res := beego.ParseForm(form, obj)
	setDateTime(controller.GetString("createTime"), obj,"CreateTime")
	setDateTime(controller.GetString("updateTime"), obj,"UpdateTime")
	return res
}

// 设置时间
func setDateTime(timeValue string, obj interface{}, name string)  {
	if timeValue != ""  {
		v := reflect.ValueOf(obj)
		elem := v.Elem()
		fieldV := elem.FieldByName(name)
		if fieldV.IsValid() {
			t := common.DateTime{}
			t.Time, _ = time.ParseInLocation(common.FormatDateTime, timeValue, time.Local)
			fieldV.Set(reflect.ValueOf(t))
		}
	}
}

func (controller *BaseController) end(data common.Result)  {
	controller.Data[RESPONSE] = data
	controller.ServeJSON()
}

// 解析token
func ValidateToken(tokenString string)  (*jwt.Token, error)  {
	parse, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			//验证是否是给定的加密算法
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return false, nil
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		validationError := err.(*jwt.ValidationError)
		switch validationError.Errors {
		case jwt.ValidationErrorExpired:
			return nil, common.NewErrorWithCode(common.ERROR_TOKEN_EXPIRE, err.Error())
		default:
			return nil, common.NewErrorWithCode(common.ERROR_TOKEN_PARSE, err.Error())
		}
	}
	return parse, nil
}

func (controller *BaseController) getCurUserId() int64 {
	header := controller.Ctx.Request.Header
	curUserId := header.Get(enum.CUR_USER_ID)
	if curUserId == "" {
		token := header.Get("token")
		if token == "" {
			panic(common.NewErrorWithCode(common.ERROR_TOKEN_NOEXIST,common.ERROR_MESSAGE[common.ERROR_TOKEN_NOEXIST]))
		}
		parseToken, err := ValidateToken(token)
		if err != nil {
			panic(err)
		}
		claims := parseToken.Claims.(jwt.MapClaims)
		userId := checkIsOnLineUser(claims)
		header.Add(enum.CUR_USER_ID, userId)
		header.Add(enum.CUR_USER_ROLD, fmt.Sprintf("%v", claims["role"]))
		return controller.getCurUserId()
	}
	id, _ := strconv.ParseInt(header.Get(enum.CUR_USER_ID), 10, 64)
	return id
}
func checkIsOnLineUser(claims jwt.MapClaims) string {
	userId := fmt.Sprintf("%v", claims["id"])
	createTime, _ := redis.Int64(bbsRedis.HGet(USER_TOKEN_CREATE_TIME_KEY, userId))
	if createTime != int64(claims["iat"].(float64)) {
		panic(common.NewErrorWithCode(common.ERROR_TOKEN_EXPIRE, common.ERROR_MESSAGE[common.ERROR_TOKEN_EXPIRE]))
	}
	return userId
}
func (controller *BaseController) getCurUserRole() int64 {
	header := controller.Ctx.Request.Header
	curUserRole := header.Get(enum.CUR_USER_ROLD)
	if curUserRole == "" {
		controller.getCurUserId()
	}
	roleId, _ := strconv.ParseInt(controller.Ctx.Request.Header.Get(enum.CUR_USER_ROLD), 10, 64)
	return roleId
}

func (controller *BaseController) paramError(err error) {
	controller.end(common.ErrorMeWithCode(err.Error(), common.ERROR_PARAM))
}