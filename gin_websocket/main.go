package main

import (
	"fmt"
	"gin_websocket/controller"
	"gin_websocket/websocket"
	"gin_websocket/websocket_wstoken"
	"github.com/gin-gonic/gin"
	//"gin_websocket/websocket"
)

//server
func main() {

	gin.SetMode(gin.ReleaseMode) //线上环境

	r := gin.Default()
	r.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动 websocket
	go websocket.Manager.Start()
	r.GET("/ws", websocket.WsHandler)

	// 启动 websocket
	go websocket_wstoken.Manager.Start()
	r.GET("/token/ws", websocket_wstoken.WsHandler)
	r.POST("/sendMsgToUser", controller.WebSocketMsg{}.SendMsgToUser)

	// 启动 gin 服务
	fmt.Println("[listen]0.0.0.0:8081")
	r.Run(":8081")
}
