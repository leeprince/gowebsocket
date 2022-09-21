# 一、websocket
> 生产项目：gowebsocket


## （一）一对一聊天
### １. 文件
```
js_websocket.html
websocket.go
```

### 2. 核心代码
main.go
```
// 启动 websocket
go websocket.Manager.Start()
r.GET("/ws", websocket.WsHandler)
```

### curl
```
http://localhost:63344/gin_websocket/js_websocket.html?uid=1&to_uid=2
```
```
http://localhost:63344/gin_websocket/js_websocket.html?uid=2&to_uid=1
```



## （二）转发消息到已通过 wstoken 建立的连接上
##　1. 文件
```
js_websocket_wstoken.html
websocket_wstoken.go
```
### 2. 核心代码
main.go
```
	// 启动 websocket
	go websocket_wstoken.Manager.Start()
	r.GET("/pushws", websocket_wstoken.WsHandler)
	r.POST("/sendMsgToUser", controller.WebSocketMsg{}.SendMsgToUser)

```

### 3. curl
```
curl --location --request POST 'http://127.0.0.1:8080/pushws' \
--header 'Content-Type: application/json' \
--data-raw '{
    "wstoken":"leeprince",
    "data":"{\"dd\":1}"
}'
```
```
curl -X POST 'http://127.0.0.1:8080/pushws' \
-H 'Content-Type: application/json' \
-d '{
    "wstoken":"leeprince",
    "data":"{\"dd\":1}"
}'
```
```
curl --location --request POST 'http://127.0.0.1:8080/pushws' \
--header 'Content-Type: application/json' \
--data-raw '{
    "wstoken":"BEvmoi06nLDSQbzJSBnfWOxXznneSuAt1154764793",
    "data":"{\"id\":100,\"msg_type\":1,\"msg_title\":\"标题001\",\"msg_content\":\"内容002\",\"created_at\":\"2021-08-31 15:42:00\",\"msg_mark\":\"11\",\"read_status\":2}"
}'
```