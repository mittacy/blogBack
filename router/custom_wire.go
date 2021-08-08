package router

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/app/data"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/service"
	"github.com/mittacy/blogBack/pkg/logger"
	"gorm.io/gorm"
)

func InitUserApi(db *gorm.DB, cache *redis.Pool, conf model.EmailConfig) api.User {
	customLogger := logger.NewCustomLogger("user")
	userData := data.NewUser(db, cache, customLogger)
	emailData := data.NewEmail(db, cache, conf, customLogger)
	userService := service.NewUser(userData, emailData, customLogger)
	userApi := api.NewUser(userService, customLogger)
	return userApi
}

func InitEmailApi(db *gorm.DB, cache *redis.Pool, conf model.EmailConfig) api.Email {
	customLogger := logger.NewCustomLogger("email")
	emailData := data.NewEmail(db, cache, conf, customLogger)
	emailService := service.NewEmail(emailData, customLogger)
	emailApi := api.NewEmail(emailService, customLogger)
	return emailApi
}

func InitAdminApi(db *gorm.DB) api.Admin {
	customLogger := logger.NewCustomLogger("admin")
	adminData := data.NewAdmin(db, customLogger)
	adminService := service.NewAdmin(adminData, customLogger)
	adminApi := api.NewAdmin(adminService, customLogger)
	return adminApi
}
