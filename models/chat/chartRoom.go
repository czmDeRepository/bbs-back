package chat

import (
	"time"

	"bbs-back/base/storage"
	"bbs-back/models/dao"
	"bbs-back/models/monitor"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"
)

const (
	EVENT_INIT = iota
	EVENT_JOIN
	EVENT_LEAVE
	EVENT_MESSAGE

	ONLINE_NUM
	CHANGE_CATEGORY
	CHAT_MESSAGE
	HEART_BEAT
)

type EventType int

type Event struct {
	CategoryId int64     `json:"categoryId"` // 主题id
	Type       EventType `json:"type"`       // JOIN, LEAVE, MESSAGE
	User       *dao.User `json:"user"`
	Timestamp  int       `json:"timestamp"`
	Content    *Content  `json:"content"`
	OnLineNum  int64     `json:"onLineNum"`
}

type MessageType int

type Content struct {
	CategoryIdList []int64     `json:"categoryIdList"`
	Type           MessageType `json:"type"`
	Message        string      `json:"message"`
	NewCategoryId  int64       `json:"newCategoryId"`
}

type Subscribe struct {
	User       *dao.User
	Ws         *websocket.Conn
	CategoryId int64
}

var (
	// Channel for new join users.
	subscribe = make(chan *Subscribe, 10)
	// Channel for exit users.
	unsubscribe = make(chan *Subscribe, 10)
	// Send events here to publish them.
	Publish = make(chan *Event, 10)
	// subscribers CategoryId->userId->Subscribe
	subscribersMap = map[int64]map[int64]*Subscribe{}
)

func init() {
	go chatroom()
}

func OnLineNum() int64 {
	var res int64 = 0
	for _, item := range subscribersMap {
		res += int64(len(item))
	}
	return res
}

// 获取当天最大群聊人数
func GetMaxOnlineNum() int64 {
	res, err := storage.GetRedisPool().BitCount(monitor.GetDateKey(monitor.CHART_MAX_ONLINE_NUM))
	if err != nil {
		logs.Error("chat: get %s fail: %s", monitor.GetDateKey(monitor.CHART_MAX_ONLINE_NUM), err.Error())
	}
	return res
}

// 获取群聊消息数
func GetMessageNum() int64 {
	messageNum, err := storage.GetRedisPool().GetInt64(monitor.GetDateKey(monitor.CHAT_MESSAGEG_NUM))
	if err != nil {
		logs.Error("chat get %s fail: %s", monitor.GetDateKey(monitor.CHAT_MESSAGEG_NUM), err.Error())
		return 0
	}
	return messageNum
}

func NewEvent(categoryId int64, ep EventType, user *dao.User, content *Content) *Event {
	return &Event{categoryId, ep, user, int(time.Now().Unix()), content, 0}
}

func Join(categoryId int64, user *dao.User, ws *websocket.Conn) *Subscribe {
	sub := &Subscribe{CategoryId: categoryId, User: user, Ws: ws}
	subscribe <- sub
	return sub
}

func Leave(sub *Subscribe) {
	unsubscribe <- sub
}

func addSubscribe(sub *Subscribe, isJoin bool) {
	if subscribersMap[sub.CategoryId] == nil {
		subscribersMap[sub.CategoryId] = map[int64]*Subscribe{}
	} else if subscribersMap[sub.CategoryId][sub.User.Id] != nil {
		oldSub := subscribersMap[sub.CategoryId][sub.User.Id]
		// 覆盖连接标志
		oldSub.User.Status = -100
		logs.Critical("Old user:", oldSub.User.Name, ";WebSocket:", oldSub.Ws != nil)
		if oldSub.Ws != nil {
			oldSub.Ws.Close()
		}
	}
	subscribersMap[sub.CategoryId][sub.User.Id] = sub
	if isJoin {
		Publish <- NewEvent(sub.CategoryId, EVENT_JOIN, sub.User, nil)
		logs.Critical("New user:", sub.User.Name, ";WebSocket:", sub.Ws != nil)
		storage.GetRedisPool().SetBit(monitor.GetDateKey(monitor.CHART_MAX_ONLINE_NUM), sub.User.Id)
	}
}

// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			addSubscribe(sub, true)
		case event := <-Publish:
			handlePublishEvent(event)
			if event.Type == EVENT_MESSAGE {
				logs.Critical("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <-unsubscribe:
			if subscribersMap[unsub.CategoryId][unsub.User.Id] != nil {
				// Clone connection.
				ws := unsub.Ws
				if ws != nil {
					ws.Close()
					logs.Critical("WebSocket closed:", unsub)
				}
				delete(subscribersMap[unsub.CategoryId], unsub.User.Id)
				if len(subscribersMap[unsub.CategoryId]) == 0 {
					delete(subscribersMap, unsub.CategoryId)
				}
				Publish <- NewEvent(unsub.CategoryId, EVENT_LEAVE, unsub.User, nil) // Publish a LEAVE event.
			}
		}

	}
}

func handlePublishEvent(event *Event) {
	//data, err := json.Marshal(event)
	//if err != nil {
	//	logs.Critical("Fail to marshal event:", err)
	//	return
	//}
	if event.Content != nil {
		switch event.Content.Type {
		// 获取在线人数
		case ONLINE_NUM:
			message := new(Event)
			message.Type = EVENT_INIT
			message.OnLineNum = OnLineNum()
			subscriber := subscribersMap[event.CategoryId][event.User.Id]
			subscriber.Ws.WriteJSON(message)
		//	修改订阅主题
		case CHANGE_CATEGORY:
			userSubMap := subscribersMap[event.CategoryId]
			sub := userSubMap[event.User.Id]
			delete(userSubMap, event.User.Id)
			sub.CategoryId = event.Content.NewCategoryId
			addSubscribe(sub, false)
		case CHAT_MESSAGE:
			// 广播消息
			if checkIsBroadcast(event) {
				broadcast(event)
			} else {
				directionalPush(event)
			}
			storage.GetRedisPool().Incr(monitor.GetDateKey(monitor.CHAT_MESSAGEG_NUM))
		}
	} else {
		broadcast(event)
	}

}
func checkIsBroadcast(event *Event) bool {
	if len(event.Content.CategoryIdList) == 0 {
		return true
	}
	for _, id := range event.Content.CategoryIdList {
		if id == 0 {
			return true
		}
	}
	return false
}

// 定向分类推送消息
func directionalPush(event *Event) {
	// 定向发送
	for _, categoryId := range event.Content.CategoryIdList {
		userSubMap := subscribersMap[categoryId]
		for _, sub := range userSubMap {
			ws := sub.Ws
			if ws != nil {
				if sub.User.Id == event.User.Id {
					continue
				}
				// User disconnected.
				if ws.WriteJSON(event) != nil {
					unsubscribe <- sub
				}
			} else {
				logs.Critical("delete:", *sub)
				delete(subscribersMap[event.CategoryId], sub.User.Id)
			}
		}
	}
	// 转发给订阅所有的用户
	for _, sub := range subscribersMap[0] {
		ws := sub.Ws
		if ws != nil {
			if sub.User.Id == event.User.Id {
				continue
			}
			if ws.WriteJSON(event) != nil {
				unsubscribe <- sub
			}
		} else {
			logs.Critical("delete:", *sub)
			delete(subscribersMap[event.CategoryId], sub.User.Id)
		}
	}
}

// 广播加入 或 退出 或消息发送分类为0
func broadcast(event *Event) {
	for _, userSubMap := range subscribersMap {
		for _, sub := range userSubMap {
			ws := sub.Ws
			if ws != nil {
				if sub.User.Id == event.User.Id {
					// 首次加入返回在线人数
					if event.Type == EVENT_JOIN {
						join := new(Event)
						join.Type = EVENT_INIT
						join.OnLineNum = OnLineNum()
						ws.WriteJSON(join)
					}
					continue
				}
				// User disconnected.
				if ws.WriteJSON(event) != nil {
					unsubscribe <- sub
				}
			} else {
				logs.Critical("delete:", *sub)
				delete(subscribersMap[event.CategoryId], sub.User.Id)
			}
		}
	}
}
