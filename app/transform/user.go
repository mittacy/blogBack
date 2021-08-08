package transform

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/validator/userValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
)

type User struct {
	logger *logger.CustomLogger
}

func NewUser(customLogger *logger.CustomLogger) User {
	return User{logger: customLogger}
}

// UserPack 数据库数据转化为响应数据
// @param data 数据库数据
// @return reply 响应体数据
// @return err
func (ctl *User) UserPack(data *model.User) (*userValidator.GetReply, error) {
	reply := userValidator.GetReply{}

	if err := copier.Copy(&reply, data); err != nil {
		return nil, err
	}

	return &reply, nil
}

// GetReply 详情响应包装
// @param data 数据库数据
func (ctl *User) GetReply(c *gin.Context, data *model.User) {
	reply, err := ctl.UserPack(data)
	if err != nil {
		response.CopierErrAndLog(c, ctl.logger, err)
		return
	}

	res := map[string]interface{}{
		"user": reply,
	}

	response.Success(c, res)
}
