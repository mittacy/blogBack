
package data

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mittacy/blogBack/app/service"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/store/cache"
	"gorm.io/gorm"
)

// 实现service层中的data接口

type User struct {
	db 	   *gorm.DB
	cache  cache.CustomRedis
	logger *logger.CustomLogger
}

func NewUser(db *gorm.DB, cacheConn *redis.Pool, logger *logger.CustomLogger) service.IUserData {
	r := cache.ConnRedisByPool(cacheConn, "user")

	return &User{
		db:    	db,
		cache: 	r,
		logger: logger,
	}
}

func (ctl *User) PingData() {}

