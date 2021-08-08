package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mittacy/blogBack/app/transform"
	"github.com/mittacy/blogBack/app/validator/emailValidator"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/response"
)

type Email struct {
	emailService IEmailService
	transform   transform.Email
	logger      *logger.CustomLogger
}

func NewEmail(emailService IEmailService, logger *logger.CustomLogger) Email {
	return Email{
		emailService: emailService,
		transform: transform.NewEmail(logger),
		logger:    logger,
	}
}

type IEmailService interface {
	SendRegisterCode(email string) error
}

/**
 * @apiVersion 0.1.0
 * @apiGroup Email
 * @api {get} /email/register_code 获取邮箱验证码
 * @apiName Email.GetRegisterCode
 *
 * @apiParam {string} email 邮箱
 *
 */
func (ctl *Email) GetRegisterCode(c *gin.Context) {
	d := emailValidator.RegisterCodeReq{}
	if err := c.ShouldBindQuery(&d); err != nil {
		response.ValidateErr(c, err)
		return
	}

	if err := ctl.emailService.SendRegisterCode(d.Email); err != nil {
		response.CheckErrAndLog(c, ctl.logger, "getEmailCode", err)
		return
	}

	response.Success(c, nil)
}

