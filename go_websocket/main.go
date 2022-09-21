package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn      *websocket.Conn
		messageType int
		err         error
		conn        *websocket.Conn
		data        []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	
	// 启动线程，心跳检测
	go func() {
		var (
			err error
		)
		for {
			if err = wsConn.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()
	
	for {
		messageType, data, err = wsConn.ReadMessage()
		if err != nil {
			fmt.Println("ReadMessage err: ", err)
			goto ERR
		}
		fmt.Printf("messageType:%v - data:%v", messageType, string(data))
		if err = wsConn.WriteMessage(messageType, data); err != nil {
			fmt.Println("WriteMessage err: ", err)
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("[listen]0.0.0.0:7777")
	
	http.ListenAndServe("0.0.0.0:7777", nil)
}
