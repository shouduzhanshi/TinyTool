package module

import "github.com/gorilla/websocket"

type Client struct {
	Ws *websocket.Conn
	AndroidId string
	HeartBeat int64
}
