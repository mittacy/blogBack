package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/bootstrap"
	"github.com/mittacy/blogBack/pkg/config"
	"github.com/mittacy/blogBack/router"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func init() {
	bootstrap.Init()
}

func main() {
	r := gin.New()

	// 初始化路由
	router.InitRouter(r)

	serverConfig := config.ServerConfig
	s := &http.Server{
		Addr: ":" + strconv.Itoa(serverConfig.Port),
		Handler: r,
		ReadTimeout: time.Second * serverConfig.ReadTimeout,
		WriteTimeout: time.Second * serverConfig.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	zap.S().Infof("监听端口:%d", serverConfig.Port)

	s.ListenAndServe()
}
