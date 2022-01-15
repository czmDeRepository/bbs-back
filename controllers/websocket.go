package controllers

import "C"
import (
	"bbs-back/models"
	"bbs-back/models/chat"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type WebSocketController struct {
	BaseController
}

var upgrader *websocket.Upgrader

func init()  {
	upgrader = new(websocket.Upgrader)
	// 跨域
	upgrader.CheckOrigin = func(_ *http.Request) bool {
		return true
	}
}

// @Title Join
// @Success
// @Failure 1000 :参数错误
// @router	/chat [get]
func (controller *WebSocketController) Join() {
	token := controller.GetString("token")
	parseToken, err := ValidateToken(token)
	if err != nil {
		http.Error(controller.Ctx.ResponseWriter, "token解析出错！！!" + err.Error(), http.StatusBadRequest)
		return
	}
	userId := checkIsOnLineUser(parseToken.Claims.(jwt.MapClaims))

	ws, err := upgrader.Upgrade(controller.Ctx.ResponseWriter, controller.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(controller.Ctx.ResponseWriter, "请使用websocket连接！！！", http.StatusBadRequest)
		return
	}
	id , _ := strconv.ParseInt(userId, 10, 64)
	user := new(models.User)
	user.Id= id
	user, _ = user.Read()
	user.Password = ""
	sub := chat.Join(0, user, ws)
	defer func() {
		// 被新建连接覆盖
		if user.Status != -100 {
			chat.Leave(sub)
		}
	}()

	// Message receive loop.
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return
		}
		content := new(chat.Content)
		json.Unmarshal(message, content)
		if content.Type == chat.HEART_BEAT {
			// 心跳
			continue
		}
		chat.Publish <- chat.NewEvent(sub.CategoryId, chat.EVENT_MESSAGE, user, content)
	}
}
