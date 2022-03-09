package module

import "golang.org/x/net/websocket"

type Client struct {
	Ws *websocket.Conn
	AndroidId string
}
