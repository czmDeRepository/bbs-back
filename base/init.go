package base

import (
	"bbs-back/base/baseconf"
	"bbs-back/base/storage"
)

func Init() {
	baseconf.Init()
	storage.Init()
}
