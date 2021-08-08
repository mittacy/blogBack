package service

import (
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/pkg/logger"
)

type User struct {
	userData IUserData
	logger *logger.CustomLogger
}

// 编写实现api层中的各个service接口的构建方法

func NewUser(userData IUserData, logger *logger.CustomLogger) api.IUserService {
	return &User{
		userData: userData,
		logger: logger,
	}
}

type IUserData interface {
	PingData()
}

func (ctl *User) Ping() {}

