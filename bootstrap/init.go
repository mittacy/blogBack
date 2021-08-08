package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/pkg/checker"
	"github.com/mittacy/blogBack/pkg/config"
	"github.com/mittacy/blogBack/pkg/jwt"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/store/cache"
	"go.uber.org/zap"
)

func Init() {
	// 1. 初始化配置文件
	config.InitViper()

	// 2. 设置gin的运行模式
	gin.SetMode(config.ServerConfig.Env)

	// 3. 初始化全局日志
	logger.Init()

	// 4. 初始化校验翻译器
	if err := checker.InitTrans(); err != nil {
		zap.L().Panic("初始化校验翻译器失败", zap.String("reason", err.Error()))
	}

	// 5. 初始化 Cache 配置
	cache.Init()

	// 5. 初始化token
	tokenCache := cache.ConnCustomRedis("blog", "token")
	jwt.InitToken(tokenCache)
}
