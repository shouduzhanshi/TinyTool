package server

import (
	"bytes"
	"container/list"
	"github.com/gorilla/websocket"
	"github.com/martinlindhe/notify"
	"time"
	"tiny_tool/log"
	"tiny_tool/module"
)

var onlineUser = list.New()

var channel = make(chan UserChannel)

type Msg struct {
	Data     []byte
	Classify int
}

func publishMsg(data []byte, Classify int) {
	channel <- UserChannel{
		Type:     2,
		Data:     data,
		Classify: Classify,
	}
}

func updateHeartBeat(id string) {
	channel <- UserChannel{
		Type: 3,
		Data: bytes.NewBufferString(id).Bytes(),
	}
}

func init() {
	go func() {
		for {
			data := <-channel
			if data.Type == 0 {
				_online(data.WS, data.AndroidId)
			} else if data.Type == 1 {
				_offline(data)
			} else if data.Type == 2 {
				_publish(data)
			} else if data.Type == 3 {
				_updateHeartBeat(bytes.NewBuffer(data.Data).String())
			}
		}
	}()
}

func _updateHeartBeat(s string) {
	for i := onlineUser.Front(); i != nil; i = i.Next() {
		client := i.Value.(*module.Client)
		if client.AndroidId == s {
			client.HeartBeat = time.Now().UnixMilli()
		}
	}
}

func _publish(data UserChannel) {
	if onlineUser.Len() <= 0 {
		if data.Classify!=0 {
			offlineNotify()
		}
		return
	}
	for i := onlineUser.Front(); i != nil; i = i.Next() {
		client := i.Value.(*module.Client)
		var currentTime int64
		currentTime = time.Now().UnixMilli() - int64(5000)
		if currentTime < client.HeartBeat {
			err := client.Ws.WriteMessage(websocket.TextMessage, data.Data)
			if err != nil {
				data.WS = client.Ws
				_offline(data)
			}
		} else {
			data.WS = client.Ws
			_offline(data)
		}
	}
}

func offlineNotify() {
	log.E("All devices are offline, please check the device network connection!")
	notify.Notify("Tiny CLI", "warning", "All devices are offline, please check the device network connection", "")
}

func _online(ws *websocket.Conn, id string) {
	onlineUser.PushBack(&module.Client{
		Ws:        ws,
		AndroidId: id,
		HeartBeat: time.Now().UnixMilli(),
	})
	log.V("online device quantity ", onlineUser.Len())
}

func _offline(data UserChannel) {
	if onlineUser.Len() <= 0 {
		return
	}
	ws := data.WS
	ws.Close()
	var element *list.Element
	for i := onlineUser.Front(); i != nil; i = i.Next() {
		if (i.Value.(*module.Client)).Ws == ws {
			element = i
			break
		}
	}
	if element != nil {
		onlineUser.Remove(element)
		log.V("device ", (element.Value.(*module.Client)).AndroidId, " offline")
	}
	log.V("online device quantity ", onlineUser.Len())
	if onlineUser.Len() <= 0 {
		offlineNotify()
	}
}

func online(AndroidId string, ws *websocket.Conn) {
	channel <- UserChannel{
		WS:        ws,
		Type:      0,
		AndroidId: AndroidId,
	}
}

func offline(ws *websocket.Conn) {
	channel <- UserChannel{
		WS:   ws,
		Type: 1,
	}
}

type UserChannel struct {
	Type      uint
	WS        *websocket.Conn
	Data      []byte
	AndroidId string
	Classify  int
}
