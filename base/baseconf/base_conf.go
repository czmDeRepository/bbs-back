package baseconf

import (
	"log"

	beego "github.com/beego/beego/v2/server/web"
)

// 文件限制大小
var MAX_FILE_LIMIT int64

func Init() {

	maxFileLimit, err := beego.AppConfig.Int64("maxFileLimit")
	// 10485760 = 1024 * 1024 * 4   5M
	if err != nil {
		MAX_FILE_LIMIT = 1024 * 1024 * 5
		log.Println("获取上传文件大小限制配置失败！！使用默认值10M")
	} else {
		MAX_FILE_LIMIT = maxFileLimit
	}
}
