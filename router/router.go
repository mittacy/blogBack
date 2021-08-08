package router

import (
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/middleware"
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
	categoryApi := InitCategoryApi(db.ConnectGorm("blog"))
	articleApi := InitArticleApi(db.ConnectGorm("blog"), cache.ConnRedis("blog"))

	// 2. 全局中间件
	r.Use(ginzap.Ginzap(logger.GetRequestLogger(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.GetRequestLogger(), true))
	r.Use(middleware.CorsMiddleware())

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

		// 分类
		g.GET("/categories", categoryApi.List)

		// 文章
		g.GET("/article/:id", articleApi.Get)
		g.GET("/articles", articleApi.List)
		g.GET("/articles_home", articleApi.HomeList)

		/**
		 * 需要登录的Api
		 */
		needAuth := g.Group("")
		needAuth.Use(middleware.ParseToken())
		{
			authCategory := needAuth.Group("/category")
			{
				authCategory.POST("", middleware.Operate(middleware.ActionAddCategory), categoryApi.Create)
				authCategory.DELETE("/:id", middleware.Operate(middleware.ActionDeleteCategory), categoryApi.Delete)
				authCategory.PUT("", middleware.Operate(middleware.ActionPutCategory), categoryApi.Update)
			}

			authArticle := needAuth.Group("/article")
			{
				authArticle.POST("", middleware.Operate(middleware.ActionAddArticle), articleApi.Create)
				authArticle.DELETE("/:id", middleware.Operate(middleware.ActionDeleteArticle), articleApi.Delete)
				authArticle.PUT("", middleware.Operate(middleware.ActionPutArticle), articleApi.Update)
			}
		}
	}
}
