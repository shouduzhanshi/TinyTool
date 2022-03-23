package server

import (
	"container/list"
	"golang.org/x/net/websocket"
	"time"
	"tiny_tool/log"
	"tiny_tool/module"
)

var onlineUser = list.New()

var channel = make(chan UserChannel)

func publishMsg(data []byte, start int64) {
	channel <- UserChannel{
		Type:  2,
		Data:  data,
		Start: start,
	}
}

func init() {
	go func() {
		for {
			data := <-channel
			if data.Type == 0 {
				_online(data.WS,data.AndroidId)
			} else if data.Type == 1 {
				_offline(data)
			} else if data.Type == 2 {
				_publish(data)
			}
		}
	}()
}

func _publish(data UserChannel) {
	for i := onlineUser.Front(); i != nil; i = i.Next() {
		client := i.Value.(module.Client)
		go func(client *module.Client) {
			size, err := client.Ws.Write(data.Data)
			if err != nil {
				offline(client.Ws)
			} else if data.Start != 0 {
				end := time.Now().UnixNano()
				log.E("总耗时 ", (end-data.Start)/1e6," ms")
				log.V("send data to ",client.AndroidId," ", size," bytes")
			}
		}(&client)
	}
}

func _online(ws *websocket.Conn, id string) {
	onlineUser.PushBack(module.Client{
		Ws: ws,
		AndroidId: id,
	})
	log.V("online device quantity ", onlineUser.Len())
}

func _offline(data UserChannel) {
	ws := data.WS
	ws.Close()
	var element *list.Element
	for i := onlineUser.Front(); i != nil; i = i.Next() {
		if (i.Value.(module.Client)).Ws == ws {
			element = i
			break
		}
	}
	if element != nil {
		onlineUser.Remove(element)
		log.V("device ",(element.Value.(module.Client)).AndroidId," offline")
	}
	log.V("online device quantity ", onlineUser.Len())
}

func online(AndroidId string, ws *websocket.Conn) {
	channel <- UserChannel{
		WS:   ws,
		Type: 0,
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
	Type  uint
	WS    *websocket.Conn
	Data  []byte
	Start int64
	AndroidId string
}
