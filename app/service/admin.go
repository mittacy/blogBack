package service

import (
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/validator/adminValidator"
	"github.com/mittacy/blogBack/pkg/jwt"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/utils"
	"github.com/pkg/errors"
)

type Admin struct {
	adminData IAdminData
	logger *logger.CustomLogger
}

// 编写实现api层中的各个service接口的构建方法

func NewAdmin(adminData IAdminData, logger *logger.CustomLogger) api.IAdminService {
	return &Admin{
		adminData: adminData,
		logger: logger,
	}
}

type IAdminData interface {
	GetByName(name string) (*model.Admin, error)
}

func (ctl *Admin) Login(login adminValidator.AdminLoginReq) (string, error) {
	// 1. 查询数据库中用户的信息
	realAdmin, err := ctl.adminData.GetByName(login.Name)
	if err != nil {
		if errors.Is(err, apierr.ErrUserNoExist) { // 隐藏错误信息，不让用户知道是账号不存在
			err = apierr.ErrUserOrPassword
		}
		return "", err
	}

	// 2. 校验密码
	if utils.EncryptionBySalt(login.Password, realAdmin.Salt) != realAdmin.Password {
		return "", apierr.ErrUserOrPassword
	}

	// 3. 生成 token
	token, err := jwt.Token.Create(realAdmin.Id, model.UserRoleAdmin)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

