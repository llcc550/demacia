package main

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	dialer := websocket.Dialer{}
	u := url.URL{Scheme: "ws", Host: "124.71.163.192:32303", Path: "/"}
	connect, _, err := dialer.Dial(u.String(), nil)
	if nil != err {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	for {
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			fmt.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			fmt.Println(string(messageData))
		case websocket.BinaryMessage:
			fmt.Println(messageData)
		case websocket.CloseMessage:
			fmt.Println(string(messageData))
		case websocket.PingMessage:
			fmt.Println(string(messageData))
		case websocket.PongMessage:
			fmt.Println(string(messageData))
		default:
		}
	}
}
