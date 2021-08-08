package router

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/app/data"
	"github.com/mittacy/blogBack/app/service"
	"github.com/mittacy/blogBack/pkg/logger"
	"gorm.io/gorm"
)

func InitUserApi(db *gorm.DB, cache *redis.Pool) api.User {
	customLogger := logger.NewCustomLogger("user")
	iUserService := data.NewUser(db, cache, customLogger)
	apiIUserService := service.NewUser(iUserService, customLogger)
	user := api.NewUser(apiIUserService, customLogger)
	return user
}
