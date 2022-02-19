package controllers

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"bbs-back/base/baseconf"
	"bbs-back/base/common"
	"bbs-back/base/dto/information"
	"bbs-back/base/enum"
	"bbs-back/base/storage"
	"bbs-back/base/util"
	"bbs-back/models/chat"

	beego "github.com/beego/beego/v2/server/web"
)

type PublicController struct {
	BaseController
}

func init() {
	_, err := os.Stat(enum.RESOURCES_IMAGE_PATH)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(enum.RESOURCES_IMAGE_PATH, 0666)
		if err != nil {
			log.Fatal("资源目录["+enum.RESOURCES_IMAGE_PATH+"] 不存在且创建失败！:", err)
		}
	}
	// 静态资源图片映射
	beego.SetStaticPath("/v1/resources", enum.RESOURCES_PATH)
}

// @Title	Upload
// @Param	image	File	true	"上传文件"
// @router /upload [post]
func (controller *PublicController) Upload() {
	controller.getCurUserId()
	// key is the file name
	file, fileHeader, err := controller.GetFile("image")
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	// don't forget to close
	defer file.Close()
	if fileHeader.Size > baseconf.MAX_FILE_LIMIT {
		controller.end(common.ErrorDetail(common.ERROR_PARAM, fmt.Sprintf("上传文件超过最大限制: max = %dM", baseconf.MAX_FILE_LIMIT/(1024*1024))))
		return
	}
	bufData := make([]byte, fileHeader.Size)
	file.Read(bufData)
	hashName := fmt.Sprintf("%x", md5.Sum(bufData))
	filePath := enum.RESOURCES_IMAGE_PATH + hashName + path.Ext(fileHeader.Filename)
	_, err = os.Stat(filePath)
	if err == nil {
		// 文件已存在，不再io
		controller.end(common.SuccessWithData(filePath[1:]))
		return
	}
	newFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		controller.end(common.HandleError(err))
		return
	}
	// don't forget to close
	defer newFile.Close()
	writer := bufio.NewWriter(newFile)
	writer.Write(bufData)
	//buf := make([]byte, 1024)
	//var index int64 = 0
	//for ; index < fileHeader.Size; {
	//	len, err := file.ReadAt(buf, index)
	//	// 读取到末尾退出
	//	if err == io.EOF {
	//		writer.Write(buf[:len])
	//		break
	//	}
	//	index += int64(len)
	//	writer.Write(buf)
	//}
	writer.Flush()
	controller.end(common.SuccessWithData(filePath[1:]))
}

// @Title	Information
// @router  /information [get]
func (controller *PublicController) Information() {
	totalReadNumStr, _ := storage.GetRedisPool().Get(information.TOTAL_READ_NUM)
	totalReadNum, _ := strconv.ParseInt(totalReadNumStr, 10, 64)
	totalUserNumStr, _ := storage.GetRedisPool().Get(information.TOTAL_USER_NUM)
	totalUserNum, _ := strconv.ParseInt(totalUserNumStr, 10, 64)
	totalArticleNumStr, _ := storage.GetRedisPool().Get(information.TOTAL_ARTICLE_NUM)
	totalArticleNum, _ := strconv.ParseInt(totalArticleNumStr, 10, 64)
	totalChatNum := chat.OnLineNum()
	res := new(information.Information)
	res.TotalReadNum = totalReadNum
	res.TotalUserNum = totalUserNum
	res.TotalArticleNum = totalArticleNum
	res.TotalChatPersonNum = totalChatNum
	controller.end(common.SuccessWithData(res))
}

// @Title	GetCaptcha
// @router  /captcha [get]
func (controller *PublicController) GetCaptcha() {
	id, image, err := util.CreateCaptchaBase64()
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.SuccessWithData(beego.M{
		"captchaKey": id,
		"image":      image,
	}))
}

// @Title	PutCaptcha
// @router  /captcha [put]
func (controller *PublicController) PutCaptcha() {
	captchaKey := controller.GetString("captchaKey")
	var id, image string
	var err error
	if captchaKey != "" {
		id, image, err = util.CreateCaptchaBase64(captchaKey)
	} else {
		id, image, err = util.CreateCaptchaBase64()
	}
	if err != nil {
		controller.end(common.Error(err))
		return
	}
	controller.end(common.SuccessWithData(beego.M{
		"captchaKey": id,
		"image":      image,
	}))
}
