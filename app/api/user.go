package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/app/transform"
	"github.com/mittacy/blogBack/pkg/logger"
)

type User struct {
	userService IUserService
	transform   transform.User
	logger      *logger.CustomLogger
}

func NewUser(userService IUserService, logger *logger.CustomLogger) User {
	return User{
		userService: userService,
		transform: transform.NewUser(logger),
		logger:    logger,
	}
}

type IUserService interface {
	Ping()
}

func (ctl *User) Ping(c *gin.Context) {}

