package util

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"bbs-back/base/storage"

	"github.com/beego/beego/v2/core/logs"
	captcha "github.com/mojocn/base64Captcha"
)

var driverString *captcha.DriverString

func Init() {
	driverString = &captcha.DriverString{
		Height:          80,
		Width:           240,
		NoiseCount:      20,
		ShowLineOptions: 2,
		Length:          5,
		BgColor:         nil,
		Source:          "1234567890zxcvbnmlkjhgfdsaqwertyuiopZXCVBNMLKJHGFDSAQWERTYUIOP",
	}
}

// CreateCaptchaBase64 返回base64格式
func CreateCaptchaBase64(ids ...string) (id, image string, err error) {
	id, content, _ := driverString.GenerateIdQuestionAnswer()
	if len(ids) > 0 {
		id = ids[0]
	}
	drawCaptcha, err := driverString.DrawCaptcha(content)
	if err != nil {
		logs.Error("captcha.create: %v", err)
		return
	}
	// 1分钟过期
	err = storage.GetRedisPool().SetExp(PreKey(id), content, time.Minute)
	if err != nil {
		return
	}
	return id, drawCaptcha.EncodeB64string(), err
}

// CreateCaptcha 直接写入io流
func CreateCaptcha(w io.Writer, ids ...string) (id string, err error) {
	id, content, _ := driverString.GenerateIdQuestionAnswer()
	if len(ids) > 0 {
		id = ids[0]
	}
	drawCaptcha, err := driverString.DrawCaptcha(content)
	if err != nil {
		logs.Error("captcha.create: %v", err)
		return
	}
	// 1分钟过期
	err = storage.GetRedisPool().SetExp(PreKey(id), content, time.Minute)
	if err != nil {
		return
	}
	_, err = drawCaptcha.WriteTo(w)
	return id, err
}

func VerifyCaptcha(id, param string, rm ...bool) error {
	content, err := storage.GetRedisPool().Get(PreKey(id))
	if err != nil {
		logs.Error("captcha.verify: %v", err)
		return err
	}
	if content == "" {
		return errors.New("验证码已失效")
	}
	if strings.TrimSpace(param) != content {
		return errors.New("验证码错误")
	}
	if len(rm) > 0 && rm[0] {
		storage.GetRedisPool().Del(PreKey(id))
	}
	return nil
}

func PreKey(id string) string {
	return fmt.Sprintf("captcha:%s", id)
}
