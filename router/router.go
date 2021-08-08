package router

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/pkg/config"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/store/cache"
	"github.com/mittacy/blogBack/pkg/store/db"
	"time"
)

func InitRouter(r *gin.Engine) {
	// 1. 初始化控制器
	userApi := InitUserApi(db.ConnectGorm("blog"), cache.ConnRedis("blog"))

	// 2. 全局中间件
	r.Use(ginzap.Ginzap(logger.GetRequestLogger(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.GetRequestLogger(), true))
	//r.Use(middleware.CorsMiddleware())

	// 3. 初始化路由
	relativePath := "/api/" + config.ServerConfig.Version
	g := r.Group(relativePath) // 统一前缀
	{
		g.POST("/user", userApi.Create)
		g.DELETE("/user", userApi.Delete)
		g.PUT("/user", userApi.Update)
		g.GET("/user", userApi.Get)
		g.GET("/users", userApi.List)
	}
}