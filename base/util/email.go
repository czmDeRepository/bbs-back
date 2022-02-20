package util

import (
	"encoding/json"
	"errors"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
	beego "github.com/beego/beego/v2/server/web"
)

var emailConf string

func InitEmail() {
	host := beego.AppConfig.DefaultString("email.host", "smtp.126.com")
	port := beego.AppConfig.DefaultInt("email.port", 25)
	username := beego.AppConfig.DefaultString("email.username", "czmdemailbox@126.com")
	password := beego.AppConfig.DefaultString("email.password", "QCFCIDRZHNFWOOBH")
	conf, err := json.Marshal(map[string]interface{}{
		"host":     host,
		"port":     port,
		"username": username,
		"password": password,
	})
	if err != nil {
		logs.Error("init email err: %v", err)
		return
	}
	emailConf = string(conf)
}

// 发送邮件
func SendEmail(content string, target ...string) error {
	reqMail := utils.NewEMail(emailConf)
	if reqMail == nil {
		return errors.New("utils.NewEMail err")
	}
	reqMail.Text = content
	reqMail.Subject = "在线论坛：验证码"
	reqMail.To = target
	return reqMail.Send()
}
