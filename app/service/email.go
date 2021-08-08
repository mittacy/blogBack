package service

import (
	"github.com/mittacy/blogBack/app/api"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/utils"
	"github.com/pkg/errors"
	"strings"
)

type Email struct {
	emailData IEmailData
	logger *logger.CustomLogger
}

// 编写实现api层中的各个service接口的构建方法

func NewEmail(emailData IEmailData, logger *logger.CustomLogger) api.IEmailService {
	return &Email{
		emailData: emailData,
		logger: logger,
	}
}

type IEmailData interface {
	GetEmailTpl(name string) (*model.EmailTpl, error)
	SaveCode(email string, code string) error
	GetCode(email string) (string, error)
	InvalidCode(email string) error
	SendEmail(mailTo []string, subject string, body string) error
}

func (ctl *Email) SendRegisterCode(email string) error {
	// 1. 查询邮件模板
	emailName := model.EmailRegisterTplName
	tpl, err := ctl.emailData.GetEmailTpl(emailName)
	if err != nil {
		return err
	}

	// 2. 生成验证码，存入redis
	code := utils.RandCode(6)
	if err := ctl.emailData.SaveCode(email, code); err != nil {
		return err
	}

	// 3. 处理模板，填充数据
	tpl.Content = strings.Replace(tpl.Content, "${{code}}", code, 1)

	// 4. 发送邮件
	if err := ctl.emailData.SendEmail([]string{email}, "注册验证码", tpl.Content); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
