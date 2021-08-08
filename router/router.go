package router

import (
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/pkg/config"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/store/cache"
	"github.com/mittacy/blogBack/pkg/store/db"
	"github.com/spf13/viper"
	"time"
)

func InitRouter(r *gin.Engine) {
	emailConf := model.EmailConfig{}
	if err := viper.UnmarshalKey("email", &emailConf); err != nil {
		panic(fmt.Sprintf("checkout the email config: %s\n", err))
	}

	// 1. 初始化控制器
	emailApi := InitEmailApi(db.ConnectGorm("blog"), cache.ConnRedis("blog"), emailConf)
	userApi := InitUserApi(db.ConnectGorm("blog"), cache.ConnRedis("blog"), emailConf)
	adminApi := InitAdminApi(db.ConnectGorm("blog"))

	// 2. 全局中间件
	r.Use(ginzap.Ginzap(logger.GetRequestLogger(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.GetRequestLogger(), true))
	//r.Use(middleware.CorsMiddleware())

	// 3. 初始化路由
	relativePath := "/api/" + config.ServerConfig.Version
	g := r.Group(relativePath) // 统一前缀
	{
		/**
		 * 不需要登录的Api
		 */
		// 登录
		g.POST("/session/admin/login", adminApi.Login)
		g.POST("/session/user/login", userApi.Login)

		// 邮件
		email := g.Group("/email")
		{
			email.GET("/register_code", emailApi.GetRegisterCode)
		}

		// 用户
		user := g.Group("/user")
		{
			user.POST("", userApi.Register)
			user.GET("/:id", userApi.GetInfo)
		}
	}
}
