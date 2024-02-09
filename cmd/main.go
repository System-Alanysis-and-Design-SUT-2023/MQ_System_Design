package main

import "github.com/gorilla/websocket"

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
