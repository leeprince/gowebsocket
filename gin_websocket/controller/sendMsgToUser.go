package controller

import (
	"gin_websocket/params"
	"gin_websocket/websocket_wstoken"
	"github.com/gin-gonic/gin"
)

type WebSocketMsg struct {
}

func (s WebSocketMsg) SendMsgToUser(c *gin.Context) {
	var sendMsgToUserReq params.SendMsgToUserReq

	if err := c.ShouldBindJSON(&sendMsgToUserReq); err != nil {
		c.JSON(-1, gin.H{
			"message": "fail",
		})
		return
	}

	// 发送消息到已通过 wstoken 建立的 websocket 连接上
	wstoken := sendMsgToUserReq.WsToken
	data := sendMsgToUserReq.Data
	conn, ok := websocket_wstoken.Manager.Clients[wstoken]
	if ok == false {
		return
	}
	conn.Send <- []byte(data)

	c.JSON(200, gin.H{
		"message": "successful.",
	})
}
