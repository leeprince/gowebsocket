package websocket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

// Message is return msg
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Manager define a ws server manager
var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}


func (c *Client) Health() {
	var (
		err error
	)
	for {
		if err = c.Socket.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
			return
		}
		time.Sleep(60 * time.Second)
	}
}

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context) {
	uid := c.Query("uid")
	touid := c.Query("to_uid")

	// 完成ws协议的握手操作
	// Upgrade: 升级将 HTTP 服务器连接升级到 WebSocket 协议。
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	//可以添加用户信息验证
	client := &Client{
		ID:     creatId(uid, touid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	Manager.Register <- client
	// 接收消息
	go client.Read()
	// 发送消息
	go client.Write()
	// 启动线程，心跳检测
	go client.Health()
}

// Start 项目运行前, 协程开启 start -> go Manager.Start()
func (manager *ClientManager) Start() {
	for {
		log.Println("<---管道通信:websocket.go--->")
		select {
		case conn := <-Manager.Register:
			log.Printf("新用户加入:%v", conn.ID)
			Manager.Clients[conn.ID] = conn
			jsonMessage, _ := json.Marshal(&Message{Content: "Successful connection to socket service"})
			conn.Send <- jsonMessage
		case conn := <-Manager.Unregister:
			log.Printf("用户离开:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected"})
				conn.Send <- jsonMessage
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case message := <-Manager.Broadcast:
			MessageStruct := Message{}
			json.Unmarshal(message, &MessageStruct)
			for id, conn := range Manager.Clients {
				if id != creatId(MessageStruct.Recipient, MessageStruct.Sender) {
					continue
				}
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
		}
	}
}

func creatId(uid, touid string) string {
	return uid + "_" + touid
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		log.Printf("读取到客户端的信息:%s", string(message))
		Manager.Broadcast <- message
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("发送到客户端的信息:%s", string(message))

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

