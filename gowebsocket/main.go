package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.yewifi.com/golden-cloud/common/tracer"
	"gowebsocket/config"
	"gowebsocket/logger"
	"gowebsocket/router"
	"gowebsocket/services"
	"net/http"
	"time"
)

func init() {
	logger.SetLogLevel(config.Config.LogLevel)
	logger.InitLog()
	services.InitWebSocket()
	gin.SetMode(config.Config.Mode)
	gin.DefaultErrorWriter = logger.GetLogger().Logger.Writer()

}

//server
func main() {
	tracer.InitTracer(config.Config.JaegerAgentURL, "kuaiji-gowebsocket")
	engine := router.InitRouter()
	addr := fmt.Sprintf("%s:%v", config.Config.GinHost.Host, config.Config.GinHost.Port)
	info := fmt.Sprintf("开始监听 %s 环境配置 %s 日志级别 %s", addr, config.Config.Env, logger.GetLogger().Logger.GetLevel())
	logger.Info("0000000000000", info)
	fmt.Println(info)

	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServe()
}
