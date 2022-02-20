package base

import (
	"bbs-back/base/baseconf"
	"bbs-back/base/storage"
	"bbs-back/base/util"
)

func Init() {
	baseconf.Init()
	storage.Init()
	util.InitCaptcha()
	util.InitEmail()
}
